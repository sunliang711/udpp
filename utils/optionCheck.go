package utils

import "fmt"

func OptionCheckd(o ...int) []string {
	var ret []string
	for i, v := range o {
		if v == 1 {
			ret = append(ret, fmt.Sprintf("%d", i))
		}
	}
	return ret
}
