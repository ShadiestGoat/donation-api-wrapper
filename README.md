# Donation API Wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/ShadiestGoat/donation-api-wrapper.svg)](https://pkg.go.dev/github.com/ShadiestGoat/github.com/ShadiestGoat/donation-api-wrapper)

## How to use

1. Setup a authed app in the auths.json file on the server
2. Use the token in `NewClient({token})`
3. Profit

## Websocket API

There is websocket support, made through callbacks.

```go

// ... in main()
// c is of type *Client

c.AddHandler(func (c *Client, v *EventClose)) {
    log.Error(v.Err)
    time.Sleep(30 * time.Second)
    c.OpenWS()
}

c.OpenWS()

```

## Rest API

API Routes are abstracted through the use of the `*Client` through methods.