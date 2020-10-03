package api

import (
	"errors"
	"fmt"

	"github.com/centretown/sketchit/info"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Runner is returned when a command request has been successfully parsed
type Runner struct {
	presentation *Presentation
	steps        []string
	taskRunner   *TaskRunner
}

// Run -
func (r *Runner) Run() (s string, err error) {
	return r.taskRunner.run(r.presentation, r.steps...)
}

// TaskRunner -
type TaskRunner struct {
	run func(pr *Presentation, steps ...string) (s string, err error)
}

// Error messsages
var (
	ErrHello         = errors.New("failed to greet server")
	ErrList          = errors.New("failed to list devices")
	ErrGet           = errors.New("failed to get device")
	ErrNotEnoughArgs = errors.New("not enough arguments")
	ErrDecode        = errors.New("failed to decode")
)

func (resp *Responder) buildRunners() {
	resp.runners = make([]*TaskRunner, len(Task_value))

	resp.runners[Task_exit] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_help] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_list] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			route := buildRoute(resp.route, steps...)
			accessor, err := NewAccessor(route...)
			if err != nil {
				err = info.Inform(err, ErrList, steps)
				return
			}
			request := &ListRequest{Parent: accessor.Parent}
			response, err := resp.client.List(resp.ctx, request)
			if err != nil {
				err = info.Inform(err, ErrList, steps)
				return
			}
			if len(response.Items) < 1 {
				fmt.Printf("Nothing to list for %v\n", steps)
				return
			}

			items := make([]protoreflect.ProtoMessage, len(response.Items))
			for i, any := range response.Items {
				items[i] = accessor.MakeItem()
				any.UnmarshalTo(items[i])
			}
			b, err := Marshal(items, pr)
			if err == nil {
				s = string(b)
			}
			return
		},
	}
	resp.runners[Task_goto] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			resp.route = buildRoute(resp.route, steps...)
			return
		},
	}
	resp.runners[Task_save] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_remove] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			return
		},
	}
	resp.runners[Task_hello] = &TaskRunner{
		run: func(pr *Presentation, steps ...string) (s string, err error) {
			message := "hello " + fmt.Sprintln(steps)
			response, err := resp.client.SayHello(resp.ctx, &PingMessage{Greeting: "hello"})
			if err != nil {
				info.Inform(err, ErrHello, message)
				return
			}
			b, err := Marshal(response, pr)
			if err == nil {
				s = string(b)
			}
			return
		},
	}
}
