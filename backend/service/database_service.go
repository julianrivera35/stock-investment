package service

import (
	"context"
	"fmt"
	"log"
	"stock-investment-backend/connection"
	"time"

	"github.com/jackc/pgx/v4"
)

type RecommendationData struct {
	Ticker     string    `json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

// Ctx is context.Context that handles request timeouts, cancellations and deadlines. Prevents queries from hanging indefinitely. Allows to stop operations if needed. Required for database operations.
func InsertOrGetCompany(conn *pgx.Conn, ctx context.Context, ticker, name string) (string, error) {
	var companyID string

	err := conn.QueryRow(ctx,
		"SELECT id FROM company WHERE ticker = $1",
		ticker).Scan(&companyID)

	if err == pgx.ErrNoRows {
		err = conn.QueryRow(ctx,
			`INSERT INTO company (ticker, name)
			VALUES ($1, $2)
			ON CONFLICT (ticker) DO UPDATE SET
			name = EXCLUDED.name, updated_at = now()
			RETURNING id`,
			ticker, name).Scan(&companyID)

		if err != nil {
			return "", fmt.Errorf("failed to query company: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to query company: %w", err)
	}

	return companyID, nil
}

func InsertOrGetBrokerage(conn *pgx.Conn, ctx context.Context, name string) (string, error) {
	var brokerageID string

	err := conn.QueryRow(ctx,
		"SELECT id FROM brokerage WHERE name = $1",
		name).Scan(&brokerageID)

	if err == pgx.ErrNoRows {
		err = conn.QueryRow(ctx,
			`INSERT INTO brokerage (name)
				VALUES ($1)
				ON CONFLICT (name) DO UPDATE SET
				updated_at = now()
				RETURNING id`,
			name).Scan(&brokerageID)

		if err != nil {
			return "", fmt.Errorf("failed to query brokerage: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to query brokerage: %w", err)
	}

	return brokerageID, nil
}

func InsertRecommendation(conn *pgx.Conn, ctx context.Context, data RecommendationData) error {
	//Get or create company
	companyID, err := InsertOrGetCompany(conn, ctx, data.Ticker, data.Company)
	if err != nil {
		return fmt.Errorf("failed to insert or get company: %v", err)
	}

	//Get or create brokerage
	brokerageID, err := InsertOrGetBrokerage(conn, ctx, data.Brokerage)
	if err != nil {
		return fmt.Errorf("failed to insert or get brokerage: %v", err)
	}

	//Parse target prices and handle empty strings
	var targetFrom, targetTo *float64
	if data.TargetFrom != "" && data.TargetFrom != "$0.00" {
		var tf float64
		_, err := fmt.Sscanf(data.TargetFrom, "$%f", &tf)
		if err == nil {
			targetFrom = &tf
		}
	}

	if data.TargetTo != "" && data.TargetTo != "$0.00" {
		var tt float64
		_, err := fmt.Sscanf(data.TargetTo, "$%f", &tt)
		if err == nil {
			targetTo = &tt
		}
	}

	// Insert analyst recommendation
	_, err = conn.Exec(ctx,
		`INSERT INTO analyst_recommendation 
		(company_id, brokerage_id, target_from, target_to, rating_from, rating_to, action, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		companyID, brokerageID, targetFrom, targetTo, data.RatingFrom, data.RatingTo, data.Action, data.Time)
	if err != nil {
		return fmt.Errorf("failed to insert analyst recommendation: %v", err)
	}

	return nil
}

// SaveRecommendations saves multiple recommendations to the database
func SaveRecommendations(recommendations []RecommendationData) error {
	conn, err := connection.GetDatabaseConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}
	defer conn.Close(context.Background())

	ctx := context.Background()

	//Start transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	successCount := 0
	for i, rec := range recommendations {
		err := InsertRecommendation(conn, ctx, rec)
		if err != nil {
			log.Printf("Error inserting recommendation %d: %v", i, err)
			continue
		}
		successCount++
	}

	//Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully inserted %d recommendations", successCount)
	return nil
}
