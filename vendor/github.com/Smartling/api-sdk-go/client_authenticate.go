package smartling

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	endpointAuthenticate        = "/auth-api/v2/authenticate"
	endpointAuthenticateRefresh = "/auth-api/v2/authenticate/refresh"
)

// Authenticate checks that access and refresh tokens are valid and refreshes
// them if needed.
func (client *Client) Authenticate() error {
	if client.Credentials.AccessToken.IsSafe() {
		return nil
	}

	if client.Credentials.UserID == "" {
		return fmt.Errorf("user ID in the client is not set")
	}

	if client.Credentials.Secret == "" {
		return fmt.Errorf("token secret in the client is not set")
	}

	var (
		url    string
		params map[string]string
	)

	client.Credentials.AccessToken = nil

	if client.Credentials.RefreshToken.IsSafe() {
		url = endpointAuthenticateRefresh
		params = map[string]string{
			"refreshToken": client.Credentials.RefreshToken.Value,
		}
	} else {
		url = endpointAuthenticate
		params = map[string]string{
			"userIdentifier": client.Credentials.UserID,
			"userSecret":     client.Credentials.Secret,
		}

		client.Credentials.RefreshToken = nil
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf(
			"unable to encode authenticate params: %s", err,
		)
	}

	var response struct {
		AccessToken      string
		ExpiresIn        time.Duration
		RefreshToken     string
		RefreshExpiresIn time.Duration
	}

	_, _, err = client.Post(url, payload, &response, WithoutAuthentication)
	if err != nil {
		if _, ok := err.(NotAuthorizedError); ok {
			return err
		}

		return fmt.Errorf(
			"authenticate request failed: %s", err,
		)
	}

	client.Credentials.AccessToken = &Token{
		Value: response.AccessToken,
		ExpirationTime: time.Now().Add(
			response.ExpiresIn * time.Second,
		),
	}

	client.Credentials.RefreshToken = &Token{
		Value: response.RefreshToken,
		ExpirationTime: time.Now().Add(
			time.Second * response.RefreshExpiresIn,
		),
	}

	return nil
}
