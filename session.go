package betfair

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Call betfair api timeoout
var ClientTimeout = time.Second * 10

// Maximum of parallel requests per session
var HTTPClientPoolSize = 100

// Login betfair endpoint
var InteractiveLoginEndpoint = "https://identitysso-api.betfair.com/api/login"

// Certificate login method endpoint
var NonInteractiveLoginEndpoint = "https://identitysso-api.betfair.com/api/certlogin"

// Keep session token alive endpoint
var KeepAliveEndpoint = "https://identitysso.betfair.com/api/keepAlive"

type pooledHTTPClient struct {
	*http.Client
	poolCh chan bool
}

func (client *pooledHTTPClient) Do(req *http.Request) ([]byte, error) {
	client.checkoutConnection()
	defer client.releaseConnection()

	res, err := client.Client.Do(req)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

func (client *pooledHTTPClient) checkoutConnection() {
	<-client.poolCh
}

func (client *pooledHTTPClient) releaseConnection() {
	client.poolCh <- true
}

func initializeHTTPClient(account *Account) (*pooledHTTPClient, error) {
	pool := make(chan bool, HTTPClientPoolSize)

	// fill pool with values
	for i := 0; i < HTTPClientPoolSize; i++ {
		pool <- true
	}

	if account.LoginMethod == NoneInteractive {
		ssl := &tls.Config{
			Certificates:       []tls.Certificate{account.Certificate},
			InsecureSkipVerify: true,
		}

		ssl.Rand = rand.Reader

		var httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(network, addr, time.Duration(ClientTimeout))
				},
				TLSClientConfig: ssl,
			},
		}

		return &pooledHTTPClient{httpClient, pool}, nil
	}

	var httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, time.Duration(ClientTimeout))
			},
		},
	}

	return &pooledHTTPClient{httpClient, pool}, nil
}

type Session struct {
	ssoid      string
	account    *Account
	httpClient *pooledHTTPClient
	m          sync.Mutex
}

func NewSession(account *Account) (*Session, error) {
	httpClient, err := initializeHTTPClient(account)

	if err != nil {
		return nil, err
	}

	return &Session{account: account, httpClient: httpClient}, nil
}

func (session *Session) GetToken() (string, error) {
	if session.ssoid == "" {
		session.m.Lock()
		defer session.m.Unlock()

		if session.ssoid == "" {
			ssoid, err := session.requestSsoid()

			if err != nil {
				return ssoid, err
			}

			session.ssoid = ssoid

			if session.account.KeepAlive {
				go session.startKeepAliveLoop()
			}
		}
	}

	return session.ssoid, nil
}

var keepAliveReader = strings.NewReader("")

func (session *Session) KeepAlive() (bool, error) {
	var payload keepAliveResult
	err := session.doRequest(&payload, KeepAliveEndpoint, keepAliveReader)

	if err != nil {
		return false, err
	}

	if payload.Status == "SUCCESS" {
		return true, nil
	}

	return false, errors.New(payload.Error)
}

func (session *Session) doRequest(payload interface{}, endpoint string, body io.Reader) error {
	resBody, err := session.doRawRequest("POST", endpoint, body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, payload)

	if err != nil {
		return err
	}

	return nil
}

func (session *Session) doRawRequest(httpMethod, endpoint string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Application", session.account.ApplicationKey)

	token, err := session.GetToken()

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Authentication", token)

	return session.httpClient.Do(req)
}

func (session *Session) startKeepAliveLoop() error {
	for {
		time.Sleep(10 * time.Minute)
		session.KeepAlive()
	}
}

func (session *Session) requestSsoid() (string, error) {
	body := strings.NewReader(fmt.Sprintf("username=%s&password=%s", session.account.Username, session.account.Password))
	req, err := http.NewRequest("POST", session.loginEndpoint(), body)

	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Application", session.account.ApplicationKey)

	resBody, err := session.httpClient.Do(req)

	if err != nil {
		return "", err
	}

	if session.account.LoginMethod == Interactive {
		var response = InteractiveSessionResponse{}
		err = json.Unmarshal(resBody, &response)

		if err != nil {
			return "", err
		}

		if response.Status == "SUCCESS" {
			return response.Token, nil
		}

		return "", errors.New(response.Error)
	}

	var response = NoneInteractiveSessionResponse{}
	err = json.Unmarshal(resBody, &response)

	if err != nil {
		return "", err
	}

	if response.LoginStatus == "SUCCESS" {
		return response.SessionToken, nil
	}

	return "", errors.New(response.LoginStatus)
}

func (session *Session) loginEndpoint() string {
	if session.account.LoginMethod == Interactive {
		return InteractiveLoginEndpoint
	}

	return NonInteractiveLoginEndpoint
}
