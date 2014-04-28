package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("EmptyOrderBook", func() {

  var orderBook *OrderBook
  BeforeEach(func() {
    orderBook = NewOrderBook()
  })

  Context("New Order Book", func() {

    It("should have empty OpenOrders", func() {
      Expect(orderBook.OpenOrders).To(BeEmpty())
    })

    It("should have empty price points", func() {
      Expect(orderBook.PricePoints.Len()).To(BeZero())
    })

    It("should have a Lowest Ask of 0", func() {
      Expect(orderBook.LowestAsk).To(Equal(uint64(0)))
    })

    It("should have a Highest Bid of 0", func() {
      Expect(orderBook.HighestBid).To(Equal(uint64(0)))
    })

  })
})
