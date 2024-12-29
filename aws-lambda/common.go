package awslambda

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LambdaResponseToAPIGW struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
}

func createResponse(challange bool, session string, body map[string]any) (LambdaResponseToAPIGW, error) {
	var respStruct LambdaResponseToAPIGW
	var headers map[string]string
	bodyToSend := ""

	if challange && len(session) > 0 {
		return respStruct, fmt.Errorf("invalid arguments: either challange or session, not both")
	}

	if challange {
		return LambdaResponseToAPIGW{
			StatusCode:        4011,
			Headers:           map[string]string{"WWW-Authenticate": "Basic"},
			IsBase64Encoded:   false,
			MultiValueHeaders: nil,
			Body:              "",
		}, nil
	}

	if len(session) > 0 {
		cookieToSet := &http.Cookie{
			Name:     sessionCookieName,
			Value:    session,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		if headers == nil {
			headers = map[string]string{"Set-Cookie": cookieToSet.String()}
		}
	}

	if len(body) > 0 {
		bodyToSendInBytes, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return respStruct, marshalErr
		}
		bodyToSend = string(bodyToSendInBytes)
	}

	respStruct = LambdaResponseToAPIGW{
		StatusCode:        4011,
		Headers:           headers,
		IsBase64Encoded:   false,
		MultiValueHeaders: nil,
		Body:              bodyToSend,
	}

	return respStruct, nil
}
