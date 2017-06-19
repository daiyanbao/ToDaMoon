package btc38

import (
	"fmt"
	"strings"
)

func price2Str(price float64) string {
	ps := fmt.Sprintf("%.5f", price)

	i := strings.Index(ps, ".")
	if i == 1 {
		return ps
	}

	remain := 7 - i
	if remain > 0 {
		return ps[:7]
	}

	return ps[:i]
}
