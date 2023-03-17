package donations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	sync.Mutex
	httpClient *http.Client
	wsConn *websocket.Conn
	handlers map[eventType]any
	// the location of the api w/o a protocol
	location string
	token string
}

type clientOpt func (*Client)

// A client opt that will use a custom http.Client for it's requests
func WithCustomHTTPClient(httpClient *http.Client) clientOpt {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// A client opt that will make the client use a different location. This should be a hostname without any protocols. 
// If you are running the API under a non-root path, you should also specify that here. 
// Do not include a trailing slash.
func WithCustomLocation(loc string) clientOpt {
	return func(c *Client) {
		c.location = loc
	}
} 

func NewClient(token string, opts ...clientOpt) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
		wsConn:     nil,
		location:   "donate.shadygoat.eu/api",
		token: 		token,
		handlers: map[eventType]any{},
	}

	for _, o := range opts {
		o(c)
	}
	
	return c
}

var ErrNilResp = errors.New(`the response received was nil, but a response was specified`)
var ErrAlrConnected = errors.New(`the ws was already connected`)

type HTTPError struct {
	Status int
	Message string `json:"error"`
}

func (err HTTPError) Error() string {
	return fmt.Sprint(err.Status) + ": " + err.Message
}

// Fetch from the API, where m
// m is the http method
// path is the api path,
// body is a pointer representing the body to be sent. This should be set to nil if there is no body.
// resp is a pointer to what will be used as the response body. This should be set to nil if the response body can be disregarded.
func (c *Client) fetch(m string, path string, body any, resp any) error {
	uri := "https://" + c.location + path
	
	var bodyToSend io.Reader

	if body != nil {
		bRaw := bytes.NewBuffer(nil)
		
		err := json.NewEncoder(bRaw).Encode(body)
		if err != nil {
			return err
		}

		bodyToSend = bRaw
	}

	req, err := http.NewRequest(m, uri, bodyToSend)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.token)

	respRaw, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if respRaw == nil {
		return ErrNilResp
	}

	if respRaw.StatusCode != http.StatusOK {
		s := &HTTPError{
			Status:  respRaw.StatusCode,
			Message: "",
		}

		if respRaw.Body != nil {
			json.NewDecoder(respRaw.Body).Decode(s)
		}

		return s
	}

	if resp != nil {
		if respRaw == nil || respRaw.Body == nil {
			return ErrNilResp
		}
		err = json.NewDecoder(respRaw.Body).Decode(resp)
		if err != nil {
			return err
		}
	}

	return nil
}
