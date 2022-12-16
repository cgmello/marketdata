# VWAP calculation engine for Go

[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)


The quick goal of this project is to develop a [VWAP (Volume-Weighted Average Price)](https://en.wikipedia.org/wiki/Volume-weighted_average_price) calculation engine for Golang that uses the [Coinbase's websockets feed](https://docs.cloud.coinbase.com/exchange/docs/websocket-overview) to stream in trade executions and update the VWAP values for each trading pair as new prices and quantities become available.

## Marketdata
The big goal of this project would be to create a scalable market data platform to collect real-time data and calculate indicators. We created an architecture that should be ready to evolve, easy to maintain and already organized to become a production platform.

### Usage (Linux and MacOS)

To run this project you need [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) and [Golang 1.18 or newer](https://go.dev/doc/devel/release#policy):

```bash
# clone the repository
git clone https://github.com/cgmello/marketdata.git

# run
cd marketdata
make start
```

To run the unity and integration tests:

```bash
make tests
```

### For Windows users

To use **make** you need to install [Chocolatey](https://chocolatey.org/install) and then run ```choco install make```. 

If not, you can run the project directly:

```bash
# clone the repository
git clone https://github.com/cgmello/marketdata.git

# run
cd marketdata
cd cmd/coinbase
go mod tidy
go run .
```

And to run the tests:

```bash
go test -v -cover
cd ../../internal/indicator 
go mod tidy
go test -v -cover
```

## Features

- Ingests real-time data from Coinbase websockets feed available at the URL ```wss://ws-feed.exchange.coinbase.com``` for the **matches** channel subscription
- Calculates the real-time **VWAP** (Volume-Weighted Average Price) indicator as new data points become available (price and volume) from the websockets stream
- The VWAP calculation uses a sliding window of data points -  for now we use 200 points but you change that at the configurations. Meaning that when we have 200 points and a new point arrives, the oldest one will be removed from the calculation, and the new one enters, such that we always have the only and last 200 points for the VWAP's calculation
- Pulls data for any trading pair that is available at Coinbase, for this project we use the following trading pairs: **BTC-USD, ETH-USD and ETH-BTC**, but you can also change that in the configurations
- Prints the VWAP real-time values in the standard output (**stdout**), put we could send data to files or brokers
- Uses the library [Gorilla Websockets for Go](https://github.com/gorilla/websocket) for connecting to the Coinbase's websockets stream
- Uses the [Decimal](https://github.com/shopspring/decimal) library to avoid float precision issues
- Uses [Testify](https://github.com/stretchr/testify) to run unity and integration tests
- You can stop the execution pressing **CTRL-C**

## How this project was designed

This project was designed with 2 goals in mind:

1. Create an architecture that would be ready to grow into a production platform; and
2. Use the architecture to implement 1 ingestion service and 1 indicator as a POC small project

## Assumptions

We made some assumptions when implementing this code:

- The first 200 updates will have less than 200 data points included in the calculation
- In order to gain performance and storage efficiency we always keep the current sums of ```price*quantity```and ```quantity``` for each trading pair, so that when a new Point(p,q) arrives, we remove the oldest ```p0*q0``` value from the numerator and ```q0``` value from the denominator, and add the new related values ```p*q``` and ```q```
> VWAP equation = sum(price*quantity) / sum(quantity)
- We were aware of the [Coinbase's websocket feed rate limits](https://docs.cloud.coinbase.com/exchange/docs/websocket-rate-limits) but we have only 1 websocket connection at a time and we only send 1 subscription message, so it is not an issue
- In [Coinbase's websocket feed best practices](https://docs.cloud.coinbase.com/exchange/docs/websocket-best-practices) they suggest that we should split the connections for each subscription, but for efficiency purposes we choose to run only 1 connection. Maybe we should make some tests to check if it would be necessary to split the connections when we have a big number of trading pairs
- We choose not to Authenticate as it is optional and because we are not executing trades
- We set the compression settings ON when using the Gorilla websockets to achieve lower bandwidth consumption
- We used the original channel, not the **_batch** one, because we need real-time values

#### Folder structure

We suggested the following directory structure to handle the complete version of the project, so that we start with a production ready folder organization from day-1.

> **/** - in the root we have the Makefile and other administrative stuff. We could have here the docker-compose yaml file for example.

> **/api** - all kind of REST, Websockets and GRPC public (and private) API's specifications and related documents such as Postman's or Thunder Client's collections and online documentations such as Swaggo. We didn't implemented this for now.

> **/cmd** - the main folder of the project. Here we have the apps. In each app we can have the main file, auxiliar code, api routes, Dockerfile and test files (e.g. main_test.go). Here we can store our microservices and BFF (backend-for-frontend) codes. We have our coinbase app stored here.

> **/config** - configuration module. We use hard-coded values for now. Here we could set the config to access different resources depending upon the environment:
    - *.env* file for local development;
    - SQL/NoSQL databases or [Kubernetes etcd](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/) for non-sensitive data; and
    - [AWS Secrets Manager](https://aws.amazon.com/pt/secrets-manager/), [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) or some vault such as [HashiCorp](https://www.vaultproject.io) to store more sensitive information like credentials, tokens and passwords

> **/deployments** - Kubernetes deployments and services files, Terraform files, Helm chart files, etc. We didn't implemented this for now.

> **/internal** - every code that are common to our apps, but cannot be imported from external projects (we can have a /pkg folder for that) such as the /indicator (our VWAP), the /lib  (our stdout output) and the /model (our coinbase structs) folders. In the /lib folder we can put all the adapters that we should need to access external resources such as:
    - database adapters for MySQL, PostgreSQL, Timestream, InfluxDB, DynamoDB, Redis, etc.
    - stream adapters for Kinesis Data Streams, Data Firehose, etc.
    - blockchain adapters for our smart contracts
    - for S3 buckets or similar long-term storage
    - for external adapters such as e-mail providers.

> **/scripts** - shell scripts or cronjobs that are not running on Github Actions or Kubernetes. It's empty for now.

> **/tools** - tools such as database migration sql scripts. It's empty for now.

> **/website** - UI frontend apps that uses the backend code. We don't have any UI for now.

## Other development features

Some other features that we used on the development process, but that are not integrated automatically (yet):

- We run [Staticcheck](https://staticcheck.io) manually on our terminal as a linter for the Go programming language, besides the lint plugin that also checks our code on real-time on Visual Source Code, in order to improve the overall quality of code
- We run [Snyk CLI](https://docs.snyk.io/snyk-cli/getting-started-with-the-cli) manually to scan our project against known vulnerabilities

### Next steps

From this VWAP calculation engine quick project that was a POC to validate our architecture, we can list some possible next steps not necessarily in order of priority:

1. Implement some **DevOps** workflow to handle all kind of automatic tests and deployments for our project. We could use Github Actions for CI/CD pipelines or other similar tools like Jenkins
2. Create a **Kubernetes cluster** to run our microservices and some tools like Prometheus, Grafana and Alerts to have a more high availability, secure and monitored system
3. Design an **Event-driven** or Event-source architectures for example, to run our microservices using the best practices for system designing. We could use for example a Kinesis Data Stream instance to receive the real-time data, microservices to enrich and process data, microservices to store data in SQL/NoSQL databases and cache-like databases such as Redis, send data to a SQS/SNS solution to work as brokers and publish/subscribe options, implement public (or private) APIs to give access to the datalake to the own company and/or clients, etc.
4. Generate domains for our different environments (develop, staging, production). Configure AWS Route 53 to handle our routes and AWS Certification Manager to use https in our Load Balancers that handles our Kubernetes services. Enable AWS WAF for firewall protection. And so on.
