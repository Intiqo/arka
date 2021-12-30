package scheduler

import "github.com/robfig/cron/v3"

type cronManager struct {
	cron *cron.Cron
}

func (c *cronManager) AddFunc(interval string, handler func()) (cron.EntryID, error) {
	return c.cron.AddFunc(interval, handler)
}
