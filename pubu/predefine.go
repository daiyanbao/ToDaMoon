package pubu

//Msg 用于创建消息。
func Msg(text string) *Incoming {
	return &Incoming{
		Text: text,
	}
}

//Att 创建带附件的消息。
func att(title, description, url string, color color) *IncomingAttachment {
	return &IncomingAttachment{
		Title:       title,
		Description: description,
		URL:         url,
		Color:       color,
	}
}

//Warning 添加WARNING消息。
func (ic *Incoming) Warning(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, warning))
}

//Info 添加Info消息。
func (ic *Incoming) Info(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, info))
}

//Primary 添加Primary消息。
func (ic *Incoming) Primary(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, primary))
}

//Mistake 添加error消息
func (ic *Incoming) Mistake(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, mistake))
}

//Muted 添加Muted消息
func (ic *Incoming) Muted(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, muted))
}

//Success 添加Success消息
func (ic *Incoming) Success(title, description, url string) {
	ic.Attachments = append(ic.Attachments, *att(title, description, url, success))
}
