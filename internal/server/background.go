package server

import (
	"github.com/cyc1ones/go-kit/background"
	"github.com/go-kratos/kratos/v2/log"
)

func NewBackground(logger log.Logger) *background.Background {
	bg := background.New(
		background.WithLogger(logger),
	)

	// TODO: register background goroutines
	return bg
}
