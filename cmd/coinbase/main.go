package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/cgmello/marketdata/config"
	"github.com/cgmello/marketdata/indicator"
	"github.com/cgmello/marketdata/model"
	"github.com/cgmello/marketdata/outputs"
)

var (
	interrupt chan os.Signal
)

func init() {
	log.Println("Starting..")

	// Channel to listen for interrupt signal to terminate gracefully
	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
}

func main() {
	// create a new stream
	s := NewCoinbaseStream()

	// Connect to the websockets stream
	if err := s.Connect(); err != nil {
		log.Fatal(err)
	}

	// Subscribe to channel
	if err := s.Subscribe(); err != nil {
		log.Fatal(err)
	}

	// Close connection when finished
	defer s.Conn.Close()

	// Comm channels
	done := make(chan interface{})
	in := make(chan model.CoinbaseResponseMessage)

	go func() {
		// Receive messages from stream
		s.Receive(done, in)
	}()

	slidingWindow, err := strconv.ParseInt(
		config.CONFIG_MAP["COINBASE_SLIDING_WINDOW"], 0, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Init VWAP data series
	vwap := indicator.NewVWAP(slidingWindow)

	message := model.CoinbaseResponseMessage{}

	for {
		select {

		case <-done:
			return

		case message = <-in:
			price, err1 := strconv.ParseFloat(message.Price, 64)
			qty, err2 := strconv.ParseFloat(message.Size, 64)
			if err1 == nil && err2 == nil {
				// New point received: trading pair, price and quantity
				p := model.NewPoint(message.ProductID, price, qty)

				// Process new point
				v := vwap.Process(p)

				// Send to the right output
				out := outputs.Stdout{
					Point: p,
					Name:  "VWAP",
					Value: v,
				}
				out.Send()
			}

		case <-interrupt:
			// Received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := s.Close()
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting..")
			case <-time.After(time.Duration(2) * time.Second):
				log.Println("Timeout in closing receiver channel. Exiting..")
			}
			return
		}
	}
}
