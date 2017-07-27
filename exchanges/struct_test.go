package exchanges

import "testing"

var q = Quotations{
	Quotation{
		Price:  3,
		Amount: 3,
	},
	Quotation{
		Price:  1,
		Amount: 1,
	},
	Quotation{
		Price:  2,
		Amount: 2,
	},
}

func Test_Depth_SortAsks(t *testing.T) {
	q.SortAsks()
	if !q.IsAskSorted() {
		t.Error("无法进行Ask序排序")
		t.Log(q)
	}
}

func Test_Depth_SortBids(t *testing.T) {
	q.SortBids()
	if !q.IsBidSorted() {
		t.Error("无法进行Bid序排序")
		t.Log(q)
	}
}
