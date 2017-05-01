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
