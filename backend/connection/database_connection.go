package connection

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func GetDatabaseConnection() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file for database connection")
	}

	database_user := os.Getenv("DATABASE_USER")
	database_password := os.Getenv("DATABASE_PASSWORD")
	database_url := os.Getenv("DATABASE_URL")

	dsn := "postgresql://" + database_user + ":" + database_password + database_url + "?sslmode=verify-full"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	return conn, nil
}

func CloseDatabaseConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func TestDatabaseConnection() {
	conn, err := GetDatabaseConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer CloseDatabaseConnection(conn)

	ctx := context.Background()
	var now time.Time
	err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("Failed to execute test query:", err)
	}
	fmt.Println("Database connected successfully. Current time:", now)
}
