module github.com/cgmello/marketdata/indicator

go 1.18

require (
	github.com/cgmello/marketdata/config v0.0.0-00010101000000-000000000000
	github.com/cgmello/marketdata/model v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cgmello/marketdata/model => ../model

replace github.com/cgmello/marketdata/config => ../../config
