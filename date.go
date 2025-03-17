package parser

import "regexp"

// checkAbove12 checks if days come before months in dates by looking for numbers > 12
func checkAbove12(numericDates [][]int) *bool {
	for _, date := range numericDates {
		if date[0] > 12 {
			result := true
			return &result
		}
	}

	for _, date := range numericDates {
		if date[1] > 12 {
			result := false
			return &result
		}
	}

	return nil
}

// checkDecreasing checks if days come before months by looking for decreasing numbers
func checkDecreasing(numericDates [][]int) *bool {
	datesByYear := groupArrayByValueAtIndex(numericDates, 2)
	var anyTrue, anyFalse bool

	for _, dates := range datesByYear {
		var daysFirst, daysSecond bool

		for i := 1; i < len(dates); i++ {
			diff1 := dates[i][0] - dates[i-1][0]
			if isNegative(diff1) {
				daysFirst = true
			}

			diff2 := dates[i][1] - dates[i-1][1]
			if isNegative(diff2) {
				daysSecond = true
			}
		}

		if daysFirst {
			anyTrue = true
		}

		if daysSecond {
			anyFalse = true
		}
	}

	if anyTrue {
		result := true
		return &result
	}

	if anyFalse {
		result := false
		return &result
	}

	return nil
}

// changeFrequencyAnalysis analyzes which number changes more frequently
func changeFrequencyAnalysis(numericDates [][]int) *bool {
	if len(numericDates) <= 1 {
		return nil
	}

	diffs := make([][]int, 0)

	for i := 1; i < len(numericDates); i++ {
		diff := []int{
			abs(numericDates[i][0] - numericDates[i-1][0]),
			abs(numericDates[i][1] - numericDates[i-1][1]),
		}
		diffs = append(diffs, diff)
	}

	first := 0
	second := 0

	for _, diff := range diffs {
		first += diff[0]
		second += diff[1]
	}

	if first > second {
		result := true
		return &result
	}
	if first < second {
		result := false
		return &result
	}

	return nil
}

// daysBeforeMonths tries to determine if days come before months in dates
func daysBeforeMonths(numericDates [][]int) *bool {
	firstCheck := checkAbove12(numericDates)
	if firstCheck != nil {
		return firstCheck
	}

	secondCheck := checkDecreasing(numericDates)
	if secondCheck != nil {
		return secondCheck
	}

	return changeFrequencyAnalysis(numericDates)
}

// normalizeDate takes year, month, and day as strings and pads them
func normalizeDate(year, month, day string) [3]string {
	// 2-digit years are assumed to be in the 2000-2099 range
	if len(year) < 4 {
		year = "20" + year
	}

	// Pad month and day with a leading zero if needed
	if len(month) < 2 {
		month = "0" + month
	}

	if len(day) < 2 {
		day = "0" + day
	}

	return [3]string{year, month, day}
}

// orderDateComponents pushes the longest number to the end (assumed to be the year)
func orderDateComponents(date string) [3]string {
	re := regexp.MustCompile(`[-/.] ?`)
	parts := re.Split(date, -1)

	a, b, c := parts[0], parts[1], parts[2]
	maxLength := max(len(a), max(len(b), len(c)))

	if len(c) == maxLength {
		return [3]string{a, b, c}
	}
	if len(b) == maxLength {
		return [3]string{a, c, b}
	}
	return [3]string{b, c, a}
}
