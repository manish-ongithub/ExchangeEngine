package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"../engine"
	kafka "github.com/segmentio/kafka-go"
)

var book engine.OrderBook
var tradelist []engine.Trade

//PrintCurrentBuySellOrders : prints list of current orders in queue
func PrintCurrentBuySellOrders() {
	lb := len(book.BuyOrders)

	fmt.Printf("BUY ORDERS LIST (%d)\n", lb)
	fmt.Println("________________________________________________________")
	fmt.Printf(" ID                               \tPrice\tAmount\tTime\n")

	for i := 0; i < lb; i++ {
		order := book.BuyOrders[i]
		fmt.Printf("%s\t%d\t%d\t%d\n", order.ID, order.Price, order.Amount, order.Ctime)
	}
	ls := len(book.SellOrders)

	fmt.Printf("SELL ORDERS LIST (%d)\n", ls)
	fmt.Println("________________________________________________________")
	fmt.Printf(" ID                               \tPrice\tAmount\tTime\n")

	for i := 0; i < ls; i++ {
		order := book.SellOrders[i]
		fmt.Printf("%s\t%d\t%d\t%d\n", order.ID, order.Price, order.Amount, order.Ctime)
	}
}

//StartWebServer : this function starts the webserver
func StartWebServer() {
	fmt.Println("starting httpserver on port 8000")
	http.HandleFunc("/orderlist", func(w http.ResponseWriter, r *http.Request) {

		ls := len(book.SellOrders)
		s := fmt.Sprintf("SELL ORDERS LIST (%d)\n", ls)
		w.Write([]byte(s))
		w.Write([]byte(" ID                               \tPrice\tAmount\tTime\n"))

		for i := 0; i < ls; i++ {
			order := book.SellOrders[i]
			str := fmt.Sprintf("%s\t%d\t%d\t%d\n", order.ID, order.Price, order.Amount, order.Ctime)
			w.Write([]byte(str))
		}

		lb := len(book.BuyOrders)
		w.Write([]byte("\n\n"))
		s = fmt.Sprintf("BUY ORDERS LIST (%d)\n", lb)
		w.Write([]byte(s))
		w.Write([]byte(" ID                               \tPrice\tAmount\tTime\n"))

		for i := 0; i < lb; i++ {
			order := book.BuyOrders[i]
			s = fmt.Sprintf("%s\t%d\t%d\t%d\n", order.ID, order.Price, order.Amount, order.Ctime)
			w.Write([]byte(s))
		}
		lt := len(tradelist)
		w.Write([]byte("\n\n"))
		s = fmt.Sprintf("TRADE LIST (%d)\n", lb)
		w.Write([]byte(s))
		w.Write([]byte(" Taker ID                               \t Maker ID                               \t"))
		w.Write([]byte("Price\tAmount\tTime\n"))

		for i := 0; i < lt; i++ {
			t := tradelist[i]
			s = fmt.Sprintf("%s\t%s\t%d\t%d\n", t.TakerOrderID, t.MakerOrderID, t.Amount, t.Price)
			w.Write([]byte(s))
		}

	})
	http.ListenAndServe(":8000", nil)
}
func main() {
	// create the order book
	book = engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 100),
		SellOrders: make([]engine.Order, 0, 100),
	}
	tradelist = make([]engine.Trade, 0, 200)
	go StartWebServer()
Start:

	done := make(chan bool)
	tradechannel := make(chan bool)

	go func() {
		fmt.Println("Starting orders topic consumer ")
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
			PrintCurrentBuySellOrders()
			for _, trade := range trades {
				tradelist = append(tradelist, trade)
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
	go func() {
		fmt.Println("Starting trades topic consumer ")
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{"192.168.0.52:9092"},
			GroupID:        "trade-consumer-group-id",
			Topic:          "trades",
			Partition:      0,
			MinBytes:       10e3,        // 10KB
			MaxBytes:       10e6,        // 10MB
			CommitInterval: time.Second, // flushes commits to Kafka every second
		})
		for {
			fmt.Println("calling trades r.ReadMessage")
			ctx := context.Background()
			msg, err := r.ReadMessage(ctx)
			//msg, err := r.FetchMessage(ctx)
			if err != nil {
				fmt.Println(err)
				break
			}
			var trade engine.Trade
			trade.FromJSON(msg.Value)
			fmt.Print("trade reader ->")
			fmt.Println(trade)

		}
		tradechannel <- true
	}()

	// wait until we are done
	<-done
	<-tradechannel
	fmt.Println("going to start label")
	goto Start

}
