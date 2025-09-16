package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the MQTT TUI application
	app := NewApp()

	// Create the Bubble Tea program with options for proper terminal handling
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
	)

	// Set the program reference in MQTT client for sending messages
	if app.mqtt != nil {
		app.mqtt.SetProgram(p)
	}

	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
} // App represents the main application state
type App struct {
	mqtt     *MQTTClient
	ui       *UI
	config   Config
	quitting bool
}

// Config holds the MQTT broker configuration
type Config struct {
	BrokerURL string
	Username  string
	Password  string
	ClientID  string
}

// NewApp creates a new application instance
func NewApp() *App {
	config := Config{
		BrokerURL: getEnvOrDefault("MQTT_BROKER", "tcp://localhost:1883"),
		Username:  getEnvOrDefault("MQTT_USERNAME", ""),
		Password:  getEnvOrDefault("MQTT_PASSWORD", ""),
		ClientID:  getEnvOrDefault("MQTT_CLIENT_ID", "mqttui"),
	}

	app := &App{
		config: config,
		ui:     NewUI(),
	}

	// Initialize MQTT client
	mqtt, err := NewMQTTClient(config)
	if err != nil {
		log.Printf("Failed to create MQTT client: %v", err)
		// Continue without MQTT for now - allow offline mode
	} else {
		app.mqtt = mqtt
	}

	return app
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	if a.mqtt != nil {
		return tea.Batch(
			a.ui.Init(),
			a.mqtt.ConnectCmd(),
		)
	}
	return a.ui.Init()
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Pass window size to UI first
		var uiCmd tea.Cmd
		a.ui, uiCmd = a.ui.Update(msg)
		if uiCmd != nil {
			cmds = append(cmds, uiCmd)
		}
		return a, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.quitting = true
			if a.mqtt != nil {
				a.mqtt.Disconnect()
			}
			return a, tea.Quit
		}
	case MQTTConnectedMsg:
		// Start topic discovery when connected
		if a.mqtt != nil {
			cmds = append(cmds, a.mqtt.DiscoverTopicsCmd())
		}
	case MQTTTopicsDiscoveredMsg:
		// Update UI with discovered topics
		a.ui.SetTopics(msg.Topics)
	case MQTTMessageMsg:
		// Update UI with new message
		a.ui.AddMessage(msg.Topic, msg.Payload, msg.Timestamp)
	case MQTTErrorMsg:
		// Handle MQTT errors
		a.ui.SetError(fmt.Sprintf("MQTT Error: %v", msg.Error))
	}

	// Update UI and handle subscription changes
	oldSubscribed := a.ui.GetSubscribedTopics()
	var uiCmd tea.Cmd
	a.ui, uiCmd = a.ui.Update(msg)
	if uiCmd != nil {
		cmds = append(cmds, uiCmd)
	}

	// Check for subscription changes
	if a.mqtt != nil && a.mqtt.IsConnected() {
		newSubscribed := a.ui.GetSubscribedTopics()
		cmds = append(cmds, a.handleSubscriptionChanges(oldSubscribed, newSubscribed)...)
	}

	return a, tea.Batch(cmds...)
}

// View implements tea.Model
func (a *App) View() string {
	if a.quitting {
		return "\nDisconnecting from MQTT broker...\nGoodbye!\n"
	}
	return a.ui.View()
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// handleSubscriptionChanges handles topic subscription/unsubscription
func (a *App) handleSubscriptionChanges(oldSubscribed, newSubscribed []string) []tea.Cmd {
	var cmds []tea.Cmd

	// Create maps for easier comparison
	oldMap := make(map[string]bool)
	for _, topic := range oldSubscribed {
		oldMap[topic] = true
	}

	newMap := make(map[string]bool)
	for _, topic := range newSubscribed {
		newMap[topic] = true
	}

	// Subscribe to new topics
	for topic := range newMap {
		if !oldMap[topic] {
			cmds = append(cmds, a.subscribeToTopicCmd(topic))
		}
	}

	// Unsubscribe from removed topics
	for topic := range oldMap {
		if !newMap[topic] {
			cmds = append(cmds, a.unsubscribeFromTopicCmd(topic))
		}
	}

	return cmds
}

// subscribeToTopicCmd creates a command to subscribe to a topic
func (a *App) subscribeToTopicCmd(topic string) tea.Cmd {
	return func() tea.Msg {
		if err := a.mqtt.SubscribeToTopic(topic); err != nil {
			return MQTTErrorMsg{Error: fmt.Errorf("failed to subscribe to %s: %v", topic, err)}
		}
		return nil
	}
}

// unsubscribeFromTopicCmd creates a command to unsubscribe from a topic
func (a *App) unsubscribeFromTopicCmd(topic string) tea.Cmd {
	return func() tea.Msg {
		if err := a.mqtt.UnsubscribeFromTopic(topic); err != nil {
			return MQTTErrorMsg{Error: fmt.Errorf("failed to unsubscribe from %s: %v", topic, err)}
		}
		return nil
	}
}
