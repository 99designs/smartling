package smartling

// Credentials represents represents user credentials used to authenticate
// user in the Smartling API.
type Credentials struct {
	// UserID is a unique user ID for accessing Smartling API.
	UserID string

	// Secret is a secret token for accessing Smartling API.
	Secret string

	// AccessToken is a access token, which is obtained by UserID/Secret pair.
	AccessToken *Token

	// RefreshToken is a token for refreshing access token. It has longer
	// lifespan.
	RefreshToken *Token
}
