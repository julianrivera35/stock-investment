package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"github.com/jackc/pgconn" 
)

type ConfigLoader interface{
	Load() error
	GetEnv(key string) string
}

type DBConnector interface{
	Connect(ctx context.Context, dsn string) (DBConnection, error)
}

type DBConnection interface {
	BeginConn(ctx context.Context) (pgx.Tx, error)
	CloseConn(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec (ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type EnvConfigLoader struct{}

func (e *EnvConfigLoader) Load() error {
	return godotenv.Load()
}

func (e *EnvConfigLoader) GetEnv(key string) string {
	return os.Getenv(key)
}

type PgxConnector struct{}

func (p *PgxConnector) Connect(ctx context.Context, dsn string) (DBConnection, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &PgxConnection{conn: conn}, nil
}

type PgxConnection struct{
	conn *pgx.Conn
}

func (p *PgxConnection) BeginConn(ctx context.Context) (pgx.Tx, error) {
	return p.conn.Begin(ctx)
}
func (p *PgxConnection) CloseConn(ctx context.Context) error {
	return p.conn.Close(ctx)
}

func (p *PgxConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return p.conn.QueryRow(ctx, sql, args...)
}

func (p *PgxConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	rows, err := p.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *PgxConnection) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return p.conn.Exec(ctx, sql, args...)
}

//Global variables for dependency injection
var(
	configLoader ConfigLoader = &EnvConfigLoader{}
	dbConnector DBConnector = &PgxConnector{}
)

func BuildDSN(user, password, url, sslMode string) string {
	return fmt.Sprintf("postgresql://%s:%s%s?sslmode=%s", user, password, url, sslMode)
}

func GetDatabaseConnection() (DBConnection, error) {
	err := configLoader.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w",err)
	}

	database_user := configLoader.GetEnv("DATABASE_USER")
	database_password := configLoader.GetEnv("DATABASE_PASSWORD")
	database_url := configLoader.GetEnv("DATABASE_URL")
	ssl_mode_local := configLoader.GetEnv("SSL_LOCAL");

	dsn := BuildDSN(database_user, database_password, database_url, ssl_mode_local)
	ctx := context.Background()
	conn, err := dbConnector.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return conn, nil
}

func CloseDatabaseConnection(conn DBConnection) error {
	if conn == nil {
		return nil
	}
	return conn.CloseConn(context.Background())
}

func TestDatabaseConnection() error {
	conn, err := GetDatabaseConnection()
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer CloseDatabaseConnection(conn)

	ctx := context.Background()
	var now time.Time
	err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}
	fmt.Println("Database connected successfully. Current time:", now)
	return nil
}
