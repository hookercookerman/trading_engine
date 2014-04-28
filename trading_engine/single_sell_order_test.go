package trading_engine_test

import (
  "fmt"
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Single Sell Order", func() {

  var orderBook *OrderBook
  var tradeExecuter *Executer

  BeforeEach(func() {
    orderBook = NewOrderBook()
    tradeExecuter = &Executer{}
    orderBook.Executer = tradeExecuter
  })

  Context("#Limit", func() {
    var sellOrder *LimitOrder
    BeforeEach(func() {
      sellOrder = &LimitOrder{Symbol: "BTC", Size: 10, Price: 10, Id: "1", Side: SELL}
      orderBook.Limit(sellOrder)
    })

    It("trade executer does not execute a trade", func() {
      Expect(tradeExecuter.called).To(BeZero())
    })

    It("Open Orders includes the sell order", func() {
      Expect(orderBook.OpenOrders).To(ContainElement(sellOrder))
    })

    It("Price Points created for sell order price", func() {
      _, ok := orderBook.PricePoints.Get(sellOrder.Price)
      Expect(ok).To(BeTrue())
    })

    It("Book Entry for sell order is created", func() {
      value, _ := orderBook.PricePoints.Get(sellOrder.Price)
      pricePoint := value.(*PricePoint)
      Expect(len(pricePoint.SellBookEntries)).ToNot(BeZero())
    })

    It("sets order price to the new lowestAsk", func() {
      fmt.Println("LOWEST", orderBook.LowestAsk)
      fmt.Println(" LOWEST", sellOrder.Price)
      Expect(orderBook.LowestAsk).To(Equal(sellOrder.Price))
    })
  })
})
