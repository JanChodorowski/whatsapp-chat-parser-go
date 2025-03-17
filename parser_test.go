package parser

import (
	"os"
	"testing"
	"time"
)

// TestParseString tests the main ParseString function
func TestParseString(t *testing.T) {

	t.Run("Empty string", func(t *testing.T) {
		messages, err := ParseString("", nil)
		if err != nil {
			t.Errorf("Expected no error for empty string, got %v", err)
		}
		if len(messages) != 0 {
			t.Errorf("Expected empty array for empty string, got %d messages", len(messages))
		}
	})

	for _, chatExample := range chatExamples {
		fileContents, err := os.ReadFile(chatExample.filePath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		t.Run("Parse messages count "+chatExample.description, func(t *testing.T) {
			messages, err := ParseString(string(fileContents), nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(messages) != chatExample.messagesCount {
				t.Errorf("Expected %d messages, got %d", chatExample.messagesCount, len(messages))
			}
		})

		t.Run("Get authors from messages "+chatExample.description, func(t *testing.T) {
			messages, err := ParseString(string(fileContents), nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			authors := GetAuthorsFromMessages(&messages)

			test := areStringArraysIdenticalIgnoringOrder(authors, chatExample.authors)

			if !test {
				t.Errorf("Expected authors to be %v, got %v", chatExample.authors, authors)
			}

		})

		t.Run("Get first and last message dates "+chatExample.description, func(t *testing.T) {
			messages, err := ParseString(string(fileContents), nil)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			firstDate, lastDate := GetFirstAndLastMessageDates(&messages)

			if firstDate == nil || lastDate == nil {
				t.Errorf("Expected dates to be non-nil")
				return
			}

			if firstDate.Day() != chatExample.firstDate.Day() ||
				firstDate.Month() != chatExample.firstDate.Month() ||
				firstDate.Year() != chatExample.firstDate.Year() ||
				firstDate.Hour() != chatExample.firstDate.Hour() ||
				firstDate.Minute() != chatExample.firstDate.Minute() {
				t.Errorf("Expected first date to be %v, got %v", chatExample.firstDate, firstDate)
			}

			if lastDate.Day() != chatExample.lastDate.Day() ||
				lastDate.Month() != chatExample.lastDate.Month() ||
				lastDate.Year() != chatExample.lastDate.Year() ||
				lastDate.Hour() != chatExample.lastDate.Hour() ||
				lastDate.Minute() != chatExample.lastDate.Minute() {
				t.Errorf("Expected first date to be %v, got %v", chatExample.lastDate, lastDate)
			}

		})
	}
}

// TestTime tests time-related functions
func TestTime(t *testing.T) {
	t.Run("convertTime12to24", func(t *testing.T) {
		tests := []struct {
			time   string
			ampm   string
			expect string
		}{
			{"12:00", "PM", "12:00"},
			{"12:00", "AM", "00:00"},
			{"05:06", "PM", "17:06"},
			{"07:19", "AM", "07:19"},
			{"01:02:34", "PM", "13:02:34"},
			{"02:04:54", "AM", "02:04:54"},
		}

		for _, test := range tests {
			result := convertTime12to24(test.time, test.ampm)
			if result != test.expect {
				t.Errorf("convertTime12to24(%q, %q) = %q, want %q",
					test.time, test.ampm, result, test.expect)
			}
		}
	})

	t.Run("normalizeAMPM", func(t *testing.T) {
		tests := []struct {
			input  string
			expect string
		}{
			{"am", "AM"},
			{"pm", "PM"},
			{"a.m.", "AM"},
			{"p.m.", "PM"},
			{"A.M.", "AM"},
			{"P.M.", "PM"},
		}

		for _, test := range tests {
			result := normalizeAMPM(test.input)
			if result != test.expect {
				t.Errorf("normalizeAMPM(%q) = %q, want %q",
					test.input, result, test.expect)
			}
		}
	})

	t.Run("normalizeTime", func(t *testing.T) {
		tests := []struct {
			input  string
			expect string
		}{
			{"12:34", "12:34:00"},
			{"1:23:45", "01:23:45"},
			{"12:34:56", "12:34:56"},
		}

		for _, test := range tests {
			result := normalizeTime(test.input)
			if result != test.expect {
				t.Errorf("normalizeTime(%q) = %q, want %q",
					test.input, result, test.expect)
			}
		}
	})
}

// TestUtils tests utility functions
func TestUtils(t *testing.T) {
	t.Run("indexAboveValue", func(t *testing.T) {
		array := []int{34, 16}

		if !indexAboveValue(0, 33)(array) {
			t.Errorf("Expected indexAboveValue(0, 33)([34, 16]) to be true")
		}
		if indexAboveValue(0, 34)(array) {
			t.Errorf("Expected indexAboveValue(0, 34)([34, 16]) to be false")
		}
		if !indexAboveValue(1, 15)(array) {
			t.Errorf("Expected indexAboveValue(1, 15)([34, 16]) to be true")
		}
		if indexAboveValue(1, 16)(array) {
			t.Errorf("Expected indexAboveValue(1, 16)([34, 16]) to be false")
		}
	})

	t.Run("isNegative", func(t *testing.T) {
		if !isNegative(-1) {
			t.Errorf("Expected isNegative(-1) to be true")
		}
		if !isNegative(-15) {
			t.Errorf("Expected isNegative(-15) to be true")
		}
		if isNegative(0) {
			t.Errorf("Expected isNegative(0) to be false")
		}
		if isNegative(1) {
			t.Errorf("Expected isNegative(1) to be false")
		}
	})

	t.Run("groupArrayByValueAtIndex", func(t *testing.T) {
		array := [][]int{
			{8, 30, 3},
			{9, 50, 3},
			{6, 30, 3},
		}

		result0 := groupArrayByValueAtIndex(array, 0)
		if len(result0) != 3 {
			t.Errorf("Expected 3 groups, got %d", len(result0))
		}

		result1 := groupArrayByValueAtIndex(array, 1)
		if len(result1) != 2 {
			t.Errorf("Expected 2 groups, got %d", len(result1))
		}

		result2 := groupArrayByValueAtIndex(array, 2)
		if len(result2) != 1 {
			t.Errorf("Expected 1 group, got %d", len(result2))
		}
	})
}

// TestDate tests date-related functions
func TestDate(t *testing.T) {
	t.Run("checkAbove12", func(t *testing.T) {
		daysFirst := [][]int{
			{13, 6, 2017},
			{26, 11, 2017},
		}
		monthsFirst := [][]int{
			{4, 13, 2017},
			{6, 15, 2017},
		}
		undetectable := [][]int{
			{4, 6, 2017},
			{11, 10, 2017},
		}

		if result := checkAbove12(daysFirst); result == nil || !*result {
			t.Errorf("Expected checkAbove12(daysFirst) to be true")
		}
		if result := checkAbove12(monthsFirst); result == nil || *result {
			t.Errorf("Expected checkAbove12(monthsFirst) to be false")
		}
		if result := checkAbove12(undetectable); result != nil {
			t.Errorf("Expected checkAbove12(undetectable) to be nil")
		}
	})

	t.Run("normalizeDate", func(t *testing.T) {
		result := normalizeDate("11", "3", "4")
		expected := [3]string{"2011", "03", "04"}

		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		result2 := normalizeDate("2011", "03", "04")
		if result2 != expected {
			t.Errorf("Already normalized date should not change: expected %v, got %v", expected, result2)
		}
	})
}

// TestParseMessages tests parsing messages with different formats
func TestParseMessages(t *testing.T) {
	t.Run("Parse different date formats", func(t *testing.T) {
		formats := []string{
			"3/6/18, 1:55 p.m. - a: m",
			"03-06-2018, 01.55 PM - a: m",
			"13.06.18 21.25.15: a: m",
			"[06.13.18 21:25:15] a: m",
			"13.6.2018 klo 21.25.15 - a: m",
			"13. 6. 2018. 21:25:15 a: m",
			"[3/6/18 1:55:00 p. m.] a: m",
			"[2018/06/13, 21:25:15] a: m",
		}

		for _, format := range formats {
			messages, err := ParseString(format, nil)
			if err != nil {
				t.Errorf("Failed to parse format %q: %v", format, err)
				continue
			}

			if len(messages) != 1 {
				t.Errorf("Expected 1 message for format %q, got %d", format, len(messages))
				continue
			}

			message := messages[0]
			if message.Author == nil {
				t.Errorf("Expected author for format %q, got nil", format)
				continue
			}

			if *message.Author != "a" {
				t.Errorf("Expected author 'a' for format %q, got %q", format, *message.Author)
			}

			if message.Message != "m" {
				t.Errorf("Expected message 'm' for format %q, got %q", format, message.Message)
			}
		}
	})

	t.Run("Parse attachments", func(t *testing.T) {
		formats := []string{
			"3/6/18, 1:55 p.m. - a: < attached: photo.jpg >",
			"3/6/18, 1:55 p.m. - a: IMG-20210428-WA0001.jpg (file attached)",
		}

		options := ParseStringOptions{
			ParseAttachments: true,
		}

		for _, format := range formats {
			messages, err := ParseString(format, &options)
			if err != nil {
				t.Errorf("Failed to parse format %q: %v", format, err)
				continue
			}

			if len(messages) != 1 {
				t.Errorf("Expected 1 message for format %q, got %d", format, len(messages))
				continue
			}

			message := messages[0]
			if message.Attachment == nil {
				t.Errorf("Expected attachment for format %q, got nil", format)
				continue
			}
		}
	})

	t.Run("Multiline messages", func(t *testing.T) {

		multilineMessage := `09/04/2017, 01:50 - +410123456789: How are you?

Is everything alright?`

		messages, err := ParseString(multilineMessage, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		lastMessage := messages[0]

		expected := "How are you?\n\nIs everything alright?"
		if lastMessage.Message != expected {
			t.Errorf("Expected multiline message to be %q, got %q", expected, lastMessage.Message)
		}
	})

}

type chatTestExample struct {
	description   string
	filePath      string
	authors       []string
	firstDate     time.Time
	lastDate      time.Time
	messagesCount int
}

var chatExamples = []chatTestExample{
	{
		description:   "default test",
		filePath:      "test_data/default.txt",
		authors:       []string{"Sample User", "TestBot", "+410123456789"},
		firstDate:     time.Date(2017, 6, 3, 0, 45, 0, 0, time.UTC),
		lastDate:      time.Date(2017, 9, 4, 1, 50, 0, 0, time.UTC),
		messagesCount: 5,
	},
	{
		description:   "english iPhone saved contacts",
		filePath:      "test_data/english_iphone-saved_contacts.txt",
		authors:       []string{"Andrew", "Josh", "Marthy"},
		firstDate:     time.Date(2018, 11, 29, 10, 50, 0, 0, time.UTC),
		lastDate:      time.Date(2018, 11, 29, 11, 36, 0, 0, time.UTC),
		messagesCount: 15,
	},
	{
		// the first system message is newer than the rest
		// becase the chat was cleared
		description:   "English Android, no saved contacts",
		filePath:      "test_data/english_android-unsaved_contacts.txt",
		authors:       []string{"+33 6 99 88 77 66", "+48 777 666 555", "Andrew"},
		firstDate:     time.Date(2025, 3, 18, 17, 49, 0, 0, time.UTC),
		lastDate:      time.Date(2025, 3, 10, 16, 44, 0, 0, time.UTC),
		messagesCount: 20,
	},
	// {
	// 	description:   "English iPhone, no saved contacts",
	// 	filePath:      "test_data/english_iphone-no_saved_contacts.txt",
	// 	authors:        []string{"Sample User", "TestBot", "+410123456789"},
	// 	firstDate:     time.Date(2017, 3, 6, 0, 45, 0, 0, time.UTC),
	// 	lastDate:      time.Date(2017, 4, 9, 1, 50, 0, 0, time.UTC),
	// 	messagesCount: 1,
	// },
	// {
	// 	description:   "",
	// 	filePath:      "test_data/english_iphone-saved_contacts.txt",
	// 	authors:        []string{"Sample User", "TestBot", "+410123456789"},
	// 	firstDate:     time.Date(2017, 3, 6, 0, 45, 0, 0, time.UTC),
	// 	lastDate:      time.Date(2017, 4, 9, 1, 50, 0, 0, time.UTC),
	// 	messagesCount: 1,
	// },
}
