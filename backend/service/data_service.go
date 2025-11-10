package service

import (
	"context"
	"fmt"
	"stock-investment-backend/connection"
	"time"
)

type Company struct {
	ID        string    `json:"id"`
	Ticker    string    `json:"ticker"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Brokerage struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Recommendation struct {
	ID         string     `json:"id"`
	Company    Company    `json:"company"`
	Brokerage  *Brokerage `json:"brokerage,omitempty"`
	TargetFrom *float64   `json:"target_from"`
	TargetTo   *float64   `json:"target_to"`
	RatingFrom string     `json:"rating_from"`
	RatingTo   string     `json:"rating_to"`
	Action     string     `json:"action"`
	Time       time.Time  `json:"time"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Retrieve all companies
func GetAllCompanies() ([]Company, error) {
	conn, err := connection.GetDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.CloseConn(context.Background())

	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT id, ticker, name, created_at, updated_at
		FROM company
		ORDER BY ticker ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		//scan assigns values from each column in the row to the corresponding struct fields
		err := rows.Scan(&c.ID, &c.Ticker, &c.Name, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		companies = append(companies, c)
	}

	return companies, nil
}

// Get company by ticker
func GetCompanyByTicker(ticker string) (*Company, error) {
	conn, err := connection.GetDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.CloseConn(context.Background())

	ctx := context.Background()
	var c Company
	err = conn.QueryRow(ctx, `
		SELECT id, ticker, name, created_at, updated_at
		FROM company
		WHERE ticker = $1`,
		ticker).Scan(&c.ID, &c.Ticker, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("company not found: %v", err)
	}
	return &c, nil
}

// Retrieve all brokerages
func GetAllBrokerages() ([]Brokerage, error) {
	conn, err := connection.GetDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.CloseConn(context.Background())

	ctx := context.Background()
	rows, err := conn.Query(ctx, `
		SELECT id, name, created_at, updated_at
		FROM brokerage
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var brokerages []Brokerage
	for rows.Next() {
		var b Brokerage
		err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		brokerages = append(brokerages, b)
	}

	return brokerages, nil
}

// Retrieve recommendations
func GetRecommendations(limit, offset int, ticker, brokerageID string) ([]Recommendation, int, error) {
	conn, err := connection.GetDatabaseConnection()
	if err != nil {
		return nil, 0, fmt.Errorf("database connection failed: %v", err)
	}
	defer conn.CloseConn(context.Background())

	ctx := context.Background()

	//Filtered query
	baseQuery := `
		SELECT
			ar.id, ar.target_from, ar.target_to, ar.rating_from, ar.rating_to,
			ar.action, ar.time, ar.created_at, ar.updated_at,
			c.id, c.ticker, c.name, c.created_at, c.updated_at,
			b.id, b.name, b.created_at, b.updated_at
		FROM analyst_recommendation ar
		JOIN company c ON ar.company_id = c.id
		LEFT JOIN brokerage b ON ar.brokerage_id = b.id
	`

	countQuery := `
		SELECT COUNT(*)
		FROM analyst_recommendation ar
		JOIN company c ON ar.company_id = c.id
		LEFT JOIN brokerage b ON ar.brokerage_id = b.id
	`

	whereClause := ""
	args := []interface{}{}
	argsIndex := 1

	if ticker != "" {
		whereClause += fmt.Sprintf(" WHERE c.ticker = $%d", argsIndex)
		args = append(args, ticker)
		argsIndex++
	}

	if brokerageID != "" {
		if whereClause == "" {
			whereClause += fmt.Sprintf(" WHERE ar.brokerage_id = $%d", argsIndex)
		} else {
			whereClause += fmt.Sprintf(" AND ar.brokerage_id = $%d", argsIndex)
		}
		args = append(args, brokerageID)
		argsIndex++
	}

	var totalCount int
	err = conn.QueryRow(ctx, countQuery+whereClause, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %v", err)
	}

	//Get recommendations
	finalQuery := baseQuery + whereClause + fmt.Sprintf(" ORDER BY ar.time DESC LIMIT $%d OFFSET $%d", argsIndex, argsIndex+1)
	args = append(args, limit, offset)

	rows, err := conn.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %v", err)
	}
	defer rows.Close()

	var recommendations []Recommendation
	for rows.Next() {
		var r Recommendation
		var dbBrokerageID, dbBrokerageName *string
		var dbBrokerageCreatedAt, dbBrokerageUpdatedAt *time.Time

		err := rows.Scan(
			&r.ID, &r.TargetFrom, &r.TargetTo, &r.RatingFrom, &r.RatingTo,
			&r.Action, &r.Time, &r.CreatedAt, &r.UpdatedAt,
			&r.Company.ID, &r.Company.Ticker, &r.Company.Name, &r.Company.CreatedAt, &r.Company.UpdatedAt,
			&dbBrokerageID, &dbBrokerageName, &dbBrokerageCreatedAt, &dbBrokerageUpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %v", err)
		}

		//Handle multiple brokerages
		if dbBrokerageID != nil {
			r.Brokerage = &Brokerage{
				ID:        *dbBrokerageID,
				Name:      *dbBrokerageName,
				CreatedAt: *dbBrokerageCreatedAt,
				UpdatedAt: *dbBrokerageUpdatedAt,
			}
		}
		recommendations = append(recommendations, r)
	}
	return recommendations, totalCount, nil
}
