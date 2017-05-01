package pubu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

//Incoming 消息格式
type Incoming struct {
	Text        string               `json:"text"`
	Attachments []IncomingAttachment `json:"attachments,omitempty"`
}

//IncomingAttachment 附件格式
type IncomingAttachment struct {
	Text  string `json:"title"`
	Color color  `json:"color,omitempty"`
}

//Build 构建了消息
func (m *Incoming) Build() (io.Reader, error) {
	if m.Text == "" {
		return nil, fmt.Errorf("text is required")
	}

	for _, attachment := range m.Attachments {
		if attachment.Text == "" {
			return nil, fmt.Errorf("title or URL is required")
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
