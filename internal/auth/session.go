package auth

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"
)

type SessionStorage struct {
	session sessions.Session
}

func (storage *SessionStorage) GetCodeVerifier() (string, error) {
	v := storage.session.Get("code_verifier")
	if v == nil {
		return "", fmt.Errorf("code_verifier not found in Session")
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid code_verifier type in Session")
	}
	return s, nil
}

func (storage *SessionStorage) SetCodeVerifier(codeVerifier string) error {
	storage.session.Set("code_verifier", codeVerifier)
	if err := storage.session.Save(); err != nil {
		return fmt.Errorf("failed to save Session: %w", err)
	}
	return nil
}

func (storage *SessionStorage) GetRawToken() (*oauth2.Token, error) {
	tokenData := storage.session.Get("kinde_token")
	if tokenData == nil {
		return nil, fmt.Errorf("token not found in Session")
	}

	// Handle both JSON string (new format) and direct token (old format for backwards compatibility)
	var tokenBytes []byte
	switch v := tokenData.(type) {
	case string:
		tokenBytes = []byte(v)
	case []byte:
		tokenBytes = v
	default:
		// Try to marshal if it's already a token (backwards compatibility)
		if t, ok := tokenData.(*oauth2.Token); ok {
			return t, nil
		}
		return nil, fmt.Errorf("invalid token type in Session")
	}

	var t oauth2.Token
	if err := json.Unmarshal(tokenBytes, &t); err != nil {
		return nil, fmt.Errorf("invalid token type in Session")
	}
	return &t, nil
}

func (storage *SessionStorage) SetRawToken(token *oauth2.Token) error {
	if token == nil {
		storage.session.Set("kinde_token", nil)
	} else {
		// Serialize token to JSON to avoid gob serialization issues
		tokenBytes, err := json.Marshal(token)
		if err != nil {
			return fmt.Errorf("failed to marshal token: %w", err)
		}
		storage.session.Set("kinde_token", string(tokenBytes))
	}
	if err := storage.session.Save(); err != nil {
		return fmt.Errorf("failed to save Session: %w", err)
	}
	return nil
}

func (storage *SessionStorage) GetPostAuthRedirect() (string, error) {
	v := storage.session.Get("post_auth_redirect")
	if v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid post_auth_redirect type in Session")
	}
	return s, nil
}

func (storage *SessionStorage) GetState() (string, error) {
	v := storage.session.Get("auth_state")
	if v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid auth_state type in Session")
	}
	return s, nil
}

func (storage *SessionStorage) SetPostAuthRedirect(redirect string) error {
	storage.session.Set("post_auth_redirect", redirect)
	if err := storage.session.Save(); err != nil {
		return fmt.Errorf("failed to save Session: %w", err)
	}
	return nil
}

func (storage *SessionStorage) SetState(state string) error {
	storage.session.Set("auth_state", state)
	if err := storage.session.Save(); err != nil {
		return fmt.Errorf("failed to save Session: %w", err)
	}
	return nil
}

func (storage *SessionStorage) GetItem(key string) string {
	value := storage.session.Get(key)
	if value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

func (storage *SessionStorage) SetItem(key, value string) {
	storage.session.Set(key, value)
	_ = storage.session.Save() // Error ignored for backwards compatibility
}
