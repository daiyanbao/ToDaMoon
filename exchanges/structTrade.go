package exchanges

import (
	"ToDaMoon/database"
	"ToDaMoon/util"
	"fmt"
	"sort"
)

//Trade 记录一个成交记录的细节
type Trade struct {
	Tid    int64
	Date   int64
	Price  float64
	Amount float64
	Type   string
}

func (t Trade) String() string {
	str := "*****************\n"
	str += fmt.Sprintf("Tid   :%d\n", t.Tid)
	str += fmt.Sprintf("Date  :%d (%s)\n", t.Date, util.DateOf(t.Date))
	str += fmt.Sprintf("Price :%f\n", t.Price)
	str += fmt.Sprintf("Amount:%f\n", t.Amount)
	str += fmt.Sprintf("Type  :%s\n", t.Type)
	return str
}

//Attributes 实现了database.Attributer接口
func (t *Trade) Attributes() []interface{} {
	return []interface{}{&t.Tid, &t.Date, &t.Price, &t.Amount, &t.Type}
}

//newTrade 返回了一个*Trade变量。
func newTrade() database.Attributer {
	return &Trade{}
}

//Trades 是*Trade的切片
type Trades []*Trade

//Len returns length of ts
func (ts Trades) Len() int {
	return len(ts)
}

//Less 决定了是升序还是降序
func (ts Trades) Less(i, j int) bool {
	return ts[i].Tid < ts[j].Tid
}

//Swap 是交换方式
func (ts Trades) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

//Sort 对Trades进行原地排序
func (ts Trades) Sort() {
	sort.Sort(ts)
}

//After 返回ts中大于tid的部分
func (ts Trades) After(tid int64) Trades {
	i := ts.IndexOf(tid)
	return ts[i+1:]
}

//IndexOf 返回tid所在的Index
func (ts Trades) IndexOf(tid int64) int {
	for i, t := range ts {
		if t.Tid == tid {
			return i
		}
	}

	return ts.Len() - 1
}

/*

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
*/
