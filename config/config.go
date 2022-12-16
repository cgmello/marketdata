package config

import "sync"

var (
	CONFIG_MAP = make(map[string]string)
)

var onceConfig sync.Once

func init() {
	onceConfig.Do(func() {
		// Get values from Secrets, .env file or hard-coded
		CONFIG_MAP["COINBASE_URL"] = "wss://ws-feed.exchange.coinbase.com"
		CONFIG_MAP["COINBASE_TRADING_PAIRS"] = "BTC-USD,ETH-USD,ETH-BTC"
		CONFIG_MAP["COINBASE_CHANNELS"] = "matches"
		CONFIG_MAP["COINBASE_SLIDING_WINDOW"] = "200"
	})
}
