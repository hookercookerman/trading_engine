package trading_engine

const (
  BUY  = 0
  SELL = 1
)

type LimitOrder struct {
  Id        string
  Symbol    string
  TraderId  string
  Side      int
  Price     uint64
  Size      uint64
  BookEntry *BookEntry
}
