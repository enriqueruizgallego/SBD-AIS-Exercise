package main

import (
	"exc8/client"
	"exc8/server"
	"fmt"
	"time"
)

func main() {
	go func() {
		// todo start server
		server.StartGrpcServer()
	}()
	time.Sleep(1 * time.Second)
	// todo start client
	cli, err := client.NewGrpcClient()
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	err = cli.Run()
	if err != nil {
		fmt.Println("Client error:", err)
		return
	}

	println("Orders complete!")
}
