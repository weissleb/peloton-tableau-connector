package extractors

import (
	"github.com/weissleb/peloton-tableau-connector/service/clients"
	"github.com/weissleb/peloton-tableau-connector/service/peloservice"
)

type PelotonClient struct {
	httpClient  *clients.HttpClientInterface
	userSession *peloservice.UserSession
}

func NewPelotonClient(httpClient clients.HttpClientInterface, username string, password string) (*PelotonClient, error) {
	session, err := peloservice.GetSession(httpClient, username, password)
	return &PelotonClient{
		httpClient:  &httpClient,
		userSession: &session,
	}, err
}

func ExistingPelotonClient(httpClient clients.HttpClientInterface, userid string, cookies string) (*PelotonClient, error) {
	session := peloservice.UserSession{
		UserId:  userid,
		Cookies: cookies,
	}

	return &PelotonClient{
		httpClient:  &httpClient,
		userSession: &session,
	}, nil
}

func (c PelotonClient) getHttpClient() *clients.HttpClientInterface {
	return c.httpClient
}

func (c PelotonClient) getUserSession() *peloservice.UserSession {
	return c.userSession
}

func (c PelotonClient) GetSessionUser() string {
	return c.userSession.UserId
}

func (c PelotonClient) GetSessionCookie() string {
	return c.userSession.Cookies
}
