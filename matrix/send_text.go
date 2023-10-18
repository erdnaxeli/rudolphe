package matrix

func (c Client) SendText(text string) error {
	err := c.sendText(c.roomID, text)
	if err != nil {
		return ErrSendMessage
	}

	return nil
}
