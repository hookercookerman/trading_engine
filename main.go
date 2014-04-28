package main

import . "github.com/hookercookerman/trading_engine/trading_engine"

import "net/http"
import "encoding/json"
import "io/ioutil"
import "fmt"
import "strconv"

import "github.com/nu7hatch/gouuid"

type Executer struct {
  sum uint64
}

func (self *Executer) Execution(price uint64, fullOrder bool, order *LimitOrder, matchedBookEntry *BookEntry) (error, bool) {
  var orderSize uint64
  if fullOrder {
    orderSize = order.Size
  } else {
    orderSize = matchedBookEntry.Size
  }
  fmt.Printf("Original Order: Side: %b, Price: $%.2f,  Size: BTC %.9f  \n\n", order.Side, float64(order.Price)/100, float64(order.Size)/100000000)
  fmt.Printf("Trade Price: $%6.2f \n", float64(price)/100)
  fmt.Printf("Trade Amount: BTC %.9f \n", float64(orderSize)/100000000)
  self.sum += orderSize
  fmt.Printf("Trade Final Amount Size %.9f \n\n", float64(self.sum)/100000000)
  return nil, true
}

func main() {
  url := "https://www.bitstamp.net/api/order_book/"
  resp, err := http.Get(url)

  if err != nil {
    panic(err.Error())
  }

  body, _ := ioutil.ReadAll(resp.Body)
  var objmap map[string]*json.RawMessage
  json.Unmarshal(body, &objmap)

  var bids [][]string
  var asks [][]string
  json.Unmarshal(*objmap["bids"], &bids)
  json.Unmarshal(*objmap["asks"], &asks)

  bidOrders := []*LimitOrder{}
  sellOrders := []*LimitOrder{}
  orderBook := NewOrderBook()
  executer := &Executer{}
  orderBook.Executer = executer

  for _, bid := range bids {
    id, _ := uuid.NewV4()
    price, _ := strconv.ParseFloat(bid[0], 64)
    size, _ := strconv.ParseFloat(bid[1], 64)
    price_int := uint64(price * 100)
    size_int := uint64(size * 100000000)
    order := &LimitOrder{Id: id.String(), Side: BUY, Size: size_int, Price: price_int}
    bidOrders = append(bidOrders, order)
    orderBook.Limit(order)
  }
  for _, ask := range asks {
    id, _ := uuid.NewV4()
    price, _ := strconv.ParseFloat(ask[0], 64)
    size, _ := strconv.ParseFloat(ask[1], 64)
    price_int := uint64(price * 100)
    size_int := uint64(size * 100000000)
    order := &LimitOrder{Id: id.String(), Side: SELL, Size: size_int, Price: price_int}
    sellOrders = append(sellOrders, order)
    orderBook.Limit(order)
  }

  id, _ := uuid.NewV4()
  mega := &LimitOrder{Id: id.String(), Side: BUY, Size: 100000000, Price: 45000}
  orderBook.Limit(mega)
  fmt.Println(orderBook.LowestAsk)
}
