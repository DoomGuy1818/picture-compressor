package PictureWorker

import (
	"log/slog"
	"picCompressor/internal/lib/compressor"
	"sync"
)

type Job struct {
	Path  string
	Alias string
	Reply chan Result
	Wg    *sync.WaitGroup
}

type Result struct {
	Path  string
	Error error
}

type Pool struct {
	jobs chan Job
}

func New(workers int, c compressor.Compressor, log *slog.Logger) *Pool {
	p := &Pool{
		jobs: make(chan Job),
	}

	for i := 0; i < workers; i++ {
		go worker(i, p.jobs, c, log)
	}

	return p
}

func (p *Pool) Enqueue(job Job) {
	p.jobs <- job
}

func worker(id int, jobs <-chan Job, c compressor.Compressor, log *slog.Logger) {
	for job := range jobs {
		log.Info(
			"compress image",
			"worker", id,
			"path", job.Path,
		)

		out, err := c.Compress(job.Path, job.Alias)

		res := Result{
			Path:  out,
			Error: err,
		}

		job.Reply <- res

		job.Wg.Done()
	}
}
