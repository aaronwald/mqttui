#!/bin/bash

# MQTT TUI - Example startup script
# This script demonstrates how to run the MQTT TUI with different configurations

echo "MQTT TUI Browser - Startup Examples"
echo "=================================="
echo

# Example 1: Default local broker
echo "1. Running with default localhost broker:"
echo "   MQTT_BROKER='tcp://localhost:1883'"
echo "   ./mqttui"
echo

# Example 2: Public Eclipse test broker
echo "2. Running with Eclipse public broker:"
echo "   MQTT_BROKER='tcp://mqtt.eclipseprojects.io:1883' ./mqttui"
echo

# Example 3: HiveMQ public broker
echo "3. Running with HiveMQ public broker:"
echo "   MQTT_BROKER='tcp://broker.hivemq.com:1883' ./mqttui"
echo

# Example 4: With authentication
echo "4. Running with authentication:"
echo "   MQTT_BROKER='tcp://your-broker.com:1883' \\"
echo "   MQTT_USERNAME='your_user' \\"
echo "   MQTT_PASSWORD='your_pass' \\"
echo "   ./mqttui"
echo

echo "Choose an option (1-4) or press Ctrl+C to exit:"
read -p "Option: " choice

case $choice in
    1)
        echo "Starting with localhost broker..."
        MQTT_BROKER="tcp://localhost:1883" ./mqttui
        ;;
    2)
        echo "Starting with Eclipse public broker..."
        MQTT_BROKER="tcp://mqtt.eclipseprojects.io:1883" ./mqttui
        ;;
    3)
        echo "Starting with HiveMQ public broker..."
        MQTT_BROKER="tcp://broker.hivemq.com:1883" ./mqttui
        ;;
    4)
        echo "Enter your broker details:"
        read -p "Broker URL (e.g., tcp://broker.com:1883): " broker
        read -p "Username: " username
        read -p "Password: " -s password
        echo
        echo "Starting with custom broker..."
        MQTT_BROKER="$broker" MQTT_USERNAME="$username" MQTT_PASSWORD="$password" ./mqttui
        ;;
    *)
        echo "Invalid option. Exiting."
        exit 1
        ;;
esac