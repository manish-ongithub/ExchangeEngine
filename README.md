# ExchangeEngine
A simple trading engine built in Go language. Needs Kafka for reading and writing messages.
This has been tested on Windows 10 (Go 1.11.14)

It uses "github.com/segmentio/kafka-go" library for communication with kafka


ExchangeEngine/main/main.go is the main server which reads and processes oders and create trades.

ExchangeEngince/client/main.go can be used to test and puch orders to the server.

Both client and server use Kafka for reading and writing messages.


RIGHT NOW EVERYTHING IS IN RAW FORM. Will be updated shortly.
