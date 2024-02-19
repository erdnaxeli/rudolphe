package matrix

import "context"

func (c Client) StartSync(ctx context.Context) error {
	return c.client.SyncWithContext(ctx)
}
