package trading_engine

import "github.com/ryszard/goskiplist/skiplist"

type TradeExecuter interface {
  Execution(price uint64, fullOrder bool, order *LimitOrder, matchedBookEntry *BookEntry) (error, bool)
}

type BookEntry struct {
  Size  uint64
  Side  int
  Order *LimitOrder
}

type PricePoint struct {
  SellBookEntries []*BookEntry
  BuyBookEntries  []*BookEntry
}

type OrderBook struct {
  PricePoints *skiplist.SkipList
  OpenOrders  map[string]*LimitOrder
  Executer    TradeExecuter
  LowestAsk   uint64
  HighestBid  uint64
}

func NewSkipList() *skiplist.SkipList {
  return skiplist.NewCustomMap(func(l, r interface{}) bool {
    return l.(uint64) < r.(uint64)
  })
}

func NewOrderBook() *OrderBook {
  pricePoints := NewSkipList()
  orders := make(map[string]*LimitOrder)
  return &OrderBook{PricePoints: pricePoints, OpenOrders: orders, LowestAsk: uint64(0), HighestBid: uint64(0)}
}

func (self *OrderBook) Cancel(id string) {
  if order, ok := self.OpenOrders[id]; ok {
    if value, ok := self.PricePoints.Get(order.Price); ok {
      self.deleteBookEntry(value.(*PricePoint), order)
      self.deletePricePoint(value.(*PricePoint), order.Price)
    }
    delete(self.OpenOrders, id)
  }
}

func (self *OrderBook) deleteBookEntry(pricePoint *PricePoint, order *LimitOrder) {
  if order.BookEntry != nil {
    if order.Side == BUY {
      for i, bookEntry := range pricePoint.BuyBookEntries {
        if bookEntry.Order == order {
          pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries[:i], pricePoint.BuyBookEntries[i+1:]...)
          break
        }
      }
    } else {
      for i, bookEntry := range pricePoint.SellBookEntries {
        if bookEntry.Order == order {
          pricePoint.SellBookEntries = append(pricePoint.SellBookEntries[:i], pricePoint.SellBookEntries[i+1:]...)
          break
        }
      }
    }
  }
}

func (self *OrderBook) deletePricePoint(pricePoint *PricePoint, price uint64) {
  if len(pricePoint.BuyBookEntries) == 0 && len(pricePoint.SellBookEntries) == 0 {
    self.PricePoints.Delete(price)
  }
}

func (self *OrderBook) deleteOrder(order *LimitOrder) {
  delete(self.OpenOrders, order.Id)
}

func (self *OrderBook) createPricePoint(price uint64, bookEntry *BookEntry) {
  if value, ok := self.PricePoints.Get(price); ok {
    pricePoint := value.(*PricePoint)
    if bookEntry.Side == BUY {
      pricePoint.BuyBookEntries = append(pricePoint.BuyBookEntries, bookEntry)
    } else {
      pricePoint.SellBookEntries = append(pricePoint.SellBookEntries, bookEntry)
    }
  } else {
    buyBookEntries := []*BookEntry{}
    sellBookEntries := []*BookEntry{}
    if bookEntry.Side == BUY {
      buyBookEntries = append(buyBookEntries, bookEntry)
    } else {
      sellBookEntries = append(sellBookEntries, bookEntry)
    }
    pricePoint := &PricePoint{BuyBookEntries: buyBookEntries, SellBookEntries: sellBookEntries}
    self.PricePoints.Set(price, pricePoint)
  }
}

func (self *OrderBook) createBookEntryForOrder(order *LimitOrder) *BookEntry {
  bookEntry := &BookEntry{Size: order.Size, Side: order.Side, Order: order}
  order.BookEntry = bookEntry
  return bookEntry
}

