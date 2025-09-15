package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UI represents the user interface state
type UI struct {
	topics           []string
	selectedTopic    int
	subscribedTopics map[string]bool
	messages         []Message
	messageScroll    int
	width            int
	height           int
	activePane       Pane
	error            string
	styles           Styles
}

// Pane represents which pane is currently active
type Pane int

const (
	TopicsPane Pane = iota
	MessagesPane
)

// Message represents an MQTT message
type Message struct {
	Topic     string
	Payload   string
	Timestamp time.Time
}

// Styles holds all the styling for the UI
type Styles struct {
	Border         lipgloss.Style
	Title          lipgloss.Style
	SelectedItem   lipgloss.Style
	UnselectedItem lipgloss.Style
	Message        lipgloss.Style
	MessageTopic   lipgloss.Style
	MessageTime    lipgloss.Style
	Error          lipgloss.Style
	Help           lipgloss.Style
	ActivePane     lipgloss.Style
	InactivePane   lipgloss.Style
}

// NewUI creates a new UI instance
func NewUI() *UI {
	styles := Styles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1),
		SelectedItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("57")).
			Bold(true),
		UnselectedItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
		Message: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 0, 1, 0),
		MessageTopic: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true),
		MessageTime: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true),
		ActivePane: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")),
		InactivePane: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")),
	}

	return &UI{
		topics:           []string{},
		subscribedTopics: make(map[string]bool),
		messages:         []Message{},
		activePane:       TopicsPane,
		styles:           styles,
	}
}

// Init implements tea.Model
func (ui *UI) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (ui *UI) Update(msg tea.Msg) (*UI, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ui.width = msg.Width
		ui.height = msg.Height
	case tea.KeyMsg:
		return ui.handleKeyPress(msg)
	}
	return ui, nil
}

// handleKeyPress handles keyboard input
func (ui *UI) handleKeyPress(msg tea.KeyMsg) (*UI, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Switch between panes
		if ui.activePane == TopicsPane {
			ui.activePane = MessagesPane
		} else {
			ui.activePane = TopicsPane
		}
	case "up", "k":
		if ui.activePane == TopicsPane {
			if ui.selectedTopic > 0 {
				ui.selectedTopic--
			}
		} else {
			if ui.messageScroll > 0 {
				ui.messageScroll--
			}
		}
	case "down", "j":
		if ui.activePane == TopicsPane {
			if ui.selectedTopic < len(ui.topics)-1 {
				ui.selectedTopic++
			}
		} else {
			ui.messageScroll++
		}
	case "enter", " ":
		if ui.activePane == TopicsPane && len(ui.topics) > 0 {
			topic := ui.topics[ui.selectedTopic]
			ui.subscribedTopics[topic] = !ui.subscribedTopics[topic]
		}
	case "r":
		// Reset messages
		ui.messages = []Message{}
		ui.messageScroll = 0
	}
	return ui, nil
}

// View implements tea.Model
func (ui *UI) View() string {
	if ui.width == 0 {
		return "Initializing..."
	}

	// Calculate dimensions
	topicsWidth := ui.width / 3
	messagesWidth := ui.width - topicsWidth - 2

	// Create the topics view
	topicsView := ui.renderTopicsPane(topicsWidth, ui.height-4)

	// Create the messages view
	messagesView := ui.renderMessagesPane(messagesWidth, ui.height-4)

	// Combine the views
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		topicsView,
		messagesView,
	)

	// Add title and help
	title := ui.styles.Title.Render("MQTT TUI Browser")
	help := ui.renderHelp()

	// Add error if present
	if ui.error != "" {
		errorMsg := ui.styles.Error.Render(fmt.Sprintf("Error: %s", ui.error))
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			content,
			errorMsg,
			help,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		help,
	)
}

