package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAdapter struct {
	config *oauth2.Config
}

func NewGoogleAdapter(clientID, clientSecret, redirectURL string) *GoogleAdapter {
	return &GoogleAdapter{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (a *GoogleAdapter) AuthURL(_ context.Context) string {
	return a.config.AuthCodeURL("state", oauth2.AccessTypeOnline)
}

type GoogleUser struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (a *GoogleAdapter) Exchange(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := a.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	client := a.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("fetch userinfo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GoogleUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &user, nil
}
