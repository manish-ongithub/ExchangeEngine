package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"../../engine"
	kafka "github.com/segmentio/kafka-go"
)

//pseudouuid
func pseudouuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func createOrder(amount uint64, price uint64, ordertype string) *engine.Order {
	ord := new(engine.Order)
	tm := time.Now().Unix()
	ord_id := pseudouuid()

	if ordertype == "buy" {
		fmt.Println("creating buy order")
		ord.ID = ord_id
		ord.Amount = amount
		ord.Price = price
		ord.Side = 1
		ord.Ctime = tm

	}
	if ordertype == "sell" {
		fmt.Println("creating sell order")
		ord.ID = ord_id
		ord.Amount = amount
		ord.Price = price
		ord.Side = 0
		ord.Ctime = tm

	}

	return ord
}
func cmdHandler(command string) {
	switch command {
	case "buy":
		fmt.Println("buy order")
		var amount uint64
		var price uint64
		fmt.Print("Enter price per quantity :")
		fmt.Scanln(&price)
		fmt.Print("Enter amount to purchase :")
		fmt.Scanln(&amount)

		order := createOrder(amount, price, "buy")
		fmt.Println(order)
		// create producer
		// to produce messages

		go func() {
			fmt.Println("writing buy order")
			topic := "orders"
			partition := 0

			conn, _ := kafka.DialLeader(context.Background(), "tcp", "192.168.0.52:9092", topic, partition)

			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			conn.WriteMessages(
				kafka.Message{Value: order.ToJSON()},
			)
			fmt.Println("conn closes after writing message")
			conn.Close()

		}()
	case "sell":
		fmt.Println("sell order")
		var amount uint64
		var price uint64
		fmt.Print("Enter price per quantity :")
		fmt.Scanln(&price)
		fmt.Print("Enter amount to sell :")
		fmt.Scanln(&amount)
		order := createOrder(amount, price, "sell")
		fmt.Println(order)
		go func() {
			fmt.Println("writing sell order")
			topic := "orders"
			partition := 0

			conn, _ := kafka.DialLeader(context.Background(), "tcp", "192.168.0.52:9092", topic, partition)

			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			conn.WriteMessages(
				kafka.Message{Value: order.ToJSON()},
			)
			fmt.Println("conn closes after writing message")
			conn.Close()

		}()

	case "orderlist":
		/*
			fmt.Println("buy order List")
			fmt.Print(len(book.BuyOrders))
			for _, order := range book.BuyOrders {
				fmt.Println(order)
			}
			fmt.Println("sell order List")
			fmt.Print(len(book.SellOrders))
			for _, order := range book.SellOrders {
				fmt.Println(order)
			}
		*/

	default:
		fmt.Println("invalid command")
	}
}
func main() {
	for {
		fmt.Print("Enter command: ")
		var input string
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}
		cmdHandler(input)
	}
}
