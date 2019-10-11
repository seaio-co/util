package math

import "math"

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

// BinarySearch
func BinarySearch(nums []int,left,right,val int) int {
	k := (left+right)/2
	if nums[k]>val {
		return BinarySearch(nums,left,k,val)
	}else if nums[k] < val {
		return BinarySearch(nums,k,right,val)
	}else{
		return k
	}
}

// FibonacciRecursion
func FibonacciRecursion(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return FibonacciRecursion(n-1) + FibonacciRecursion(n-2)
	}
}

// SelectSort
func SelectSort(buf []int) []int {
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
	return buf
}

// InsertSort
func InsertSort(buf []int) []int {
	times := 0
	for i := 1; i < len(buf); i++ {
		for j := i; j > 0; j-- {
			if buf[j] < buf[j-1] {
				times++
				tmp := buf[j-1]
				buf[j-1] = buf[j]
				buf[j] = tmp
			} else {
				break
			}
		}
	}
	return buf
}

// ShellSort
func ShellSort(arr []int) []int {
	increment:=len(arr)
	for {
		increment = increment / 2
		for i := 0; i < increment; i++ {
			for j := i + increment; j < len(arr); j = j + increment {
				for k := j; k > i; k = k - increment {
					if arr[k] < arr[k-increment] {
						arr[k], arr[k-increment] = arr[k-increment], arr[k]
					} else {
						break
					}
				}
			}
		}
		if increment == 1 {
			break
		}
	}
	return arr
}

// FibonacciFind
func FibonacciFind(n int)int{
	x,y,fib := 0,1,0
	for i:=0;i<=n;i++{
		if i==0 {
			fib=0
		}else if i== 1{
			fib = x+y
		}else{
			fib=x+y
			x,y = y,fib
		}
	}
	return fib
}

// Reverse
func Reverse(s []int) {
	for i, j := 0, len(s) -1 ; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func binarySearch(arr []int,  k int) int {
	left, right, mid := 1, len(arr), 0
	for {
		mid = int(math.Floor(float64((left + right) / 2)))
		if arr[mid] > k {
			right = mid - 1
		} else if arr[mid] < k {
			left = mid + 1
		} else {
			break
		}
		if left > right {
			mid = -1
			break
		}
	}
	return mid
}

func binarySearch2(sortedArray []int, lookingFor int) int {
	var low int = 0
	var high int = len(sortedArray) - 1
	for low <= high {
		var mid int =low + (high - low)/2
		var midValue int = sortedArray[mid]
		if midValue == lookingFor {
			return mid
		} else if midValue > lookingFor {
			high = mid -1
		} else {
			low = mid + 1
		}
	}
	return -1
}


func binarySearch3(arr []int, k int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := low + (high-low)>>1
		if k < arr[mid] {
			high = mid - 1
		} else if k > arr[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func binarySearch4(arr []int, k int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := low & high  + (low ^ high) >> 1
		if k < arr[mid] {
			high = mid - 1
		} else if k > arr[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}