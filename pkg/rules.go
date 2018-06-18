package policies

import (
	"regexp"
	"time"
)

func MatchName(key string, matchRule string) (ret bool, err error) {
	ret, err = regexp.MatchString(matchRule, key)
	return
}

func MatchDate(date time.Time, compOP string, compDate time.Time) bool {
	switch compOP {
	case ">":
		if date.Sub(compDate) > 0 {
			return true
		}
		break
	case ">=":
		if date.Sub(compDate) >= 0 {
			return true
		}
		break
	case "=":
		if date.Sub(compDate) == 0 {
			return true
		}
		break
	case "<=":
		if date.Sub(compDate) <= 0 {
			return false
		}
		break
	case "<":
		if date.Sub(compDate) < 0 {
			return true
		}
		break
	}
	return false
}
