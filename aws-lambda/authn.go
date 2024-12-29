package awslambda

import (
	"context"
	"fmt"
	"net/http"
	"slices"
)

const (
	sessionCookieName = "xcaliapp-session"
)

type SessionStore interface {
	GetAllowedCredentials(ctx context.Context) (string, error)
	CreateSession(ctx context.Context) (string, error)
	ListSessions(ctx context.Context) ([]string, error)
}

type SessionManager struct {
	store SessionStore
}

type Challange struct{}

func (challange *Challange) Error() string {
	return "WWW-Authenticate" // "Basic"
}

func (manager *SessionManager) checkCreateSession(ctx context.Context, headers map[string]string) (string, error) {
	cookieHeader, sessionCookieFound := headers[sessionCookieName]

	if sessionCookieFound {
		sessions, listSessionErr := manager.store.ListSessions(ctx)
		if listSessionErr != nil {
			return "", fmt.Errorf("failed to list sessions: %w", listSessionErr)
		}
		cookies, parseCookieErr := http.ParseCookie(cookieHeader)
		if parseCookieErr != nil {
			return "", fmt.Errorf("failed to parse cookie: %w", parseCookieErr)
		}
		for _, cookie := range cookies {
			if slices.Contains(sessions, cookie.Value) {
				return "", nil
			}
		}
		return "", &Challange{}
	}

	authrHeader, authrHeaderFound := headers["Authorization"]
	if !authrHeaderFound {
		return "", &Challange{}
	}

	allowedCred, getAllowedCredErr := manager.store.GetAllowedCredentials(ctx)
	if getAllowedCredErr != nil {
		return "", fmt.Errorf("failed to get allowed credentials: %w", getAllowedCredErr)
	}

	if allowedCred == authrHeader {
		return "", nil
	}

	sessionId, createSessIdErr := manager.store.CreateSession(ctx)
	if createSessIdErr != nil {
		return "", fmt.Errorf("failed to create session: %w", createSessIdErr)
	}

	return sessionId, &Challange{}
}
