package math

import (
	"fmt"
	"testing"
)

func Test_QuickSort(t *testing.T) {

	arr := RandomArray(10)

	fmt.Println("Initial array is:", arr)
	fmt.Println("")
	fmt.Println("Sorted array is: ", QuickSort(arr))
}
