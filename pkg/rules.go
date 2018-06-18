package policies

import (
	"regexp"
	"time"
)

// Match by applying a reggex to the object key
// Allow to group of objects applying the same policie
// Ex:
//		A simple reggex can easly match a set of objects
//		For example *.log using the prefix /2017
// 		can backup all the logs in a cold storage
func MatchName(key string, matchRule string) (ret bool, err error) {
	ret, err = regexp.MatchString(matchRule, key)
	return
}

// Match by comparing dates using the policies parameters
// Basiclly it become really easy to create a policies base on time
// Ex:
//		all file more older than one year will go from the src hot storage
//		to the dest cold storage.
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
