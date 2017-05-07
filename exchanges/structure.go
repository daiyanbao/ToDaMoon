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

//Depth 记录深度信息
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

//Trade 记录一个成交记录的细节
type Trade struct {
	Tid    int64
	Date   int64
	Price  float64
	Amount float64
	Type   string
}

//Attributes 返回Trade记录的细节
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

//Trades 是成交记录Trade的切片
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

//Sort 对Trades进行了排序。
func (ts Trades) Sort() {
	sort.Sort(ts)
}

//IsUnique 检查交易记录切片是否具有重复项
//TODO: 修改这个方法的名称，或者删除
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

//PrintIDDiff 是输出ID的差值
func (ts Trades) PrintIDDiff() {
	for i := 0; i < ts.Len()-1; i++ {
		fmt.Print(ts[i+1].Tid-ts[i].Tid, ",")
	}
}

//IndexOf 返回Trades中date >= 参数date的最小索引值
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

//CopyBetween Result included startDate, But WITHOUT endDate.
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

//DropBefore return Droped Data, Keep the others.
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

//Trans 把Trade的数据类型转变成了*pb.coin型
func (ts Trades) Trans() *pb.Coin {
	result := new(pb.Coin)
	result.Trades = []*pb.Trade{}
	for i := range ts {
		result.Trades = append(result.Trades, ts[i].trans())
	}
	return result
}

//KLine 封装K线
type KLine [][6]float64