func (self *OrderBook) limitSell(order *LimitOrder) {
  if order.Price <= self.HighestBid && self.HighestBid != 0 {
    finished := false
    iterator := self.PricePoints.Seek(self.HighestBid)
    defer iterator.Close()
    for iterator != nil && !finished && order.Price <= self.HighestBid {
      pricePoint := iterator.Value().(*PricePoint)
      for _, bookEntry := range pricePoint.BuyBookEntries {
        if bookEntry.Size < order.Size {
          self.Execution(iterator.Key().(uint64), false, order, bookEntry)
          order.Size -= bookEntry.Size
          self.deleteBookEntry(pricePoint, bookEntry.Order)
          self.deleteOrder(bookEntry.Order)
        } else {
          // we execute the trade we then get rid of the bookEntry
          self.Execution(iterator.Key().(uint64), true, order, bookEntry)
          if bookEntry.Size > order.Size {
            bookEntry.Size -= order.Size
            self.deleteOrder(order)
          } else {
            // this is obviously shit!!!
            self.deleteBookEntry(pricePoint, bookEntry.Order)
            self.deletePricePoint(pricePoint, order.Price)
            self.deleteOrder(bookEntry.Order)
          }
          return
        }
      }
      if ok := iterator.Previous(); ok {
        if len(iterator.Value().(*PricePoint).BuyBookEntries) > 0 {
          self.HighestBid = iterator.Key().(uint64)
        }
      } else {
        finished = true
        self.HighestBid = 0
      }
    }
  }
  self.OpenOrders[order.Id] = order
  bookEntry := self.createBookEntryForOrder(order)
  self.createPricePoint(order.Price, bookEntry)
  if self.LowestAsk > order.Price || self.LowestAsk == 0 {
    self.LowestAsk = order.Price
  }
}

func (self *OrderBook) limitBuy(order *LimitOrder) {
  if order.Price >= self.LowestAsk && self.LowestAsk != 0 {
    iterator := self.PricePoints.Seek(self.LowestAsk)
    finished := false
    for iterator != nil && !finished && order.Price >= self.LowestAsk {
      pricePoint := iterator.Value().(*PricePoint)
      for _, bookEntry := range pricePoint.SellBookEntries {
        if bookEntry.Size < order.Size {
          self.Execution(iterator.Key().(uint64), false, order, bookEntry)
          order.Size -= bookEntry.Size
          self.deleteBookEntry(pricePoint, bookEntry.Order)
          self.deleteOrder(bookEntry.Order)
        } else {
          // we execute the trade we then get rid of the bookEntry
          self.Execution(iterator.Key().(uint64), true, order, bookEntry)
          if bookEntry.Size > order.Size {
            bookEntry.Size -= order.Size
            self.deleteOrder(order)
          } else {
            // this is obviously shit!!!
            self.deleteBookEntry(pricePoint, bookEntry.Order)
            self.deletePricePoint(pricePoint, order.Price)
            self.deleteOrder(bookEntry.Order)
          }
          return
        }
      }
      if ok := iterator.Next(); ok {
        if len(iterator.Value().(*PricePoint).SellBookEntries) > 0 {
          self.LowestAsk = iterator.Key().(uint64)
        }
      } else {
        finished = true
        self.LowestAsk = 0
      }
    }
  }
  self.OpenOrders[order.Id] = order
  bookEntry := self.createBookEntryForOrder(order)
  self.createPricePoint(order.Price, bookEntry)
  if self.HighestBid < order.Price || self.HighestBid == 0 {
    self.HighestBid = order.Price
  }
}

func (self *OrderBook) Limit(order *LimitOrder) {
  if order.Side == BUY {
    self.limitBuy(order)
  } else {
    self.limitSell(order)
  }
}

func (self *OrderBook) MarketOrder(order *LimitOrder) {
  if order.Side == BUY {
    order.Price = self.LowestAsk
    self.Limit(order)
  } else {
    order.Price = self.HighestBid
    self.Limit(order)
  }
}

func (self *OrderBook) Execution(price uint64, fullOrder bool, order *LimitOrder, matchedBookEntry *BookEntry) (error, bool) {
  self.Executer.Execution(price, fullOrder, order, matchedBookEntry)
  return nil, true
}
