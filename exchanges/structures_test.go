// structure_test.go
package exchanges

import (
	"testing"
)

var ts = Trades{
	Trade{Tid: 2},
	Trade{Tid: 1},
	Trade{Tid: 3},
	Trade{Tid: 0},
}

func Test_ts_Sort(t *testing.T) {
	ts.Sort()
	for i, td := range ts {
		if td.Tid != int64(i) {
			t.Error("Trades.Sort() Does NOT Work.")
		}
	}
}
func Benchmark_ts_Sort1200(b *testing.B) {
	length := 1200
	ts := make(Trades, length)
	for i := 0; i < length; i++ {
		ts = append(ts, Trade{Tid: int64(i)})
	}
	tsa := ts[:length/2]
	tsb := ts[length/2:]
	ts = append(tsb, tsa...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.Sort()
		b.StopTimer()
		tsa = ts[:length/2]
		tsb = ts[length/2:]
		ts = append(tsb, tsa...)
		b.StartTimer()
	}
}

func Benchmark_ts_SortSorted1200(b *testing.B) {
	length := 1200
	ts := make(Trades, length)
	for i := 0; i < length; i++ {
		ts = append(ts, Trade{Tid: int64(i)})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.Sort()
	}
}

func Benchmark_ts_Sort120000(b *testing.B) {
	length := 120000
	ts := make(Trades, length)
	for i := 0; i < length; i++ {
		ts = append(ts, Trade{Tid: int64(i)})
	}
	tsa := ts[:length/2]
	tsb := ts[length/2:]
	ts = append(tsb, tsa...)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.Sort()
		b.StopTimer()
		tsa = ts[:length/2]
		tsb = ts[length/2:]
		ts = append(tsb, tsa...)
		b.StartTimer()
	}
}

func Benchmark_ts_SortSorted120000(b *testing.B) {
	length := 120000
	ts := make(Trades, length)
	for i := 0; i < length; i++ {
		ts = append(ts, Trade{Tid: int64(i)})
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts.Sort()
	}
}

func Test_ts_CheckUnique(t *testing.T) {
	unique, id := ts.IsUnique()
	if !unique {
		t.Error("ts.CheckUnique() 未能识别出唯一序列。")
	}
	if id != nil {
		t.Error("ts.CheckUnique()为唯一序列时，未能返回0.")
	}
	ts = append(ts, Trade{Tid: 3})
	unique, id = ts.IsUnique()
	if unique {
		t.Error("ts有重复的tid，却被认为是unique。")
	}
	if id[0] != 3 {
		t.Error("ts.CheckUnique()未能找出重复的tid==3")
	}
	ts = ts[0:5]
}
