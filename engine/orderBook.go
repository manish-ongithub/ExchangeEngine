package engine

// OrderBook type
type OrderBook struct {
	BuyOrders  []Order
	SellOrders []Order
}

func (book *OrderBook) addBuyOrder(order Order) {
	l := len(book.BuyOrders)
	if l == 0 {
		book.BuyOrders = append(book.BuyOrders, order)
		return
	}

	var counter int
	for counter = l - 1; counter >= 0; counter-- {
		buyOrder := book.BuyOrders[counter]
		if buyOrder.Price < order.Price {
			break
		}
	}

	if counter == l-1 {
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		copy(book.BuyOrders[counter+1:], book.BuyOrders[counter:]) // copy(dst, src []Type) int
		book.BuyOrders[counter] = order
	}
}

func (book *OrderBook) addSellOrder(order Order) {
	len := len(book.SellOrders)
	if len == 0 {
		book.SellOrders = append(book.SellOrders, order)
		return
	}
	var counter int
	for counter = len - 1; counter >= 0; counter-- {
		sellOrder := book.SellOrders[counter]
		if sellOrder.Price > order.Price {
			break
		}
	}
	if counter == len-1 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		copy(book.SellOrders[counter+1:], book.SellOrders[counter:])
		book.SellOrders[counter] = order
	}
}

// Remove a buy order from the order book at a given index
func (book *OrderBook) removeBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

// Remove a sell order from the order book at a given index
func (book *OrderBook) removeSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
