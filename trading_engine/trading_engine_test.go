package trading_engine_test

import (
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
  // . "github.com/onsi/gomega"
)

type Executer struct {
  called           int
  Order            *LimitOrder
  MatchedBookEntry *BookEntry
  Price            uint64
  returnError      error
  returnOk         bool
}

func (self *Executer) Execution(price uint64, fullOrder bool, order *LimitOrder, matchedBookEntry *BookEntry) (error, bool) {
  self.called += 1
  self.Order = order
  self.Price = price
  self.MatchedBookEntry = matchedBookEntry
  return self.returnError, self.returnOk
}

var _ = Describe("TradingEngine", func() {
  var orderBook *OrderBook
  BeforeEach(func() {
    orderBook = NewOrderBook()
  })

  Context("#Limit, Given BUY Order with Price 12", func() {
    limitOrder := new(LimitOrder)
    BeforeEach(func() {
      limitOrder.Id = "1"
      limitOrder.Symbol = "BTC"
      limitOrder.Price = 12
      limitOrder.Side = 0
      limitOrder.Size = 10
    })
  })
})

//     Context("Matching Sell Order with same pice and size", func() {
//       var sellOrder *LimitOrder
//       var executer *Executer

//       BeforeEach(func() {
//         executer = &Executer{}
//         orderBook.Executer = executer
//         sellOrder = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: 1, Size: 10}
//         orderBook.Limit(sellOrder)
//       })

//       It("executes the trade as we have a matching sell and buy", func() {
//         Expect(executer.called).To(Equal(1))
//       })

//       It("deletes both OpenOrders", func() {
//         Expect(orderBook.OpenOrders).To(BeEmpty())
//       })

//       It("deletes the price Point as there are now no entries", func() {
//         _, ok := orderBook.PricePoints.Get(int64(12))
//         Expect(ok).To(BeFalse())
//       })
//     })

//     Context("Two Buy OpenOrders same Price Point", func() {
//       var buyOrder2 *LimitOrder
//       var executer *Executer

//       BeforeEach(func() {
//         executer = &Executer{}
//         orderBook.Executer = executer
//         buyOrder2 = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: 0, Size: 5}
//         orderBook.Limit(buyOrder2)
//       })

//       It("should have 2 BuyOrders for that price point", func() {
//         value, _ := orderBook.PricePoints.Get(int64(12))
//         pricePoint := value.(*PricePoint)
//         Expect(len(pricePoint.BuyBookEntries)).To(Equal(2))
//       })

//       Context("Sell Order that is more then the current Buy Price", func() {
//         var sellOrder *LimitOrder
//         var executer *Executer
//         BeforeEach(func() {
//           executer = &Executer{}
//           orderBook.Executer = executer
//           sellOrder = &LimitOrder{Id: "3", Symbol: "BTC", Price: 12, Side: 1, Size: 20}
//           orderBook.Limit(sellOrder)
//         })

//         It("should execute 2 Trades", func() {
//           Expect(executer.called).To(Equal(2))
//         })

//         It("should delete the original BookEntry For the Buy", func() {
//           Expect(orderBook.OpenOrders[limitOrder.Id]).To(BeNil())
//           Expect(orderBook.OpenOrders[buyOrder2.Id]).To(BeNil())
//         })
//       })

//     })

//     Context("Sell Order that is more then the current Buy Price", func() {
//       var sellOrder *LimitOrder
//       var executer *Executer
//       BeforeEach(func() {
//         executer = &Executer{}
//         orderBook.Executer = executer
//         sellOrder = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: 1, Size: 20}
//         orderBook.Limit(sellOrder)
//       })

//       It("should execute an order", func() {
//         Expect(executer.called).To(Equal(1))
//       })

//       It("should delete the original BookEntry For the Buy", func() {
//         Expect(orderBook.OpenOrders[limitOrder.Id]).To(BeNil())
//       })
//     })

//     Context("Sell Order that only partially matches forfills the Buy Order", func() {
//       var sellOrder *LimitOrder
//       var executer *Executer

//       BeforeEach(func() {
//         executer = &Executer{}
//         orderBook.Executer = executer
//         sellOrder = &LimitOrder{Id: "2", Symbol: "BTC", Price: 12, Side: 1, Size: 1}
//         fmt.Println(orderBook.OpenOrders)
//         orderBook.Limit(sellOrder)
//         fmt.Println(orderBook.OpenOrders)
//       })

//       It("should execute a trade", func() {
//         Expect(executer.called).To(Equal(1))
//       })

//       It("should not delete the Buy Sell Price Point", func() {
//         _, ok := orderBook.PricePoints.Get(sellOrder.Price)
//         Expect(ok).To(BeTrue())
//       })

//       It("should delete the sell order as it has been fully forfilled", func() {
//         Expect(orderBook.OpenOrders[sellOrder.Id]).To(BeNil())
//       })

//       It("should delete the Buy Book Entry as it its not fully exhasted", func() {
//         value, _ := orderBook.PricePoints.Get(int64(12))
//         pricePoint := value.(*PricePoint)
//         Expect(pricePoint.BuyBookEntries[0].Size).To(Equal(int64(9)))
//       })

//     })

//   })
// })

// })
