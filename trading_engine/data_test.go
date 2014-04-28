package trading_engine_test

import (
  "strconv"
  . "github.com/hookercookerman/trading_engine/trading_engine"
  . "github.com/onsi/ginkgo"
)

var _ = Describe("Data Test", func() {
  var orderBook *OrderBook
  var orders []*LimitOrder

  BeforeEach(func() {
    orders = []*LimitOrder{}
    orderBook = NewOrderBook()
    orderBook.Executer = &Executer{}
    for i, datum := range DATA {
      order := &LimitOrder{Price: datum.Price, Side: datum.Side, Id: strconv.Itoa(i + 1), Size: datum.Size}
      orders = append(orders, order)
    }
  })

  Measure("Performance from test Data", func(b Benchmarker) {
    b.Time("runtime", func() {
      for _, order := range orders {
        if order.Price == 0 {
          orderBook.Cancel(order.Id)
        } else {
          orderBook.Limit(order)
        }
      }
    })
  }, 100)

})
