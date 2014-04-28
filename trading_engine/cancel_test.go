package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("TradingEngine, #Cancel", func() {
  var orderBook *OrderBook
  var sellOrder *LimitOrder
  BeforeEach(func() {
    orderBook = NewOrderBook()
    sellOrder = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: 1, Size: 10}
  })

  Context("Given Order has been made", func() {
    BeforeEach(func() {
      orderBook.Limit(sellOrder)
      orderBook.Cancel(sellOrder.Id)
    })

    It("should not have contain canceled order", func() {
      Expect(orderBook.OpenOrders).ToNot(ContainElement(sellOrder))
    })

    It("clears the price point for that order", func() {
      _, ok := orderBook.PricePoints.Get(sellOrder.Price)
      Expect(ok).To(BeFalse())
    })
  })

  Context("Given Order already exist for same price point", func() {
    var anotherSellOrder *LimitOrder
    BeforeEach(func() {
      anotherSellOrder = &LimitOrder{Id: "4", Symbol: "BTC", Price: 12, Side: 1, Size: 10}
      orderBook.Limit(anotherSellOrder)
      orderBook.Limit(sellOrder)
      orderBook.Cancel(sellOrder.Id)
    })

    It("clears the book Entry for the order", func() {
      value, _ := orderBook.PricePoints.Get(sellOrder.Price)
      pricePoint := value.(*PricePoint)
      Expect(len(pricePoint.SellBookEntries)).To(Equal(1))
    })

    It("keeps the entry for the other sell order", func() {
      value, _ := orderBook.PricePoints.Get(sellOrder.Price)
      pricePoint := value.(*PricePoint)
      Expect(pricePoint.SellBookEntries).To(ContainElement(anotherSellOrder.BookEntry))
    })

  })
})
