package connection

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/jackc/pgconn" 
)

// MockConfigLoader implemenets ConfigLoader interface for testing
type MockConfigLoader struct {
	LoadFunc   func() error
	GetEnvFunc func(key string) string
}

func (m *MockConfigLoader) Load() error {
	if m.LoadFunc != nil {
		return m.LoadFunc()
	}
	return nil
}

func (m *MockConfigLoader) GetEnv(key string) string {
	if m.GetEnvFunc != nil {
		return m.GetEnvFunc(key)
	}
	return ""
}

// MockDBConnector implemenets DBConnector interface for testing
type MockDBConnector struct {
	ConnectFunc func(ctx context.Context, dsn string) (DBConnection, error)
}

func (m *MockDBConnector) Connect(ctx context.Context, dsn string) (DBConnection, error) {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(ctx, dsn)
	}
	return nil, nil
}

// Mock DBConnection implemenets DBConnection interface for testing
type MockDBConnection struct {
	CloseCalled  bool
	CloseError   error
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

func (m *MockDBConnection) BeginConn(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}
func (m *MockDBConnection) CloseConn(ctx context.Context) error {
	m.CloseCalled = true
	return m.CloseError
}

func (m *MockDBConnection) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return nil
}

func (m *MockDBConnection) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return nil, nil
}

func (m *MockDBConnection) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return nil, nil
}

//Mock Row
type MockRow struct {
	ScanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error{
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil
}

// Setup and teardown
func setupTest() func() {
	originalConfigLoader := configLoader
	originalDBConnector := dbConnector

	return func() {
		configLoader = originalConfigLoader
		dbConnector = originalDBConnector
	}
}

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		url      string
		sslMode  string
		expected string
	}{
		{
			name:     "Whole DSN with SSL disabled",
			user:     "testuser",
			password: "testpass",
			url:      "@localhost:5432/testdb",
			sslMode:  "disable",
			expected: "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name:     "Whole DSN with SSL enabled",
			user:     "testuser",
			password: "testpass",
			url:      "@localhost:5432/testdb",
			sslMode:  "verify-full",
			expected: "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=verify-full",
		},
		{
			name:     "Whole DSN without values",
			user:     "",
			password: "",
			url:      "",
			sslMode:  "",
			expected: "postgresql://:?sslmode=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildDSN(tt.user, tt.password, tt.url, tt.sslMode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetDatabaseConnection_Success(t *testing.T) {
	teardown := setupTest()
	defer teardown()

	mockConfig := &MockConfigLoader{
		LoadFunc: func() error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			envVars := map[string]string{
				"DATABASE_USER":     "testuser",
				"DATABASE_PASSWORD": "testpass",
				"DATABASE_URL":      "@localhost:5432/testdb",
				"SSL_LOCAL":         "disable",
			}
			return envVars[key]
		},
	}
	dsnWasCorrect := false
	mockConn := &MockDBConnection{}

	mockConnector := &MockDBConnector{
		ConnectFunc: func(ctx context.Context, dsn string) (DBConnection, error) {
			expectedDSN := "postgresql://testuser:testpass@localhost:5432/testdb?sslmode=disable"
			if dsn == expectedDSN {
				dsnWasCorrect = true
			}
			return mockConn, nil
		},
	}

	configLoader = mockConfig
	dbConnector = mockConnector

	conn, err := GetDatabaseConnection()

	assert.True(t, dsnWasCorrect, "DSN should match expected format")
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}

func TestGetDatabaseConnection_LoadConfigError(t *testing.T) {
	teardown := setupTest()
	defer teardown()

	expectedErr := errors.New("config file not found")

	mockConfig := &MockConfigLoader{
		LoadFunc: func() error {
			return expectedErr
		},
	}

	configLoader = mockConfig

	conn, err := GetDatabaseConnection()

	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "error loading .env file")
}

