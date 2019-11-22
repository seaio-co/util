package math

import (
	"math/rand"
	"time"
	"fmt"
)

// QuickSort 实现快速排序int[]
func QuickSort(arr []int) []int {

	if len(arr) <= 1 {
		return arr
	}

	median := arr[rand.Intn(len(arr))]

	low_part := make([]int, 0, len(arr))
	high_part := make([]int, 0, len(arr))
	middle_part := make([]int, 0, len(arr))

	for _, item := range arr {
		switch {
		case item < median:
			low_part = append(low_part, item)
		case item == median:
			middle_part = append(middle_part, item)
		case item > median:
			high_part = append(high_part, item)
		}
	}

	low_part = QuickSort(low_part)
	high_part = QuickSort(high_part)

	low_part = append(low_part, middle_part...)
	low_part = append(low_part, high_part...)

	return low_part
}

// RandomArray 生成随机int[]
func RandomArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i <= n-1; i++ {
		arr[i] = rand.Intn(n)
	}
	return arr
}

func bubbleSort(arr []int) []int {
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-1; j++ {
			if arr[j+1] < arr[j] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

func bubbleSortb(arr []int) []int {
	var flag bool
	for i := 0; i < len(arr)-1; i++ {
		flag = false
		for j := 0; j < len(arr)-1; j++ {
			if arr[j+1] < arr[j] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				flag=true
			}
		}
		if !flag {
			break
		}
	}
	return arr
}

// shuffle
func shuffle(arr []int){
	rand.Seed(time.Now().UnixNano())
	var i, j int
	var temp int
	for i = len(arr) - 1; i > 0; i-- {
		j = rand.Intn(i + 1)
		temp = arr[i]
		arr[i] = arr[j]
		arr[j] = temp
	}
}

func randomMoney(remainCount, remainMoney int)int{
	if remainCount == 1{
		return remainMoney
	}

	rand.Seed(time.Now().UnixNano())

	var min = 1
	max := remainMoney / remainCount * 2
	money := rand.Intn(max) + min
	return money
}

// redPackage
func redPackage(count, money int)  {
	for i := 0; i < count; i++ {
		m := randomMoney(count - i, money)
		fmt.Printf("%d  ",  m)
		money -= m
	}
}

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}