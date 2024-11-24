package task

import (
	"github.com/lesomnus/bring/config"
	"github.com/lesomnus/bring/thing"
)

type Job struct {
	NumTasks int
}

type Task struct {
	BringConfig config.BringConfig

	Thing thing.Thing
	Job   Job
	Order int
	Dest  string
}
