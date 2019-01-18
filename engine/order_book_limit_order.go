package engine

import "fmt"

// Process an order and return the trades generated before adding the remaining amount to the market
func (book *OrderBook) Process(order Order) []Trade {
	if order.Side == 1 {
		fmt.Println("calling processLimitBuy")
		return book.processLimitBuy(order)
	}
	fmt.Println("calling processLimitSell")
	return book.processLimitSell(order)
}

func (book *OrderBook) processLimitBuy(order Order) []Trade {
	/*
		make([]int, 0, 1) allocates an underlying array
		of size 1 and returns a slice of length 0 and capacity 1 that is
		backed by this underlying array.
	*/
	trades := make([]Trade, 0, 1)
	l := len(book.SellOrders)
	fmt.Printf("inside processLimitBuy len of SellOrders -> %d\n", l)
	if l == 0 {
		book.addBuyOrder(order)
		return trades
	}
	if l != 0 || book.SellOrders[l-1].Price <= order.Price {
		// traverse all orders that match

		for i := l - 1; i >= 0; i-- {
			sellOrder := book.SellOrders[i]
			if sellOrder.Price < order.Price {
				break
			}
			// fill the entire order
			fmt.Printf("processlimitbuy ** for loop completed value of i => %d\n", i)
			if sellOrder.Amount >= order.Amount {
				trades = append(trades, Trade{order.ID, sellOrder.ID, order.Amount, sellOrder.Price})
				sellOrder.Amount -= order.Amount
				if sellOrder.Amount == 0 {
					book.removeSellOrder(i)
				}
				fmt.Print("")
				return trades
			}
			// fill a partial order and continue
			fmt.Println("*** filling partial order")
			if sellOrder.Amount < order.Amount {
				trades = append(trades, Trade{order.ID, sellOrder.ID, sellOrder.Amount, sellOrder.Price})
				order.Amount -= sellOrder.Amount
				book.removeSellOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addBuyOrder(order)

	return trades
}

func (book *OrderBook) processLimitSell(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	l := len(book.BuyOrders)
	fmt.Printf("*** inside processLimitSell len of BuyOrders -> %d\n", l)
	if l == 0 {
		book.addSellOrder(order)
		return trades
	}

	if l != 0 || book.BuyOrders[l-1].Price <= order.Price {
		// traverse all orders that match
		for i := l - 1; i >= 0; i-- {
			buyOrder := book.BuyOrders[i]
			if buyOrder.Price > order.Price {
				break
			}
			// fill the entire order
			fmt.Printf("*** prcesslimitsell for loop completed value of i => %d\n", i)
			if buyOrder.Amount >= order.Amount {
				trades = append(trades, Trade{order.ID, buyOrder.ID, order.Amount, buyOrder.Price})
				buyOrder.Amount -= order.Amount
				if buyOrder.Amount == 0 {
					book.removeBuyOrder(i)
				}
				return trades
			}
			// fill a partial order and continue
			fmt.Println("*** processlimitsell fill a partial order and continue")
			if buyOrder.Amount < order.Amount {
				trades = append(trades, Trade{order.ID, buyOrder.ID, buyOrder.Amount, buyOrder.Price})
				order.Amount -= buyOrder.Amount
				book.removeBuyOrder(i)
				continue
			}
		}
	}
	// finally add the remaining order to the list
	book.addSellOrder(order)
	return trades
}
