//Package exchanges 的structure.go定义了交易所需用到的数据格式。
package exchanges

import (
	pb "ToDaMoon/exchanges/trades"
	"fmt"
	"sort"
)

//Ticker 是ticker的数据结构。
type Ticker struct {
	Last float64
	Buy  float64
	Sell float64
	High float64
	Low  float64
	Vol  float64
}

func (t *Ticker) String() string {
	str := fmt.Sprintf("Last:%f\n", t.Last)
	str += fmt.Sprintf("Buy :%f\n", t.Buy)
	str += fmt.Sprintf("Sell:%f\n", t.Sell)
	str += fmt.Sprintf("High:%f\n", t.High)
	str += fmt.Sprintf("Low :%f\n", t.Low)
	str += fmt.Sprintf("Vol :%f\n", t.Vol)
	return str
}

type Depth struct {
	Asks [][2]float64
	Bids [][2]float64
}

func (d *Depth) String() string {
	str := fmt.Sprintln("Asks:")
	for _, v := range d.Asks {
		str += "\t" + fmt.Sprintln(v)
	}
	str += fmt.Sprintln("Bids:")
	for _, v := range d.Bids {
		str += "\t" + fmt.Sprintln(v)
	}
	return str
}

type Trade struct {
	Tid    int64
	Date   int64
	Price  float64
	Amount float64
	Type   string
}

func (t Trade) Attributes() (int64, int64, float64, float64, string) {
	return t.Tid, t.Date, t.Price, t.Amount, t.Type
}

func (t Trade) String() string {
	str := "*****************\n"
	str += fmt.Sprintf("Tid   :%d\n", t.Tid)
	str += fmt.Sprintf("Date  :%d\n", t.Date)
	str += fmt.Sprintf("Price :%f\n", t.Price)
	str += fmt.Sprintf("Amount:%f\n", t.Amount)
	str += fmt.Sprintf("Type  :%s\n", t.Type)
	return str
}
func (t Trade) trans() *pb.Trade {
	result := pb.Trade(t)
	return &result
}

// TidSlice attaches the methods of sort.Interface to []int64, sorting in increasing order.
type TidSlice []int64

func (s TidSlice) Len() int           { return len(s) }
func (s TidSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s TidSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Sort is a convenience method.
func (s TidSlice) Sort() {
	sort.Sort(s)
}

// // SearchInt64s searches for x in a sorted slice of int64 and returns the index
// // as specified by sort.Search. The slice must be sorted in ascending order.
// func SearchInt64s(a []int64, x int64) int {
// 	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
// }

type Trades []Trade

func (ts Trades) Len() int {
	return len(ts)
}

func (ts Trades) Less(i, j int) bool {
	return ts[i].Tid < ts[j].Tid
}

func (ts Trades) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts Trades) Sort() {
	sort.Sort(ts)
}

func (ts Trades) IsUnique() (bool, []int64) {
	var repeatID []int64
	if ts.Len() < 2 {
		return true, nil
	}
	tids := make(TidSlice, ts.Len())
	for i := range ts {
		tids[i] = ts[i].Tid
	}
	tids.Sort()

	tempTid := tids[0]

	for i := 1; i < tids.Len(); i++ {
		if tempTid == tids[i] {
			repeatID = append(repeatID, tempTid)
		} else {
			tempTid = tids[i]
		}
	}

	if repeatID == nil {
		return true, nil
	}
	return false, repeatID
}

func (ts Trades) PrintIDDiff() {
	for i := 0; i < ts.Len()-1; i++ {
		fmt.Print(ts[i+1].Tid-ts[i].Tid, ",")
	}
}

func (ts Trades) IndexOf(date int64) int {
	length := len(ts)
	switch {
	case date <= ts[0].Date:
		return 0
	case ts[length-1].Date < date:
		return length
	default:
		for i, t := range ts {
			if date <= t.Date {
				return i
			}
		}
	}
	panic("NEVER REACH THIS.")
}

//Result included startDate, But WITHOUT endDate.
func (ts Trades) CopyBetween(startDate, endDate int64) Trades {
	ts.Sort()
	s := ts.IndexOf(startDate)
	e := ts.IndexOf(endDate)
	if s == e {
		return nil
	}
	result := make(Trades, e-s, e-s)
	copy(result, ts[s:e:e])
	return result
}

// DropBefore return Droped Data, Keep the others.
func (ts *Trades) DropBefore(startDate int64) Trades {
	s := ts.IndexOf(startDate)
	if s == 0 {
		return nil
	}
	result := make(Trades, s, s)
	copy(result, (*ts)[:s:s]) // copy丢弃的数据到新的Trades可以避免，修改丢弃数据时候，产生对ts的修改。
	*ts = (*ts)[s:ts.Len():ts.Len()]
	return result
}

// CutFirstDate 会cut掉ts中,所有包含ts[0].Date值的数据，并作为返回值。
// 这个默认 ts已经是排序好了的。
func (ts *Trades) CutFirstDate() Trades {
	result := Trades{}
	firstDate := (*ts)[0].Date

	// 首先处理ts中只有一个Date值的情况。
	if firstDate == (*ts)[ts.Len()-1].Date {
		result = append(result, (*ts)...)
		*ts = (*ts)[ts.Len():ts.Len()]
		return result
	}

	for i, t := range *ts {
		if t.Date != firstDate {
			result = append(result, (*ts)[:i]...)
			*ts = (*ts)[i:ts.Len()]

			// fmt.Println(i)
			//fmt.Println("result", result[result.Len()-1].Tid)

			return result
		}
	}
	panic("NEVER BE HERE")
}

func (ts Trades) Trans() *pb.Coin {
	result := new(pb.Coin)
	result.Trades = []*pb.Trade{}
	for i := range ts {
		result.Trades = append(result.Trades, ts[i].trans())
	}
	return result
}

type KLine [][6]float64
