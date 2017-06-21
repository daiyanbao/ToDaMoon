package exchanges

import (
	"testing"
)

var ts = Trades{
	&Trade{Tid: 2},
	&Trade{Tid: 1},
	&Trade{Tid: 3},
	&Trade{Tid: 0},
}

func Test_ts_Sort(t *testing.T) {
	ts.Sort()
	for i, td := range ts {
		//检查td经过排序后，是否是升序

		if td.Tid != int64(i) {
			t.Error("Trades.Sort() Does NOT Work.")
		}
	}
	t.Log("ts经过原地排序后，确实是升序排列")
}
