package smartling

import "fmt"
import "time"

const tokenExpirationSafetyDuration = 30 * time.Second

// Token represents authentication token, either access or refresh.
type Token struct {
	// Value is a string representation of token.
	Value string

	// ExpirationTime is a expiration time for token when it becomes invalid.
	ExpirationTime time.Time
}

// IsValid returns true if token still can be used.
func (token *Token) IsValid() bool {
	if token == nil {
		return false
	}

	if token.Value == "" {
		return false
	}

	return time.Now().Before(token.ExpirationTime)
}

// IsSafe returns true if token still can be used and it's expiration time is
// in safe bounds.
func (token *Token) IsSafe() bool {
	if !token.IsValid() {
		return false
	}

	return time.Now().Add(tokenExpirationSafetyDuration).Before(
		token.ExpirationTime,
	)
}

// String returns token representation for logging purposes only.
func (token *Token) String() string {
	if token == nil {
		return "[no token]"
	}

	return fmt.Sprintf(
		"[token=%s...{%d bytes} ttl %.2fs]",
		token.Value[:7],
		len(token.Value),
		token.ExpirationTime.Sub(time.Now()).Seconds(),
	)
}
