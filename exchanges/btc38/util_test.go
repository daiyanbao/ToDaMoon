package btc38

import "testing"

func Test_priceStr(t *testing.T) {
	data := map[float64]string{
		123456789.12345: "123456789",
		12345.6789:      "12345.6",
		123.456789:      "123.456",
		1.23456789:      "1.23457",
		0.00123456789:   "0.00123",
	}

	for k, v := range data {
		psk := priceStr(k)
		if psk != v {
			t.Errorf("%f应该被转换成%s，而不是%s", k, v, psk)
		}
	}
}
