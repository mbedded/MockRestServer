package services

import (
	"../models"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DatabaseManager struct {
	ConnectionString string
	Database         *sql.DB
}

const SqlCreateTable string = `CREATE TABLE Mocks (
									id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
									key VARCHAR(250) NOT NULL,
									value TEXT NOT NULL);`
const SqlInsertMock string = `INSERT INTO Mocks(key, value) VALUES (?, ?);`
const SqlSelectMockByKey string = `SELECT id, key, value FROM Mocks WHERE key = ?;`
const SqlSelectAllMocks string = `SELECT id, key, value FROM Mocks ORDER BY key;`
const SqlCountNumberOfMocks = `SELECT count(id) FROM Mocks;`
const SqlCountByKey string = `SELECT COUNT(key) FROM Mocks WHERE key = ?;`
const SqlUpdateMock string = `UPDATE Mocks SET value = ? WHERE key = ?;`
const SqlDeleteMockByKey = `DELETE FROM Mocks WHERE key = ?;`

func NewDatabaseManager(connectionString string) *DatabaseManager {
	instance := &DatabaseManager{
		ConnectionString: connectionString,
	}

	instance.InitializeDatabase()

	return instance
}

func (m *DatabaseManager) InitializeDatabase() {
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

func (m *DatabaseManager) CloseConnection() {
	if m.Database != nil {
		m.Database.Close()
	}
}

func (m *DatabaseManager) SaveMockToDatabase(key string, content string) (string, error) {
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

func (m *DatabaseManager) ContainsKey(key string) (isExisting bool, err error) {
	var count int
	row := m.Database.QueryRow(SqlCountByKey, key)
	err = row.Scan(&count)

	if err != nil {
		log.Fatalf("Error counting keys in database. %q", err)
	}

	return count >= 1, nil
}

func (m *DatabaseManager) UpdateMock(key string, content string) (err error) {
	_, err = m.Database.Exec(SqlUpdateMock, content, key)

	if err != nil {
		log.Fatalf("Error counting keys in database. %q", err)
	}

	return err
}

func (m *DatabaseManager) GetMock(key string) (result models.JsonMockGet, err error) {
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

func (m *DatabaseManager) GetAll() (result []models.JsonMockGet, err error) {
	var numberOfMocks int
	row := m.Database.QueryRow(SqlCountNumberOfMocks)
	err = row.Scan(&numberOfMocks)

	if err != nil {
		log.Fatalf("Error counting Mocks in database. %q", err)
	}

	rows, err := m.Database.Query(SqlSelectAllMocks)

	if err != nil {
		log.Fatalf("Error reading all mocks from database. %q", err)
	}

	result = make([]models.JsonMockGet, numberOfMocks)
	var index = 0

	for rows.Next() && index < numberOfMocks {
		item := models.JsonMockGet{}

		err = rows.Scan(&item.Id, &item.Key, &item.Content)

		if err != nil {
			log.Fatalf("Error reading row from database. %q", err)
		}

		result[index] = item
		index++
	}

	return result, nil
}

func (m *DatabaseManager) DeleteMock(key string) error {
	_, err := m.Database.Exec(SqlDeleteMockByKey, key)
	return err
}
