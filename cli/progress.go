package cli

import (
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type CalculateTask struct {
	mu             sync.Mutex
	taskName       string
	totalTasks     int
	completedTasks int
	tickInterval   time.Duration
	options        []progressbar.Option
}

func NewCalculateTask(taskName string, totalTasks int, tickInterval time.Duration, options ...progressbar.Option) *CalculateTask {
	return &CalculateTask{
		taskName:       taskName,
		totalTasks:     totalTasks,
		completedTasks: 0,
		tickInterval:   tickInterval,
		options:        options,
	}
}

func (task *CalculateTask) RefreshName(name string) {
	task.mu.Lock()
	defer task.mu.Unlock()
	task.taskName = name
}

func (task *CalculateTask) RefreshProgress(progress int) {
	task.mu.Lock()
	defer task.mu.Unlock()
	task.completedTasks = progress
}

func PrintProgress(task *CalculateTask) {
	done := make(chan struct{}, 1)
	if task.tickInterval == 0 {
		task.tickInterval = time.Millisecond * 65
	}

	// create new progress bar with default options
	options := []progressbar.Option{
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription(task.taskName),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionFullWidth(),
		progressbar.OptionOnCompletion(func() { done <- struct{}{} }),
		progressbar.OptionShowCount(),
	}
	options = append(options, task.options...)

	// initialize the progress bar and print it
	bar := progressbar.NewOptions(task.totalTasks, options...)
	for {
		select {
		case <-done:
			return
		case <-time.After(task.tickInterval):
			task.mu.Lock()
			if bar.State().Description != task.taskName {
				bar.Describe(task.taskName)
			}

			e := bar.Set(task.completedTasks)
			task.mu.Unlock()
			if e != nil {
				return
			}
		}
	}
}
