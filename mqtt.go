package main

import (
	"log"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT Message types for Bubble Tea
type MQTTConnectedMsg struct{}
type MQTTDisconnectedMsg struct{}
type MQTTTopicsDiscoveredMsg struct {
	Topics []string
}
type MQTTMessageMsg struct {
	Topic     string
	Payload   string
	Timestamp time.Time
}
type MQTTErrorMsg struct {
	Error error
}

// MQTTClient wraps the MQTT functionality
type MQTTClient struct {
	client           mqtt.Client
	config           Config
	discoveredTopics map[string]bool
	topicsMutex      sync.RWMutex
	program          *tea.Program
}

// NewMQTTClient creates a new MQTT client
func NewMQTTClient(config Config) (*MQTTClient, error) {
	client := &MQTTClient{
		config:           config,
		discoveredTopics: make(map[string]bool),
	}

	// Set up MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.BrokerURL)
	opts.SetClientID(config.ClientID)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	// Set connection handlers
	opts.SetDefaultPublishHandler(client.messageHandler)
	opts.SetOnConnectHandler(client.connectHandler)
	opts.SetConnectionLostHandler(client.connectionLostHandler)

	// Create the MQTT client
	client.client = mqtt.NewClient(opts)

	return client, nil
}

// SetProgram sets the Bubble Tea program for sending messages
func (m *MQTTClient) SetProgram(p *tea.Program) {
	m.program = p
}

// ConnectCmd returns a command to connect to the MQTT broker
func (m *MQTTClient) ConnectCmd() tea.Cmd {
	return func() tea.Msg {
		if token := m.client.Connect(); token.Wait() && token.Error() != nil {
			return MQTTErrorMsg{Error: token.Error()}
		}
		return MQTTConnectedMsg{}
	}
}

// DiscoverTopicsCmd subscribes to # wildcard to discover all topics
func (m *MQTTClient) DiscoverTopicsCmd() tea.Cmd {
	return func() tea.Msg {
		// Subscribe to all topics to discover them
		if token := m.client.Subscribe("#", 0, m.discoveryHandler); token.Wait() && token.Error() != nil {
			return MQTTErrorMsg{Error: token.Error()}
		}

		// Wait a bit to collect topics, then return discovered topics
		time.Sleep(2 * time.Second)

		m.topicsMutex.RLock()
		topics := make([]string, 0, len(m.discoveredTopics))
		for topic := range m.discoveredTopics {
			topics = append(topics, topic)
		}
		m.topicsMutex.RUnlock()

		return MQTTTopicsDiscoveredMsg{Topics: topics}
	}
}

// SubscribeToTopic subscribes to a specific topic
func (m *MQTTClient) SubscribeToTopic(topic string) error {
	if token := m.client.Subscribe(topic, 0, m.messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// UnsubscribeFromTopic unsubscribes from a specific topic
func (m *MQTTClient) UnsubscribeFromTopic(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect disconnects from the MQTT broker
func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}

// IsConnected returns true if connected to the broker
func (m *MQTTClient) IsConnected() bool {
	return m.client.IsConnected()
}

// Message handlers
func (m *MQTTClient) connectHandler(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
	if m.program != nil {
		m.program.Send(MQTTConnectedMsg{})
	}
}

func (m *MQTTClient) connectionLostHandler(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
	if m.program != nil {
		m.program.Send(MQTTErrorMsg{Error: err})
	}
}

func (m *MQTTClient) messageHandler(client mqtt.Client, msg mqtt.Message) {
	if m.program != nil {
		m.program.Send(MQTTMessageMsg{
			Topic:     msg.Topic(),
			Payload:   string(msg.Payload()),
			Timestamp: time.Now(),
		})
	}
}

func (m *MQTTClient) discoveryHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()

	m.topicsMutex.Lock()
	m.discoveredTopics[topic] = true
	m.topicsMutex.Unlock()

	// Also handle as regular message
	m.messageHandler(client, msg)
}

// GetDiscoveredTopics returns a list of discovered topics
func (m *MQTTClient) GetDiscoveredTopics() []string {
	m.topicsMutex.RLock()
	defer m.topicsMutex.RUnlock()

	topics := make([]string, 0, len(m.discoveredTopics))
	for topic := range m.discoveredTopics {
		topics = append(topics, topic)
	}
	return topics
}
