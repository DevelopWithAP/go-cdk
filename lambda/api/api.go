package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore 
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registeredUser types.RegisterUser

	error := json.Unmarshal([]byte(request.Body), &registeredUser)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, error 
	}

	if registeredUser.Username == "" || registeredUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Empty username or password",
			StatusCode: http.StatusBadRequest,
		}, error 
	}

	exists, error := api.dbStore.UserExists(registeredUser.Username)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body: "Iternal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, error
	}

	user, error := types.NewUser(registeredUser)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("unable to create new user %w", error)
	}


	if exists {
		return events.APIGatewayProxyResponse{
			Body: "User exists",
			StatusCode: http.StatusConflict,
		}, error
	}

	err := api.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Iternal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body: "User successfully registered",
		StatusCode: http.StatusCreated,
	}, err
} 

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginRequest types.LoginRequest

	error := json.Unmarshal([]byte(request.Body), &loginRequest)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, error
	}

	exists, error := api.dbStore.GetUser(loginRequest.Username)
	if error != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, error
	}

	if !types.ValidatePassword(exists.HashedPassword, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body: "Invalid user credentials",
			StatusCode: http.StatusBadRequest,
		}, nil 
	}
	return events.APIGatewayProxyResponse{
		Body: "Successfully logged in",
		StatusCode: http.StatusOK,
	}, nil 
}
