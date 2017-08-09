package exchanges

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Test_Ticker_String(t *testing.T) {
	ast := assert.New(t)

	tk := &Ticker{
		Last: 0,
		Buy:  0,
		Sell: 0,
		High: 0,
		Low:  0,
		Vol:  0,
	}

	actual := tk.String()
	expected := `Last:0.000000
Buy :0.000000
Sell:0.000000
High:0.000000
Low :0.000000
Vol :0.000000
`
	ast.Equal(expected, actual, "Ticker的格式化不对")
}
func Test_Depth_String(t *testing.T) {
	ast := assert.New(t)

	d := &Depth{
		Asks: Quotations{Quotation{1, 2}},
		Bids: Quotations{Quotation{1, 2}},
	}

	actual := d.String()
	expected := `Asks
	Price		Amount
	1.000000	2.000000
Bids
	Price		Amount
	1.000000	2.000000
`
	ast.Equal(expected, actual, "Depth的格式化不对")
}

func Test_IsAskSorted(t *testing.T) {
	ast := assert.New(t)

	data := Quotations{
		Quotation{1, 0},
		Quotation{0, 0},
	}

	ast.False(data.IsAskSorted(), "data不是升序排列的，却没有返回false")
}
func Test_IsBidSorted(t *testing.T) {
	ast := assert.New(t)

	data := Quotations{
		Quotation{0, 0},
		Quotation{1, 0},
	}

	ast.False(data.IsBidSorted(), "data不是降序排列的，却没有返回false")
}

func Test_Order_String(t *testing.T) {
	ast := assert.New(t)

	o := &Order{
		ID:     0,
		Date:   0,
		Money:  "cny",
		Price:  0,
		Coin:   "btc",
		Amount: 0,
		Type:   "buy",
	}

	actual := o.String()
	expected := `ID    :0
Date  :0
Money :cny
Price :0.000000
Coin  :btc
Amount:0.000000
Type  :buy
`
	ast.Equal(expected, actual, "Order的格式化不对")
}
