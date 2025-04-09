package main

import (
	"fmt"
	"net/http"

	"lambda-func/app"
	"lambda-func/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Take in payload and process it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfuly called by - %s", event.Username), nil 
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "You've hit a protected route",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp()

	// lambda.Start(myApp.ApiHandler.RegisterUserHandler)
	lambda.Start(func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register": 
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUser(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default: 
			return events.APIGatewayProxyResponse{
				Body: "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil 
		}
	})
}