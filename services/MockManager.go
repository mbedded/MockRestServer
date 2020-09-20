package services

import (
	"../models"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type MockManager struct {
	ConnectionString string
	Database         *sql.DB
}

const SqlCreateTable string = `CREATE TABLE Mocks (
									id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
									key VARCHAR(250) NOT NULL,
									value TEXT NOT NULL);`
const SqlInsertMock string = `INSERT INTO Mocks(key, value) VALUES (?, ?);`
const SqlSelectMockByKey string = `SELECT id, key, value FROM Mocks WHERE key = ?;`
const SqlCountByKey string = `SELECT COUNT(key) FROM Mocks WHERE key = ?;`
const SqlUpdateMock string = `UPDATE Mocks SET value = ? WHERE key = ?;`

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

func (m *MockManager) SaveMockToDatabase(key string, content string) (string, error) {
	if len(key) == 0 {
		key = uuid.New().String()
	}

	_, err := m.Database.Exec(SqlInsertMock, key, content)

	if err != nil {
		log.Fatalf("Unable to insert data to database %q", err)
		return "", err
	}

	return key, err
}

func (m *MockManager) ContainsKey(key string) (isExisting bool, err error) {
	var count int
	row := m.Database.QueryRow(SqlCountByKey, key)
	err = row.Scan(&count)

	if err != nil {
		log.Fatalf("Error counting keys in database. %q", err)
	}

	return count >= 1, nil
}

func (m *MockManager) UpdateMock(key string, content string) (err error) {
	_, err = m.Database.Exec(SqlUpdateMock, content, key)

	if err != nil {
		log.Fatalf("Error counting keys in database. %q", err)
	}

	return err
}

func (m *MockManager) GetMock(key string) (result models.JsonMockGet, err error) {
	row := m.Database.QueryRow(SqlSelectMockByKey, key)
	err = row.Scan(&result.Id, &result.Key, &result.Content)

	if err == sql.ErrNoRows {
		return result, nil
	}

	if err != nil {
		log.Fatalf("Error loading data from database. %q", err)
	}

	return result, err
}
