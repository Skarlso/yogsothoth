package yogsot

import (
	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Token is a access token for DO.
type Token struct {
	AccessToken string
}

// YogClient is a client struct for Yogsothoth.
type YogClient struct {
	*godo.Client
}

// Token provides a convinient method to setup a Token.
func (t *Token) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// NewClient produces a client for Yogsothoth.
func NewClient(token string) *YogClient {
	oauthClient := oauth2.NewClient(context.Background(), &Token{AccessToken: token})
	yogClient := YogClient{godo.NewClient(oauthClient)}
	return &yogClient
}

// NewContext provides a new context.
func NewContext() context.Context {
	return context.TODO()
}
