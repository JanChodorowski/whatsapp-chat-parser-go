package parser

import (
	"time"
)

func GetAuthorsFromMessages(messages *[]Message) []string {
	seen := make(map[string]bool)
	var uniqueAuthors []string

	for _, message := range *messages {
		if message.Author == nil {
			continue
		}

		author := message.Author

		if !seen[*author] {
			seen[*author] = true
			uniqueAuthors = append(uniqueAuthors, *author)
		}
	}

	return uniqueAuthors
}

func GetFirstAndLastMessageDates(messages *[]Message) (*time.Time, *time.Time) {
	if len(*messages) == 0 {
		return nil, nil
	}
	firstDate := (*messages)[0].Date
	lastDate := (*messages)[len(*messages)-1].Date
	return &firstDate, &lastDate
}
