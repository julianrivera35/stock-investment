package main

import (
	"flag"
	"fmt"
	"log"
	"stock-investment-backend/connection"
	"stock-investment-backend/server"
	"stock-investment-backend/service"
)

func main() {
	var (
		apiMode   = flag.Bool("api", false, "Run in API server mode")
		fetchMode = flag.Bool("fetch", false, "Fetch data from external API")
		port      = flag.String("port", "8080", "Port for API server")
	)
	flag.Parse()

	fmt.Println("Starting Stock Investment Backend...")

	// Test database connection
	fmt.Println("Testing database connection...")
	if err := connection.TestDatabaseConnection(); err != nil {
		log.Fatal(err)
	}

	if *fetchMode {
		// Fetch API data and save to database
		fmt.Println("Fetching API data...")
		service.ApiGet()
	} else if *apiMode {
		// Start API server
		srv := server.NewServer()
		srv.Start(*port)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  -fetch    Fetch data from external API")
		fmt.Println("  -api      Start API server")
		fmt.Println("  -port     Port for API server (default: 8080)")
	}
}
