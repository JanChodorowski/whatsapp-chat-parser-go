package parser

import "time"

type Message struct {
	Date       time.Time   `json:"date"`
	Author     *string     `json:"author"` // nil for system messages
	IsSystem   bool        `json:"isSystem"`
	Message    string      `json:"message"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	FileName string `json:"fileName"`
}

type RawMessage struct {
	System bool
	Msg    string
}

type ParseStringOptions struct {
	DaysFirst        *bool `json:"daysFirst"`
	ParseAttachments bool  `json:"parseAttachments"`
}
