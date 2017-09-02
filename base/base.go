package base

import (
	"sync"

	"golang.org/x/net/context"
)

var isDev = false

func SetDevelopmentMode(dev bool) {
	isDev = dev
}

func IsDev() bool {
	return isDev
}

func Parallel(ctx context.Context, jobs ...func(context.Context) error) []error {
	wg := sync.WaitGroup{}
	wg.Add(len(jobs))
	errors := make([]error, len(jobs))

	for i := range jobs {
		job := jobs[i]
		go func(i int) {
			defer wg.Done()
			errors[i] = job(ctx)
		}(i)
	}

	wg.Wait()
	return errors
}
