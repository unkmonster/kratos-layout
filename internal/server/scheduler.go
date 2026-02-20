package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c   *cron.Cron
	log *log.Helper
}

var _ transport.Server = (*Scheduler)(nil)

func NewScheduler(logger log.Logger) *Scheduler {
	c := cron.New()

	return &Scheduler{
		c:   c,
		log: log.NewHelper(log.With(logger, "module", "scheduler")),
	}
}

// Start implements [transport.Server].
func (s *Scheduler) Start(context.Context) error {
	s.c.Start()
	s.log.Infof("[Scheduler] cron launched")
	return nil
}

// Stop implements [transport.Server].
func (s *Scheduler) Stop(ctx context.Context) error {
	s.log.Infof("[Scheduler] cron stopping")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.c.Stop().Done():
		return nil
	}
}
