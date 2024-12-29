package awslambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	xcalistores3 "xcalistore-s3"
)

type InMemorySessionStore struct {
	allowedCredentials string
	sessions           []string
}

func (store *InMemorySessionStore) GetAllowedCredentials(ctx context.Context) (string, error) {
	return store.allowedCredentials, nil
}

func (store *InMemorySessionStore) CreateSession(ctx context.Context) (string, error) {
	sessionId := xcalistores3.SessionId()
	store.sessions = []string{sessionId}
	return sessionId, nil
}

func (store *InMemorySessionStore) ListSessions(ctx context.Context) ([]string, error) {
	return store.sessions, nil
}

func HandleEcho(ctx context.Context, event json.RawMessage) (LambdaResponseToAPIGW, error) {
	var response LambdaResponseToAPIGW
	var parsedEvent map[string]interface{}

	if eventParseErr := json.Unmarshal(event, &parsedEvent); eventParseErr != nil {
		log.Printf("Failed to unmarshal event: %v", eventParseErr)
		return response, eventParseErr
	}

	fmt.Printf("parsedEvent: %#v\n", parsedEvent)
	fmt.Printf("headers: %#v\n", parsedEvent["headers"])
	fmt.Printf("multiValueHeaders: %#v\n", parsedEvent["multiValueHeaders"])

	headers, headersCastOk := parsedEvent["headers"].(map[string]string)
	if headersCastOk {
		fmt.Printf("failed to cast headers:\n")
		return response, fmt.Errorf("failed to cast headers:\n")
	}

	sessMan := SessionManager{&InMemorySessionStore{}}
	sessionId, createSessErr := sessMan.checkCreateSession(ctx, headers)
	if createSessErr != nil {
		fmt.Printf("failed to create session: %v\n", createSessErr)

		var challange Challange
		if errors.As(createSessErr, &challange) {
			response, createRespErr := createResponse(true, "", nil)
			if createRespErr != nil {
				fmt.Printf("failed to create response: %v\n", createRespErr)
				return response, createSessErr
			}
			return response, nil
		}

		return response, createSessErr
	}

	response, createRespErr := createResponse(false, sessionId, map[string]string{"message": "hello, xcali!"})
	if createRespErr != nil {
		fmt.Printf("failed to create response: %v\n", createRespErr)
		return response, createSessErr
	}

	return response, nil
}
