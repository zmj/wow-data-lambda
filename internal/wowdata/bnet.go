package wowdata

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type bnetClient struct {
	httpClient  *http.Client
	accessToken string
}

var clientCredentials = (url.Values{"grant_type": {"client_credentials"}}).Encode()

func newBnet(ctx context.Context, o oauthClient) (*bnetClient, error) {
	hc := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth.battle.net/token", strings.NewReader(clientCredentials))
	if err != nil {
		return nil, fmt.Errorf("create oauth token request: %w", err)
	}
	password := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", o.ClientID, o.ClientSecret)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", password))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send oauth token request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		var body string
		if b, err := io.ReadAll(res.Body); err == nil {
			body = string(b)
		}
		return nil, fmt.Errorf("oauth token error response: %v %v", res.Status, body)
	}

	var tokenRes tokenResponse
	if err = json.NewDecoder(res.Body).Decode(&tokenRes); err != nil {
		return nil, fmt.Errorf("deserialize oauth token response: %w", err)
	} else if tokenRes.AccessToken == "" {
		return nil, fmt.Errorf("unexpected oauth token response")
	}

	exp := time.Duration(tokenRes.ExpiresIn) * time.Second
	logger(ctx).DebugContext(ctx, "got bnet oauth token", "exp", exp)
	return &bnetClient{httpClient: &hc, accessToken: tokenRes.AccessToken}, nil
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
