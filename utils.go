package parser

import (
	"sort"
	"strconv"
)

// indexAboveValue returns a function that checks if the value at the given
// index in an array is greater than the given value
func indexAboveValue(index int, value int) func([]int) bool {
	return func(array []int) bool {
		return array[index] > value
	}
}

// isNegative returns true for a negative number, false otherwise
func isNegative(number int) bool {
	return number < 0
}

// groupArrayByValueAtIndex takes an array of arrays and an index and groups
// the inner arrays by the value at the index provided
func groupArrayByValueAtIndex(array [][]int, index int) [][][]int {
	groups := make(map[string][][]int)

	for _, item := range array {
		key := "_" + strconv.Itoa(item[index])
		groups[key] = append(groups[key], item)
	}

	result := make([][][]int, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// abs returns the absolute value of n
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// max returns the larger of x or y
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func areStringArraysIdenticalIgnoringOrder(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	sortedArr1 := make([]string, len(arr1))
	sortedArr2 := make([]string, len(arr2))

	copy(sortedArr1, arr1)
	copy(sortedArr2, arr2)

	sort.Strings(sortedArr1)
	sort.Strings(sortedArr2)

	for i := range sortedArr1 {
		if sortedArr1[i] != sortedArr2[i] {
			return false
		}
	}

	return true
}
