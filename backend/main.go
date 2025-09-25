package main

import (
	"fmt"
	"stock-investment-backend/connection"
	"stock-investment-backend/service"
)

func main() {
	fmt.Println("Starting Stock Investment Backend...")
	
	fmt.Println("Testing database connection...")
	connection.TestDatabaseConnection()

	fmt.Println("Fetching data from API...")
	service.ApiGet()
}
