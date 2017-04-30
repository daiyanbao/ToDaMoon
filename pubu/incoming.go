// Incoming message builder.
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
	Title       string `json:"title"`
	Description string `json:"descrption,omitempty"`
	Color       color  `json:"color,omitempty"`
	URL         string `json:"url,omitempty"`
}

//Build 构建了消息
func (m *Incoming) Build() (io.Reader, error) {
	if m.Text == "" {
		return nil, fmt.Errorf("text is required")
	}

	for _, attachment := range m.Attachments {
		if attachment.Title == "" && attachment.URL == "" {
			return nil, fmt.Errorf("title or URL is required")
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
