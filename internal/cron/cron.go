package cron

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Runner interface {
	Start(ctx context.Context)
}

type Job interface {
	Run(ctx context.Context) error
	GetDelay() time.Duration
	GetName() string
}

type cron struct {
	jobs []Job
}

func NewCron(jobs []Job) cron {
	return cron{
		jobs: jobs,
	}
}

func (c cron) Start(ctx context.Context) {
	for _, job := range c.jobs {
		j := job

		go func() {
			wg := &sync.WaitGroup{}

			for {
				wg.Add(1)

				go func() {
					defer wg.Done()

					select {
					case <-ctx.Done():
						log.WithFields(log.Fields{
							"job": j.GetName(),
						}).Debug("cron context done")
						return
					default:
						log.WithFields(log.Fields{
							"job": j.GetName(),
						}).Debug("started job")

						err := j.Run(ctx)
						if err != nil {
							log.WithFields(log.Fields{
								"job": j.GetName(),
							}).Error(err)
						}

						log.WithFields(log.Fields{
							"job": j.GetName(),
						}).Debug("finished job")

						time.Sleep(j.GetDelay())

						log.WithFields(log.Fields{
							"job": j.GetName(),
						}).Debug("waited job delay")

						return
					}
				}()

				wg.Wait()
			}
		}()
	}
}
