package matrix

func (b Bot) SendText(text string) error {
	err := b.sendText(b.roomID, text)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}
