package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("Single Buy Order", func() {

  var orderBook *OrderBook
  var tradeExecuter *Executer

  BeforeEach(func() {
    orderBook = NewOrderBook()
    tradeExecuter = &Executer{}
    orderBook.Executer = tradeExecuter
  })

  Context("#Limit, Given LowestAsk is 10", func() {
    var buyOrder *LimitOrder
    BeforeEach(func() {
      orderBook.LowestAsk = 10
      buyOrder = &LimitOrder{Symbol: "BTC", Size: 10, Price: 10, Id: "1", Side: BUY}
      orderBook.Limit(buyOrder)
    })

    It("trade executer should execute a trade", func() {
      Expect(tradeExecuter.called).To(BeZero())
    })

    It("Open Orders include the buy order", func() {
      Expect(orderBook.OpenOrders).To(ContainElement(buyOrder))
    })

    It("Price Points created for buy order price", func() {
      _, ok := orderBook.PricePoints.Get(buyOrder.Price)
      Expect(ok).To(BeTrue())
    })

    It("Book Entry for buy order is created", func() {
      value, _ := orderBook.PricePoints.Get(buyOrder.Price)
      pricePoint := value.(*PricePoint)
      Expect(len(pricePoint.BuyBookEntries)).ToNot(BeZero())
    })

    It("sets the order book HightestBid to the order price", func() {
      Expect(orderBook.HighestBid).To(Equal(buyOrder.Price))
    })
  })

})
