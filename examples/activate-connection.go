package main

import (
	"fmt"
	"os"

	networkmanager "github.com/phoracek/networkmanager-go/src"
)

func main() {
	connectionIDToRemove := os.Args[1]

	client := networkmanager.NewClient()
	defer client.Close()

	if connection := findConnection(client, connectionIDToRemove); connection != nil {
		err := client.ActivateConnection(connection)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed", err)
			os.Exit(1)
		}
	}
}

func findConnection(client *networkmanager.Client, connectionID string) *networkmanager.Connection {
	connections := client.ListConnections()
	for _, connection := range connections {
		settings, _ := connection.GetSettings()
		if settings["connection"]["id"] == connectionID {
			return connection
		}
	}
	return nil
}
