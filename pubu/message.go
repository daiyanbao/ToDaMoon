package pubu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type color string

const (
	warning color = "warning"
	info    color = "info"
	primary color = "primary"
	mistake color = "error"
	muted   color = "muted"
	success color = "success"
)

//incoming 消息格式
type incoming struct {
	Text        string       `json:"text"`
	Attachments []attachment `json:"attachments,omitempty"`
}

//incomingAttachment 附件格式
type attachment struct {
	Text  string `json:"title"`
	Color color  `json:"color,omitempty"`
}

//build 构建了消息
func (m *incoming) build() (io.Reader, error) {
	if m.Text == "" {
		return nil, fmt.Errorf("text is required")
	}

	for _, attachment := range m.Attachments {
		if attachment.Text == "" {
			return nil, fmt.Errorf("text is required")
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

//msg 用于创建消息。
func msgMaker(text string) *incoming {
	return &incoming{
		Text: text,
	}
}

//att 创建带附件的消息。
func att(text string, color color) attachment {
	return attachment{
		Text:  text,
		Color: color,
	}
}

func (m *incoming) debug() {
	m.Attachments = append(m.Attachments, att("DEBUG: "+time.Now().Format("2006-01-02 15:04:05"), muted))
}
func (m *incoming) warning() {
	m.Attachments = append(m.Attachments, att("WARNING: "+time.Now().Format("2006-01-02 15:04:05"), warning))
}

func (m *incoming) error() {
	m.Attachments = append(m.Attachments, att("ERROR: "+time.Now().Format("2006-01-02 15:04:05"), mistake))
}

func (m *incoming) info() {
	m.Attachments = append(m.Attachments, att("INFO: "+time.Now().Format("2006-01-02 15:04:05"), info))
}

func (m *incoming) good() {
	m.Attachments = append(m.Attachments, att("GOOD: "+time.Now().Format("2006-01-02 15:04:05"), success))
}

func (m *incoming) bad() {
	m.Attachments = append(m.Attachments, att("BAD: "+time.Now().Format("2006-01-02 15:04:05"), primary))
}
