package service

import "context"

type Scheduler interface {
	Run(ctx context.Context) error
}
