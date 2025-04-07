package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	// db store initialised here
	// passed down to api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db) 

	return App {
		ApiHandler: apiHandler,
	}
}