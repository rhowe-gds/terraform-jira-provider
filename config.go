package main

import (
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "log"
  "strings"

  "golang.org/x/net/context"

  "github.com/dghubble/oauth1"
  jira "gopkg.in/andygrunwald/go-jira.v1"
)

// Config holds API and APP keys to authenticate to Datadog.
type Config struct {
  AuthMethod string
  BaseUrl string
  Username string
  Password string
  OauthConsumerKey string
  OauthPrivateKey string
}

func getJIRATokenFromWeb(config *oauth1.Config) *oauth1.Token {
  requestToken, requestSecret, err := config.RequestToken()
  if err != nil {
    log.Fatalf("Unable to get request token. %v", err)
  }
  authorizationURL, err := config.AuthorizationURL(requestToken)
  if err != nil {
    log.Fatalf("Unable to get authorization url. %v", err)
  }
  fmt.Printf("Go to the following link in your browser then type the "+
    "authorization code: \n%v\n", authorizationURL.String())

  var code string
  if _, err := fmt.Scan(&code); err != nil {
    log.Fatalf("Unable to read authorization code. %v", err)
  }

  accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, code)
  if err != nil {
    log.Fatalf("Unable to get access token. %v", err)
  }
  return oauth1.NewToken(accessToken, accessSecret)
}

// Client returns a new Jira client.
func (c *Config) Client() *jira.Client {
  switch c.AuthMethod {
  case "basic":
    transport := jira.BasicAuthTransport{
      Username: c.Username,
      Password: c.Password,
    }

    client, _ := jira.NewClient(transport.Client(), c.BaseUrl)

    return client
  case "cookie":
    transport := jira.CookieAuthTransport{
      Username: c.Username,
      Password: c.Password,
      AuthURL: fmt.Sprintf("%s/rest/auth/1/session", c.BaseUrl),
    }

    client, _ := jira.NewClient(transport.Client(), c.BaseUrl)

    return client
  case "oauth":
    keyDERBlock, _ := pem.Decode([]byte(c.OauthPrivateKey))
    if keyDERBlock == nil {
      log.Fatal("unable to decode key PEM block")
    }
    if !(keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {
      log.Fatalf("unexpected key DER block type: %s", keyDERBlock.Type)
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
    if err != nil {
      log.Fatalf("unable to parse PKCS1 private key. %v", err)
    }
    config := oauth1.Config{
      ConsumerKey: c.OauthConsumerKey,
      CallbackURL: "oob", /* for command line usage */
      Endpoint: oauth1.Endpoint{
        RequestTokenURL: c.BaseUrl + "plugins/servlet/oauth/request-token",
        AuthorizeURL:    c.BaseUrl + "plugins/servlet/oauth/authorize",
        AccessTokenURL:  c.BaseUrl + "plugins/servlet/oauth/access-token",
      },
      Signer: &oauth1.RSASigner{
        PrivateKey: privateKey,
      },
    }

    token := getJIRATokenFromWeb(&config)

    client, _ := jira.NewClient(config.Client(context.Background(), token), c.BaseUrl)

    return client
  case "none":
    client, _ := jira.NewClient(nil, c.BaseUrl)

    return client
  }
  return nil
}
