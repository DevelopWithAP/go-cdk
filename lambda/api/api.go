package api

import (
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
)

type ApiHandler struct {
	dbStore database.DynamoDBClient 
}

func NewApiHandler(dbStore database.DynamoDBClient) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("request has empty parameters")
	}

	exists, error := api.dbStore.UserExists(event.Username)
	if error != nil {
		return fmt.Errorf("an error occurred: %w", error)
	}

	if exists {
		return fmt.Errorf("username taken")
	}

	err := api.dbStore.InsertUser(event)

	if err != nil {
		return fmt.Errorf("error registering the user %w", err)
	}

	return nil
} 