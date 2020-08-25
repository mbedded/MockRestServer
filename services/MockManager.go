package services

import (
	"../models"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

type MockManager struct {
	ConnectionString string
	Database         *sql.DB
}

const SqlCreateTable string = `CREATE TABLE Mocks (
									id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
									key VARCHAR(250) NOT NULL,
									value TEXT NOT NULL);`
const SqlInsertMock string = `INSERT INTO Mocks(key, value) VALUES (?, ?)`
const SqlSelectMockByKey string = `SELECT id, key, value
										FROM Mocks
										WHERE key = ?;`

func NewMockManager(connectionString string) *MockManager {
	instance := &MockManager{
		ConnectionString: connectionString,
	}

	instance.InitializeDatabase()

	return instance
}

func (m *MockManager) InitializeDatabase() {
	db, err := sql.Open("sqlite3", m.ConnectionString)

	if err != nil {
		log.Fatal(err)
		return
	}

	m.Database = db
	_, err = m.Database.Exec(SqlCreateTable)

	if err != nil {
		log.Printf("Error creating database - %q", err)
	}
}

func (m *MockManager) CloseConnection() {
	if m.Database != nil {
		m.Database.Close()
	}
}

func (m *MockManager) SaveMockToDatabase(data models.JsonMockPost) (id string, err error) {
	key := strings.TrimSpace(data.Key)
	content := strings.TrimSpace(data.Content)

	if len(key) == 0 {
		key = uuid.New().String()
	}

	_, err = m.Database.Exec(SqlInsertMock, key, content)

	if err != nil {
		log.Fatalf("Unable to insert data to database %q", err)
	}

	return key, err
}

func (m *MockManager) GetMock(key string) (result models.JsonMockGet, err error) {
	row := m.Database.QueryRow(SqlSelectMockByKey, key)
	err = row.Scan(&result.Id, &result.Key, &result.Content)

	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}

		log.Fatalf("Error loading data from database. %q", err)
	}

	return result, err
}
