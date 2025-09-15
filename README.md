# MQTT TUI Browser

A terminal user interface (TUI) application built with Go and the Charm framework for browsing MQTT topics and viewing real-time message updates.

## Features

- **Real-time Topic Discovery**: Automatically discovers all available topics on the MQTT broker
- **Interactive Topic Browser**: Navigate through topics with keyboard controls
- **Topic Subscription Management**: Subscribe/unsubscribe to topics with space or enter
- **Live Message Display**: View real-time messages from subscribed topics
- **Dual-pane Interface**: Split view with topics on the left and messages on the right
- **Keyboard Navigation**: Full keyboard control with intuitive shortcuts
- **Styled Interface**: Beautiful TUI with colors and styling using Lip Gloss

## Installation

1. Clone this repository
2. Make sure you have Go 1.19+ installed
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the application:
   ```bash
   go build -o mqttui
   ```

## Usage

### Environment Variables

Configure the MQTT connection using environment variables:

```bash
export MQTT_BROKER="tcp://localhost:1883"    # MQTT broker URL
export MQTT_USERNAME="your_username"         # Optional: MQTT username
export MQTT_PASSWORD="your_password"         # Optional: MQTT password
export MQTT_CLIENT_ID="mqttui"              # Optional: MQTT client ID
```

### Running the Application

```bash
# Using default settings (localhost:1883)
./mqttui

# Or with custom broker
MQTT_BROKER="tcp://broker.example.com:1883" ./mqttui

# With authentication
MQTT_BROKER="tcp://broker.example.com:1883" \
MQTT_USERNAME="user" \
MQTT_PASSWORD="pass" \
./mqttui
```

### Keyboard Controls

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate up/down in the active pane |
| `Tab` | Switch between topics and messages panes |
| `Enter` or `Space` | Subscribe/unsubscribe to selected topic |
| `r` | Reset/clear all messages |
| `q` or `Ctrl+C` | Quit the application |

### Interface Layout

```
┌─ MQTT TUI Browser ─────────────────────────────────────────┐
│                                                            │
│ ┌─ Topics (5) ─────────┐ ┌─ Messages (12) ──────────────┐ │
│ │ ✓ home/temperature   │ │ home/temperature 14:30:15    │ │
│ │   home/humidity      │ │ 23.5                         │ │
│ │ ✓ sensors/motion     │ │                              │ │
│ │   devices/switch1    │ │ sensors/motion 14:30:20      │ │
│ │   system/status      │ │ true                         │ │
│ └─────────────────────┘ └─────────────────────────────┘ │
│                                                            │
│ ↑/↓ navigate • tab switch panes • enter/space toggle      │
│ subscription • r reset messages • q quit                  │
└────────────────────────────────────────────────────────────┘
```

- **Left Pane**: Shows all discovered topics. Subscribed topics are marked with ✓
- **Right Pane**: Shows real-time messages from subscribed topics
- **Active Pane**: Highlighted with colored border
- **Status**: Help text at the bottom shows available keyboard shortcuts

## Architecture

The application is built using the following components:

### Core Components

1. **main.go**: Application entry point and coordination
2. **mqtt.go**: MQTT client functionality and message handling
3. **ui.go**: Terminal user interface using Bubble Tea and Lip Gloss

### Key Features

- **Topic Discovery**: Uses wildcard subscription (`#`) to discover all topics
- **Real-time Updates**: Asynchronous message handling with Bubble Tea commands
- **State Management**: Clean separation between MQTT logic and UI state
- **Error Handling**: Graceful error display and connection management

### Dependencies

- **Bubble Tea**: TUI framework for interactive terminal applications
- **Lip Gloss**: Styling and layout for beautiful terminal interfaces
- **Eclipse Paho**: MQTT client library for Go

## Development

### Project Structure

```
mqttui/
├── main.go          # Application entry point
├── mqtt.go          # MQTT client implementation
├── ui.go           # Terminal user interface
├── go.mod          # Go module dependencies
├── go.sum          # Dependency checksums
└── README.md       # This file
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o mqttui-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o mqttui-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o mqttui.exe
```

## Testing with Public Brokers

You can test the application with public MQTT brokers:

```bash
# Eclipse test broker
MQTT_BROKER="tcp://mqtt.eclipseprojects.io:1883" ./mqttui

# HiveMQ public broker
MQTT_BROKER="tcp://broker.hivemq.com:1883" ./mqttui

# Mosquitto test broker
MQTT_BROKER="tcp://test.mosquitto.org:1883" ./mqttui
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the excellent TUI frameworks
- [Eclipse Paho](https://www.eclipse.org/paho/) for the MQTT client library
- The Go community for the amazing ecosystem