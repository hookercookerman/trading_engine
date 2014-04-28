package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Full Sell Order", func() {

  var orderBook *OrderBook
  var tradeExecuter *Executer
  var buyOrder *LimitOrder

  BeforeEach(func() {
    orderBook = NewOrderBook()
    tradeExecuter = &Executer{}
    orderBook.Executer = tradeExecuter
    buyOrder = &LimitOrder{Price: 10, Size: 10, Side: BUY, Id: "1"}
    orderBook.Limit(buyOrder)
  })

  Context("Sell order matching price and size of buy order", func() {
    var sellOrder *LimitOrder
    BeforeEach(func() {
      sellOrder = &LimitOrder{Price: 10, Size: 10, Side: SELL, Id: "2"}
      orderBook.Limit(sellOrder)
    })

    It("executes 1 trade", func() {
      Expect(tradeExecuter.called).To(Equal(1))
    })

    It("clears the buy Open Order", func() {
      _, ok := orderBook.OpenOrders[buyOrder.Id]
      Expect(ok).To(BeFalse())
    })

    It("clears the pricePoint for the orders price point", func() {
      _, ok := orderBook.PricePoints.Get(sellOrder.Price)
      Expect(ok).To(BeFalse())
    })

    It("sell order should not be in the Open Orders", func() {
      _, ok := orderBook.OpenOrders[sellOrder.Id]
      Expect(ok).To(BeFalse())

    })
  })

  Context("Another Buy Order on same price point", func() {
    var buyOrder2 *LimitOrder
    var sellOrder *LimitOrder
    BeforeEach(func() {
      buyOrder2 = &LimitOrder{Price: 10, Size: 10, Side: BUY, Id: "2"}
      orderBook.Limit(buyOrder2)
    })

    Describe("#Limit sell order crossing both buy orders", func() {
      BeforeEach(func() {
        sellOrder = &LimitOrder{Price: 10, Size: 20, Side: SELL, Id: "3"}
        orderBook.Limit(sellOrder)
      })

      It("executes 2 trades", func() {
        Expect(tradeExecuter.called).To(Equal(2))
      })

      It("clears all orders", func() {
        Expect(len(orderBook.OpenOrders)).To(BeZero())
      })

      It("clear the pricepoint for that price", func() {
        _, ok := orderBook.PricePoints.Get(sellOrder.Price)
        Expect(ok).To(BeFalse())
      })
    })

    Describe("#Limit sell order Price 10, Size: 14 crossing partially across the 2 buy orders", func() {
      BeforeEach(func() {
        sellOrder = &LimitOrder{Price: 10, Size: 14, Side: SELL, Id: "3"}
        orderBook.Limit(sellOrder)
      })

      It("executes 2 trades", func() {
        Expect(tradeExecuter.called).To(Equal(2))
      })

      It("keeps the second buy order on its open orders", func() {
        _, ok := orderBook.OpenOrders[buyOrder2.Id]
        Expect(ok).To(BeTrue())
      })

      It("keeps the price point 10 as there are still book entries there", func() {
        _, ok := orderBook.PricePoints.Get(sellOrder.Price)
        Expect(ok).To(BeTrue())
      })

      It("has only 1 bookEntry for the PricePoint", func() {
        value, _ := orderBook.PricePoints.Get(sellOrder.Price)
        pricePoint := value.(*PricePoint)
        Expect(len(pricePoint.BuyBookEntries)).To(Equal(1))
        Expect(len(pricePoint.SellBookEntries)).To(Equal(0))
      })

      It("has a bookEntry for the Buy Order with adjusted amount to 6", func() {
        value, _ := orderBook.PricePoints.Get(sellOrder.Price)
        pricePoint := value.(*PricePoint)
        bookEntry := pricePoint.BuyBookEntries[0]
        Expect(bookEntry.Size).To(Equal(uint64(6)))
      })
    })
  })

})
