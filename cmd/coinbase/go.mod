module github.com/cgmello/marketdata/coinbase

go 1.18

replace github.com/cgmello/marketdata/config => ../../config

replace github.com/cgmello/marketdata/model => ../../internal/model

replace github.com/cgmello/marketdata/indicator => ../../internal/indicator

require (
	github.com/cgmello/marketdata/config v0.0.0-00010101000000-000000000000
	github.com/cgmello/marketdata/indicator v0.0.0-00010101000000-000000000000
	github.com/cgmello/marketdata/model v0.0.0-00010101000000-000000000000
	github.com/cgmello/marketdata/outputs v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.0
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/tools v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cgmello/marketdata/outputs => ../../internal/lib/outputs
