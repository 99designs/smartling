package smartling

import (
	"fmt"
	"time"
)

func (client *Client) logRequest(
	method string,
	url string,
	body []byte,
) {
	token := "[no token]"
	if client.Credentials.AccessToken != nil {
		token = fmt.Sprintf(
			"[token=%s...{%d bytes} ttl %.2fs]",
			client.Credentials.AccessToken.Value[:7],
			len(client.Credentials.AccessToken.Value),
			client.Credentials.AccessToken.ExpirationTime.Sub(
				time.Now(),
			).Seconds(),
		)
	}

	client.Logger.Debugf(
		"<- %s %s %s [%d bytes body]", method, url, token, len(body),
	)
}
