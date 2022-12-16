package main

import (
	"github.com/cgmello/marketdata/config"
	"github.com/cgmello/marketdata/model"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestConfiguration(t *testing.T) {
	// Check hard-coded vars
	require.Equal(t, config.CONFIG_MAP["COINBASE_URL"], "wss://ws-feed.exchange.coinbase.com")
	require.Equal(t, config.CONFIG_MAP["COINBASE_TRADING_PAIRS"], "BTC-USD,ETH-USD,ETH-BTC")
	require.Equal(t, config.CONFIG_MAP["COINBASE_CHANNELS"], "matches")
	require.Equal(t, config.CONFIG_MAP["COINBASE_SLIDING_WINDOW"], "200")

	// Check sliding window integer value
	win, err := strconv.ParseInt(config.CONFIG_MAP["COINBASE_SLIDING_WINDOW"], 0, 64)
	require.NoError(t, err)
	require.Equal(t, win, int64(200))
}

func TestConnShouldSuccess(t *testing.T) {
	s := NewCoinbaseStream()

	// Check websockets connection
	require.NoError(t, s.Connect())
	defer s.Conn.Close()
}

func TestConnShouldFail(t *testing.T) {
	s := NewCoinbaseStream()
	s.Url = "wss://some.fake.url.com" // fake url
	s.Retries = 3                     // retry just 3 times to not take so long

	// Check websockets connection, should fail
	require.Error(t, s.Connect())
}

func TestSubscribe(t *testing.T) {
	s := NewCoinbaseStream()

	// Check subscribing
	require.NoError(t, s.Connect())
	require.NoError(t, s.Subscribe())
	defer s.Close()
}

func TestReceiveMessages(t *testing.T) {
	s := NewCoinbaseStream()

	require.NoError(t, s.Connect())
	require.NoError(t, s.Subscribe())
	defer s.Close()

	done := make(chan interface{})
	in := make(chan model.CoinbaseResponseMessage)

	// Check receiving messages in a new goroutine
	go func() {
		s.Receive(done, in)
	}()

	message := model.CoinbaseResponseMessage{}

	maxPoints := 30 // check for the first 30 points
	n := 0

	for {
		select {

		case <-done: // signal from receive
			return

		case message = <-in:
			// Check if price is valid
			_, err := strconv.ParseFloat(message.Price, 64)
			require.NoError(t, err)

			// Check if quantity is valid
			_, err = strconv.ParseFloat(message.Size, 64)
			require.NoError(t, err)

			n++
			if n == maxPoints {
				return
			}

		case <-interrupt:
			require.NoError(t, s.Close())

			select {
			case <-done: // signal from receive
			case <-time.After(time.Duration(2) * time.Second): // timeout
			}
			return
		}
	}
}
