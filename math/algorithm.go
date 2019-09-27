package math

import "fmt"

// BubbleSort
func BubbleSort(buf []int) []int {
	times := 0
	for i := 0; i < len(buf)-1; i++ {
		flag := false
		for j := 1; j < len(buf)-i; j++ {
			if buf[j-1] > buf[j] {
				times++
				tmp := buf[j-1]
				buf[j-1] = buf[j]
				buf[j] = tmp
				flag = true
			}
		}
		if !flag {
			break
		}
	}
	return buf
}

// SelectSort
func SelectSort(buf []int) {
	times := 0
	for i := 0; i < len(buf)-1; i++ {
		min := i
		for j := i; j < len(buf); j++ {
			times++
			if buf[min] > buf[j] {
				min = j
			}
		}
		if min != i {
			tmp := buf[i]
			buf[i] = buf[min]
			buf[min] = tmp
		}
	}
	fmt.Println("xuanze times: ", times)
}