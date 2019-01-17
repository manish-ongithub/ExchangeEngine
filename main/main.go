package main

import (
	"context"
	"fmt"
	"time"

	"../engine"
	kafka "github.com/segmentio/kafka-go"
)

var book engine.OrderBook


func main() {
	// create the order book
	book = engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}

Start:
	/*
		fmt.Print("Enter text: ")
		var input string
		fmt.Scanln(&input)
		cmdHandler(input)
	*/
	done := make(chan bool)
	fmt.Println("create the consumer and listen for new order messages")

	go func() {
		fmt.Println("reader routine called")
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{"192.168.0.52:9092"},
			GroupID:        "consumer-group-id",
			Topic:          "orders",
			Partition:      0,
			MinBytes:       10e3,        // 10KB
			MaxBytes:       10e6,        // 10MB
			CommitInterval: time.Second, // flushes commits to Kafka every second
		})
		//r.SetOffset(42)

		for {
			fmt.Println("calling r.ReadMessage")
			ctx := context.Background()
			msg, err := r.ReadMessage(ctx)
			//msg, err := r.FetchMessage(ctx)
			if err != nil {
				fmt.Println(err)
				break
			}

			//fmt.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
			var order engine.Order
			order.FromJSON(msg.Value) // decode the message
			//r.CommitMessages(ctx, msg)
			if order.Price == 0 && order.Amount == 0 && order.Ctime == 0 { //fake message
				continue
			}

			fmt.Print("order => ")
			fmt.Println(order)
			// process the order
			trades := book.Process(order)
			for _, trade := range trades {
				rawTrade := trade.ToJSON()
				fmt.Print(" trade => ")
				fmt.Println(trade)

				topic := "trades"
				partition := 0
				fmt.Println("writing trade")
				conn, _ := kafka.DialLeader(context.Background(), "tcp", "192.168.0.52:9092", topic, partition)

				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				conn.WriteMessages(
					kafka.Message{Value: rawTrade},
				)
				fmt.Println("written trade and closing connection")
				conn.Close()

			}

		}

		fmt.Println("closing reader")
		r.Close()
		done <- true
	}()

	// wait until we are done
	<-done
	fmt.Println("going to start label")
	goto Start

}
