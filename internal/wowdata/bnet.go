package wowdata

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type bnetClient struct {
	httpClient  *http.Client
	accessToken string
}

const bnetOAuthURL = "https://oauth.battle.net/token"

var clientCredentials = (url.Values{"grant_type": {"client_credentials"}}).Encode()

func newBnet(ctx context.Context, clientID, clientSecret string) (*bnetClient, error) {
	hc := http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, bnetOAuthURL, strings.NewReader(clientCredentials))
	if err != nil {
		return nil, fmt.Errorf("create oauth token request: %w", err)
	}
	password := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", password))

	res, err := hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send oauth token request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("oauth token error response: %v %v", res.StatusCode, res.Status)
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
