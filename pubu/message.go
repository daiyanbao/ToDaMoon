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
package pubu

//Msg 用于创建消息。
func Msg(text string) *Incoming {
	return &Incoming{
		Text: text,
	}
}

//Att 创建带附件的消息。
func att(text string, color color) IncomingAttachment {
	return IncomingAttachment{
		Text:  text,
		Color: color,
	}
}

//Warning 添加WARNING消息。
func (ic *Incoming) Warning(text string) {
	ic.Attachments = append(ic.Attachments, att("WARNING: "+text, warning))
}

//Info 添加Info消息。
func (ic *Incoming) Info(text string) {
	ic.Attachments = append(ic.Attachments, att("INFO: "+text, info))
}

//Primary 添加Primary消息。
func (ic *Incoming) Primary(text string) {
	ic.Attachments = append(ic.Attachments, att("PRIMARY: "+text, primary))
}

//Mistake 添加error消息
func (ic *Incoming) Error(text string) {
	ic.Attachments = append(ic.Attachments, att("ERROR: "+text, mistake))
}

//Muted 添加Muted消息
func (ic *Incoming) Muted(text string) {
	ic.Attachments = append(ic.Attachments, att("MUTED: "+text, muted))
}

//Success 添加Success消息
func (ic *Incoming) Success(text string) {
	ic.Attachments = append(ic.Attachments, att("SUCCESS: "+text, success))
}