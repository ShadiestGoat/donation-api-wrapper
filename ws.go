package donations

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type eventType int

const (
	et_NONE eventType = iota
	et_NEW_DON
	et_NEW_FUND
	et_PING

	// local handlers only!
	et_CLOSE
	et_OPEN
)

type wsEvent struct {
	Type eventType       `json:"event"`
	Body json.RawMessage `json:"body"`
}

// Fired whenever a new fund is created
type EventNewFund struct {
	*Fund
}

// An event that is fired whenever a new donation is made
type EventNewDonation struct {
	*Donation
}

// An event for whenever the WS connection closes
type EventClose struct {
	Err error
}

// An empty event for whenever the WS connection opens
type EventOpen struct {}

// Add a handler for WS events, where h is func(c *Client, h *Event{event name})
// Events supported are currently EventNewFund, and EventNewDonation, EventClose
func (c *Client) AddHandler(h any) {
	c.Lock()
	defer c.Unlock()

	switch h := h.(type) {
	case func(c *Client, v *EventNewFund):
		c.handlers[et_NEW_FUND] = h
	case func(c *Client, v *EventNewDonation):
		c.handlers[et_NEW_DON] = h
	case func(c *Client, v *EventClose):
		c.handlers[et_CLOSE] = h
	case func(c *Client, v *EventOpen):
		c.handlers[et_OPEN] = h
	}
}

func (c *Client) sendEvent(payload any) {
	var ev eventType

	switch payload.(type) {
	case *EventNewFund:
		ev = et_NEW_FUND
	case *EventNewDonation:
		ev = et_NEW_DON
	case *EventClose:
		ev = et_CLOSE
	case *EventOpen:
		ev = et_CLOSE
	}

	if h, ok := c.handlers[ev]; ok {
		switch h := h.(type) {
		case func(c *Client, v *EventNewFund):
			go h(c, payload.(*EventNewFund))
		case func(c *Client, v *EventNewDonation):
			go h(c, payload.(*EventNewDonation))
		case func(c *Client, v *EventClose):
			go h(c, payload.(*EventClose))
		case func(c *Client, v *EventOpen):
			go h(c, payload.(*EventOpen))
		}
	}
}

// Opens WS connection
func (c *Client) OpenWS() error {
	c.Lock()
	defer c.Unlock()
	if c.wsConn != nil {
		return ErrAlrConnected
	}
	conn, resp, err := websocket.DefaultDialer.Dial("wss://"+c.location+"/ws", http.Header{
		"Authorization": []string{c.token},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		r := []byte{}
		if resp.Body != nil {
			r, _ = io.ReadAll(resp.Body)
		}

		return &HTTPError{
			Status:  resp.StatusCode,
			Message: string(r),
		}
	}

	c.wsConn = conn

	go c.wsLoop()

	c.sendEvent(&EventOpen{})

	return nil
}

func (c *Client) CloseWS() {
	c.Lock()
	defer c.Unlock()

	if c.wsConn != nil {
		c.wsConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, ""), time.Time{})
		c.wsConn.Close()
		c.wsConn = nil
	}
}

func (c *Client) wsLoop() {
	var err error

	defer func() {
		// calling this just in case
		c.CloseWS()

		c.sendEvent(&EventClose{
			Err: err,
		})
	}()

	for {
		_, b, err := c.wsConn.ReadMessage()
		if err != nil {
			return
		}
		e := &wsEvent{}
		err = json.Unmarshal(b, &e)
		if err != nil {
			return
		}

		var ev any

		switch e.Type {
		case et_PING:
			c.wsConn.WriteMessage(websocket.TextMessage, []byte{'P'})
			continue
		case et_NEW_DON:
			d := &Donation{}
			err = json.Unmarshal(e.Body, &d)

			if err != nil {
				return
			}

			ev = &EventNewDonation{Donation: d}
		case et_NEW_FUND:
			f := &Fund{}
			err = json.Unmarshal(e.Body, &f)

			if err != nil {
				return
			}

			ev = &EventNewFund{Fund: f}
		}

		c.sendEvent(ev)
	}
}
