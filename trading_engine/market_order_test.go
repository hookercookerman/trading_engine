package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

var _ = Describe("TradingEngine, #Market", func() {
  var tradeExecuter *Executer

  Context("When sell and buy order exits in the marker", func() {
    var sellOrder1 *LimitOrder
    var sellOrder2 *LimitOrder
    var sellOrder3 *LimitOrder
    var orderBook *OrderBook

    BeforeEach(func() {
      tradeExecuter = &Executer{}
      orderBook = NewOrderBook()
      orderBook.Executer = tradeExecuter
      sellOrder1 = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: SELL, Size: 5}
      sellOrder2 = &LimitOrder{Id: "3", Symbol: "BTC", Price: 8, Side: SELL, Size: 15}
      sellOrder3 = &LimitOrder{Id: "4", Symbol: "BTC", Price: 9, Side: SELL, Size: 50}
      orderBook.Limit(sellOrder1)
      orderBook.Limit(sellOrder2)
      orderBook.Limit(sellOrder3)
    })

    Describe("#MarketOrder", func() {
      var marketOrder *LimitOrder
      BeforeEach(func() {
        marketOrder = &LimitOrder{Id: "5", Symbol: "BTC", Side: BUY, Size: 10}
        orderBook.MarketOrder(marketOrder)
      })

      It("executes an order at the current LowestAsk on the order book", func() {
        Expect(tradeExecuter.Order.Price).To(Equal(orderBook.LowestAsk))
      })

      It("executes a trade with a bookEntry of the lowest Ask reduce to 5", func() {
        Expect(tradeExecuter.MatchedBookEntry.Size).To(Equal(uint64(5)))
      })

      It("executes an order at the current LowestAsk", func() {
        Expect(tradeExecuter.Price).To(Equal(orderBook.LowestAsk))
      })

    })
  })
})
