package math

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