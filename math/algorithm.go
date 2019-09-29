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
func FibonacciRecursion(n int)int{
	if n==0 {
		return 0
	}else if n==1 {
		return 1
	}else{
		return FibonacciRecursion(n-1)+FibonacciRecursion(n-2)
	}
}