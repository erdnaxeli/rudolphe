package matrix

func (b Bot) StartSync() error {
	return b.client.Sync()
}
