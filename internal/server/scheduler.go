package server

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
)

// Scheduler is a cron wrapper
type Scheduler struct {
	c   *cron.Cron
	log *log.Helper
}

var _ transport.Server = (*Scheduler)(nil)

func NewScheduler(conf *conf.Scheduler, logger log.Logger) *Scheduler {
	opts := []cron.Option{}
	if conf.WithSeconds {
		opts = append(opts, cron.WithSeconds())
	}

	c := cron.New(opts...)
	// TODO: register cron jobs

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