// renderTopicsPane renders the topics list pane
func (ui *UI) renderTopicsPane(width, height int) string {
	title := "Topics"
	if len(ui.topics) > 0 {
		title += fmt.Sprintf(" (%d)", len(ui.topics))
	}

	var items []string
	for i, topic := range ui.topics {
		prefix := "  "
		if ui.subscribedTopics[topic] {
			prefix = "✓ "
		}

		item := prefix + topic
		if i == ui.selectedTopic && ui.activePane == TopicsPane {
			item = ui.styles.SelectedItem.Render(item)
		} else {
			item = ui.styles.UnselectedItem.Render(item)
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		items = append(items, ui.styles.UnselectedItem.Render("No topics discovered yet..."))
	}

	content := strings.Join(items, "\n")

	style := ui.styles.InactivePane
	if ui.activePane == TopicsPane {
		style = ui.styles.ActivePane
	}

	return style.
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			ui.styles.Title.Render(title),
			content,
		))
}

// renderMessagesPane renders the messages pane
func (ui *UI) renderMessagesPane(width, height int) string {
	title := "Messages"
	if len(ui.messages) > 0 {
		title += fmt.Sprintf(" (%d)", len(ui.messages))
	}

	var items []string
	startIdx := ui.messageScroll
	endIdx := startIdx + height - 3 // Account for title and padding

	if endIdx > len(ui.messages) {
		endIdx = len(ui.messages)
	}
	if startIdx >= endIdx {
		startIdx = endIdx - 1
		if startIdx < 0 {
			startIdx = 0
		}
	}

	for i := startIdx; i < endIdx; i++ {
		msg := ui.messages[i]
		timeStr := msg.Timestamp.Format("15:04:05")

		topicLine := ui.styles.MessageTopic.Render(msg.Topic) +
			" " + ui.styles.MessageTime.Render(timeStr)

		// Wrap payload text
		payloadLines := ui.wrapText(msg.Payload, width-4)

		messageContent := lipgloss.JoinVertical(
			lipgloss.Left,
			topicLine,
			strings.Join(payloadLines, "\n"),
		)

		items = append(items, ui.styles.Message.Render(messageContent))
	}

	if len(items) == 0 {
		items = append(items, ui.styles.UnselectedItem.Render("No messages yet..."))
	}

	content := strings.Join(items, "\n")

	style := ui.styles.InactivePane
	if ui.activePane == MessagesPane {
		style = ui.styles.ActivePane
	}

	return style.
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			ui.styles.Title.Render(title),
			content,
		))
}

// renderHelp renders the help text
func (ui *UI) renderHelp() string {
	help := "↑/↓ navigate • tab switch panes • enter/space toggle subscription • r reset messages • q quit"
	return ui.styles.Help.Render(help)
}

// wrapText wraps text to fit within the specified width
func (ui *UI) wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		if currentLength+len(word)+len(currentLine) > width {
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
				currentLine = []string{word}
				currentLength = len(word)
			} else {
				// Word is too long, split it
				lines = append(lines, word[:width])
				currentLine = []string{}
				currentLength = 0
			}
		} else {
			currentLine = append(currentLine, word)
			currentLength += len(word)
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return lines
}

// SetTopics updates the list of available topics
func (ui *UI) SetTopics(topics []string) {
	sort.Strings(topics)
	ui.topics = topics
	if ui.selectedTopic >= len(topics) {
		ui.selectedTopic = len(topics) - 1
	}
	if ui.selectedTopic < 0 {
		ui.selectedTopic = 0
	}
}

// AddMessage adds a new message to the messages list
func (ui *UI) AddMessage(topic, payload string, timestamp time.Time) {
	message := Message{
		Topic:     topic,
		Payload:   payload,
		Timestamp: timestamp,
	}

	ui.messages = append(ui.messages, message)

	// Auto-scroll to bottom when new message arrives
	if ui.activePane == MessagesPane {
		ui.messageScroll = len(ui.messages) - 1
	}
}

// SetError sets an error message
func (ui *UI) SetError(err string) {
	ui.error = err
}

// GetSubscribedTopics returns the list of subscribed topics
func (ui *UI) GetSubscribedTopics() []string {
	var subscribed []string
	for topic, isSubscribed := range ui.subscribedTopics {
		if isSubscribed {
			subscribed = append(subscribed, topic)
		}
	}
	return subscribed
}
