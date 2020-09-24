package api

// Runner is returned when a command request has been successfully parsed
type Runner struct {
	presentation *Presentation
	route        []string
	taskRunner   *TaskRunner
}

// Run -
func (r *Runner) Run() (s string, err error) {
	return r.taskRunner.run(r.presentation, r.route...)
}

// TaskRunner -
type TaskRunner struct {
	run func(pr *Presentation, route ...string) (s string, err error)
}

func (resp *Responder) buildRunners() {
	resp.runners = make([]*TaskRunner, Task_TASK_LEN)

	resp.runners[Task_exit] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_help] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_list] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_goto] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_save] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_remove] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_hello] = &TaskRunner{
		run: func(pr *Presentation, route ...string) (s string, err error) {
			return
		},
	}
}