func TestGetDatabaseConnection_ConnectionError(t *testing.T) {
	teardown := setupTest()
	defer teardown()

	expectedErr := errors.New("unable to connect database")
	mockConfig := &MockConfigLoader{
		LoadFunc: func() error {
			return nil
		},
		GetEnvFunc: func(key string) string {
			envVars := map[string]string{
				"DATABASE_USER":     "testuser",
				"DATABASE_PASSWORD": "testpass",
				"DATABASE_URL":      "@localhost:5432/testdb",
				"SSL_LOCAL":         "disable",
			}
			return envVars[key]
		},
	}
	mockConnector := &MockDBConnector{
		ConnectFunc: func(ctx context.Context, dsn string) (DBConnection, error) {
			return nil, expectedErr
		},
	}

	configLoader = mockConfig
	dbConnector = mockConnector

	conn, err := GetDatabaseConnection()

	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "failed to connect database")
}

func TestCloseDatabaseConnection(t *testing.T) {
	mockConn := &MockDBConnection{}

	err := CloseDatabaseConnection(mockConn)

	assert.NoError(t, err)
	assert.True(t, mockConn.CloseCalled)
}

func TestCloseDatabaseConnection_Error(t *testing.T) {
	expectedErr := errors.New("closed failed")
	mockConn := &MockDBConnection{
		CloseError: expectedErr,
	}

	err := CloseDatabaseConnection(mockConn)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, mockConn.CloseCalled)
}

func TestCloseDatabseConnection_WithNilConnection(t *testing.T) {
	err := CloseDatabaseConnection(nil)
	assert.NoError(t, err)
	assert.Nil(t, err)
}


func TestTestDatabaseConnection_Success(t *testing.T) {
	teardown := setupTest()
	defer teardown()

	mockConfig := &MockConfigLoader{
		LoadFunc: func() error{
			return nil
		},
		GetEnvFunc: func (key string) string {
			return "test"
		},
	}

	mockRow := &MockRow{
		ScanFunc: func(dest ...interface{}) error {
			return nil
		},
	}

	mockConn := &MockDBConnection{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return mockRow
		},
	}

	mockConnector := &MockDBConnector{
		ConnectFunc: func(ctx context.Context, dsn string) (DBConnection, error) {
			return mockConn, nil
		},
	}

	configLoader = mockConfig
	dbConnector = mockConnector

	err := TestDatabaseConnection()

	assert.NoError(t, err)
	assert.True(t, mockConn.CloseCalled, "Connection should be closed")
}

func TestTestDatabaseConnection_ConnectionFailed(t *testing.T){
	teardown := setupTest()
	defer teardown()

	expectedErr := errors.New("database connection failed")

	mockConfig := &MockConfigLoader{
		LoadFunc: func() error{
			return nil
		},
		GetEnvFunc: func (key string) string {
			return "test"
		},
	}

	mockConnector := &MockDBConnector{
		ConnectFunc: func(ctx context.Context, dsn string) (DBConnection, error) {
			return nil, expectedErr
		},
	}

	configLoader = mockConfig
	dbConnector = mockConnector

	err := TestDatabaseConnection()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestTestDatabaseConnection_QueryFail(t *testing.T){
	teardown := setupTest()
	defer teardown()

	expectedErr := errors.New("failed to execute test query")

	mockConfig := &MockConfigLoader{
		LoadFunc: func() error{
			return nil
		},
		GetEnvFunc: func (key string) string {
			return "test"
		},
	}

	mockRow := &MockRow{
		ScanFunc: func(dest ...interface{}) error {
			return expectedErr
		},
	}
	mockConn := &MockDBConnection{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return mockRow
		},
	}

	mockConnector := &MockDBConnector{
		ConnectFunc: func(ctx context.Context, dsn string) (DBConnection, error) {
			return mockConn, nil
		},
	}

	configLoader = mockConfig
	dbConnector = mockConnector

	err := TestDatabaseConnection()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute test query")
	assert.True(t, mockConn.CloseCalled, "Connection should be closed")
}