package collection

import "saas/app/helper/constraints"

func Unique[T constraints.Ordered](array []T) []T {
	temp := make([]T, 0)
	for i := 0; i < len(array); i++ {
		repeat := false
		for j := i + 1; j < len(array); j++ {
			if array[i] == array[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			temp = append(temp, array[i])
		}
	}
	return temp
}
