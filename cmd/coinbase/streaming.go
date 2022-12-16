package main

import (
	"errors"
	"math"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cgmello/marketdata/config"
	"github.com/cgmello/marketdata/model"
	"github.com/gorilla/websocket"
)

type CoinbaseStream struct {
	Conn         *websocket.Conn // WebSocket connection
	Url          string
	SubscribeMsg model.CoinbaseSubscribeMessage
	Retries      int
}

func NewCoinbaseStream() *CoinbaseStream {
	url := config.CONFIG_MAP["COINBASE_URL"]
	productIds := strings.Split(config.CONFIG_MAP["COINBASE_TRADING_PAIRS"], ",")
	channels := strings.Split(config.CONFIG_MAP["COINBASE_CHANNELS"], ",")

	return &CoinbaseStream{
		Conn: nil,
		Url:  url,
		SubscribeMsg: model.CoinbaseSubscribeMessage{
			Type:       "subscribe",
			ProductIds: productIds,
			Channels:   channels,
		},
		Retries: 10, // try to connect n times
	}
}

// connect to the endpoint
func (s *CoinbaseStream) Connect() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true

	// Try to connect n times
	for i := 0; i < s.Retries; i++ {
		c, _, err := dialer.Dial(s.Url, nil)
		if err != nil {
			// log.Println("Dial error:", err)

			// Exponential backoff: (2^(retry_attempt) - 1) / 2 * 1000 * milliseconds
			t := time.Duration((math.Pow(2, float64(i))-1)/2*1000) * time.Millisecond

			select {
			case <-sig:
				return errors.New("received interrupt signal")
			case <-time.After(t):
				// log.Printf("Trying to reconnect to %s (#%d/%d)\n", s.Url, i+1, s.Retries)
			}
		} else {
			s.Conn = c
			return nil
		}
	}

	return errors.New("could not connect")
}

// send a subscribe message
func (s *CoinbaseStream) Subscribe() error {
	return s.Conn.WriteJSON(s.SubscribeMsg)
}

func (s *CoinbaseStream) Receive(done chan interface{}, in chan<- model.CoinbaseResponseMessage) {
	defer close(done)

	m := model.CoinbaseResponseMessage{}

	for {
		if err := s.Conn.ReadJSON(&m); err != nil {
			println(err.Error())
			return
		}

		if m.Price != "" {
			in <- m
		}
	}
}

func (s *CoinbaseStream) Close() error {
	return s.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
