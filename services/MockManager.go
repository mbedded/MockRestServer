package services

import (
	"../models"
)

type MockManager struct {
	ConnectionString string
}

func NewMockManager(connectionString string) *MockManager {
	return &MockManager{
		ConnectionString: connectionString,
	}
}

func (m *MockManager) CreateMock(data models.JsonMockPost) (id string, err error) {
	return "MOCK-Data_123", nil
}
