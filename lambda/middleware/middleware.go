package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := getTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body: "Missing auth token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body: "User unauthorised",
				StatusCode: http.StatusUnauthorized,
			}, err
		}
		expires := (int64)(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body: "Invalid token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		return next(request) 
	}
}

func getTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, error := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secret := "secret-key"
		return []byte(secret), nil 
	}) 
	if error != nil {
		return nil, fmt.Errorf("unauthorised")
	}
	if !token.Valid {
		return nil, fmt.Errorf("unauthorised - invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorised claims")
	}

	return claims, nil 
}