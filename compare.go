package aoj

import (
	"bufio"
	"bytes"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func fuzzyCompare(a, b []byte) bool {
	if bytes.Compare(a, b) == 0 {
		return true
	}
	logdbg("two bytes are not the same.")

	sa := bufio.NewScanner(bytes.NewReader(a))
	sb := bufio.NewScanner(bytes.NewReader(b))

	for {
		haveMoreA := sa.Scan()
		haveMoreB := sb.Scan()
		if haveMoreA != haveMoreB {
			logdbg("fussyCompare reached  end only one output")
			return false // reached end only one side
		} else if !haveMoreA {
			logdbg("fussyCompare reached end of both output")
			return true // reached end of both data
		}
		// have more tokens
		ta := strings.Split(sa.Text(), " ")
		tb := strings.Split(sb.Text(), " ")

		if len(ta) != len(tb) {
			return false
		}

		for i := 0; i < len(ta); i++ {
			logdbg("fuzzyCompare checking '%s', '%s'", ta[i], tb[i])
			if ta[i] == tb[i] {
				continue
			}
			if eqFloatString(ta[i], tb[i]) {
				continue
			}
			return false
		}
	}
	return true
}

func isFloat(s string) bool {
	for _, r := range s {
		if r == '.' || unicode.IsDigit(r) {
			continue
		}
		return false
	}
	return true
}

func eqFloatString(s, t string) bool {
	if !isFloat(s) || !isFloat(t) {
		return false
	}

	fs, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	ft, err := strconv.ParseFloat(t, 64)
	if err != nil {
		return false
	}
	return eqFloat(fs, ft)
}

func eqFloat(a, b float64) bool {
	return math.Abs(a-b) <= Epsilon
}
