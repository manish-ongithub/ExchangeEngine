package engine

import (
	"sort"
)

type OrderSorter []Order

func (a OrderSorter) Len() int           { return len(a) }
func (a OrderSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderSorter) Less(i, j int) bool { return a[i].Price > a[j].Price }

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

func (book *OrderBook) addBuyOrder(order Order) {
	book.BuyOrders = append(book.BuyOrders, order)
	length := len(book.BuyOrders)
	if length < 2 {
		return
	}
	/*
		sort.Slice(book.BuyOrders, func(i, j int) bool {
			return book.BuyOrders[i].Price < book.BuyOrders[i].Price
		})
	*/
	sort.Sort(sort.Reverse(OrderSorter(book.BuyOrders)))

}

func (book *OrderBook) addSellOrder(order Order) {
	book.SellOrders = append(book.SellOrders, order)
	length := len(book.SellOrders)
	if length < 2 {
		return
	}

	sort.Sort(OrderSorter(book.SellOrders))
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
