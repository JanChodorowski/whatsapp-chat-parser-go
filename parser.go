package parser

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	sharedRegex           = `^(?:[\x{200E}\x{200F}])*\[?(\d{1,4}[-/.]\s?\d{1,4}[-/.]\s?\d{1,4})[,.]?\s\D*?(\d{1,2}[.:]\d{1,2}(?:[.:]\d{1,2})?)(?:\s([ap]\.?\s?m\.?))?\]?(?:\s-|:)?\s`
	authorAndMessageRegex = `(.+?):\s((?s:.*))`
	messageRegex          = `((?s:.*))`
	regexParser           = regexp.MustCompile(`(?i)` + sharedRegex + authorAndMessageRegex)
	regexParserSystem     = regexp.MustCompile(`(?i)` + sharedRegex + messageRegex)
	regexAttachment       = regexp.MustCompile(`(?:\x{200E}|\x{200F})*(?:<.+:(.+)>|([\w-]+\.\w+)\s[(<].+[)>])`)
	regexSplitTime        = regexp.MustCompile(`[:.]+`)
	newlinesRegex         = regexp.MustCompile(`(?:\r\n|\r|\n)`)
)

func isNotNewFormatSystemMessage(message string) bool {
	return strings.Count(message, "\u200E") != 1
}

// makeArrayOfMessages takes an array of lines and detects multiline messages
func makeArrayOfMessages(lines []string) []RawMessage {
	var result []RawMessage

	for _, line := range lines {
		if !regexParser.MatchString(line) && !regexParserSystem.MatchString(line) {
			// If the line doesn't match either regex pattern, it's part of a previous message
			if len(result) > 0 {
				lastIndex := len(result) - 1
				prevMsg := result[lastIndex]
				result[lastIndex] = RawMessage{
					System: prevMsg.System,
					Msg:    prevMsg.Msg + "\n" + line,
				}
			}
		} else {
			// This is a new message
			if regexParser.MatchString(line) && isNotNewFormatSystemMessage(line) {
				result = append(result, RawMessage{
					System: false,
					Msg:    line,
				})
			} else {
				// It's a system message
				result = append(result, RawMessage{
					System: true,
					Msg:    line,
				})
			}
		}
	}

	return result
}

// parseMessageAttachment parses a message to extract attachment details
func parseMessageAttachment(message string) *Attachment {
	matches := regexAttachment.FindStringSubmatch(message)
	if matches == nil {
		return nil
	}

	var fileName string
	if len(matches) > 1 && matches[1] != "" {
		fileName = strings.TrimSpace(matches[1])
	} else if len(matches) > 2 {
		fileName = strings.TrimSpace(matches[2])
	}

	return &Attachment{
		FileName: fileName,
	}
}

// parseMessages parses an array of raw messages into structured messages
func parseMessages(messages []RawMessage, options ParseStringOptions) ([]Message, error) {
	var result []Message
	var allDates [][]int

	// First pass: collect date components for format detection and create message objects
	for _, rawMsg := range messages {
		var matches []string

		if rawMsg.System {
			matches = regexParserSystem.FindStringSubmatch(rawMsg.Msg)
		} else {
			matches = regexParser.FindStringSubmatch(rawMsg.Msg)
		}

		if matches == nil {
			continue
		}

		// Extract date components for format detection
		dateStr := matches[1]
		dateParts := orderDateComponents(dateStr)

		dateComponents := make([]int, 3)
		for i, part := range dateParts {
			val, _ := strconv.Atoi(part)
			dateComponents[i] = val
		}

		allDates = append(allDates, dateComponents)

		// Create message objects with just author and text for now
		var message Message
		if rawMsg.System {
			messageText := matches[4]
			message = Message{
				Author:  nil,
				Message: messageText,
			}
		} else {
			author := matches[4]
			messageText := matches[5]
			message = Message{
				Author:  &author,
				Message: messageText,
			}
		}

		result = append(result, message)
	}

	// Determine if days come first
	var daysFirst bool
	if options.DaysFirst != nil {
		daysFirst = *options.DaysFirst
	} else {
		daysFirstPtr := daysBeforeMonths(allDates)
		if daysFirstPtr != nil {
			daysFirst = *daysFirstPtr
		} else {
			daysFirst = true // Default assumption
		}
	}

	// Second pass: add proper date objects and preserve full message content
	for i, rawMsg := range messages {
		if i >= len(result) {
			break
		}

		var matches []string
		if rawMsg.System {
			matches = regexParserSystem.FindStringSubmatch(rawMsg.Msg)
		} else {
			matches = regexParser.FindStringSubmatch(rawMsg.Msg)
		}

		if matches == nil {
			continue
		}

		dateStr := matches[1]
		timeStr := matches[2]
		var ampmStr string
		if len(matches) > 3 && matches[3] != "" {
			ampmStr = matches[3]
		}

		dateParts := orderDateComponents(dateStr)

		var day, month, year string
		if daysFirst {
			day, month, year = dateParts[0], dateParts[1], dateParts[2]
		} else {
			month, day, year = dateParts[0], dateParts[1], dateParts[2]
		}

		normalizedDate := normalizeDate(year, month, day)
		year, month, day = normalizedDate[0], normalizedDate[1], normalizedDate[2]

		var normalizedTime string
		if ampmStr != "" {
			normalizedTime = normalizeTime(convertTime12to24(timeStr, normalizeAMPM(ampmStr)))
		} else {
			normalizedTime = normalizeTime(timeStr)
		}

		timeParts := strings.Split(normalizedTime, ":")
		hour, minute, second := timeParts[0], timeParts[1], timeParts[2]

		yearInt, _ := strconv.Atoi(year)
		monthInt, _ := strconv.Atoi(month)
		dayInt, _ := strconv.Atoi(day)
		hourInt, _ := strconv.Atoi(hour)
		minuteInt, _ := strconv.Atoi(minute)
		secondInt, _ := strconv.Atoi(second)

		date := time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, secondInt, 0, time.UTC)
		result[i].Date = date

		// Use the full raw message to extract the complete message text
		if rawMsg.System {
			prefixLen := len(matches[0]) - len(matches[4])
			result[i].Message = strings.TrimSuffix(rawMsg.Msg[prefixLen:], "\n")
		} else {
			// Extract the full message text, including newlines
			prefixLen := len(matches[0]) - len(matches[5])
			result[i].Message = strings.TrimSuffix(rawMsg.Msg[prefixLen:], "\n")
		}

		// Add attachment if requested
		if options.ParseAttachments {
			if i < len(result) {
				attachment := parseMessageAttachment(result[i].Message)
				if attachment != nil {
					result[i].Attachment = attachment
				}
			}
		}
	}

	for i := range min(10, len(result)) {
		if strings.Contains(result[i].Message, "end-to-end") {
			result[i].Author = nil
		}
	}

	return result, nil
}

// ParseString parses a string containing a WhatsApp chat log.
// Returns an array of parsed messages.
func ParseString(content string, options *ParseStringOptions) ([]Message, error) {
	if options == nil {
		defaultOptions := ParseStringOptions{
			ParseAttachments: false,
		}
		options = &defaultOptions
	}

	lines := newlinesRegex.Split(content, -1)
	rawMessages := makeArrayOfMessages(lines)
	return parseMessages(rawMessages, *options)
}
