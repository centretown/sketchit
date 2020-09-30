package api

import (
	context "context"
	"errors"
	"fmt"
	"strings"

	"github.com/centretown/sketchit/info"
	"github.com/golang/glog"
	grpc "google.golang.org/grpc"
)

// Responder -
type Responder struct {
	ctx          context.Context
	conn         *grpc.ClientConn
	client       SketchitClient
	deputy       *Deputy
	presentation *Presentation
	suffix       string
	runners      []*TaskRunner
	route        []string
}

var defaultPresentation = &Presentation{
	Format:     Format_yaml,
	Projection: []Projection{Projection_full},
	Confirm:    Auto_off,
}

// NewResponder -
func NewResponder(ctx context.Context,
	conn *grpc.ClientConn,
	client SketchitClient) (responder *Responder) {

	responder = &Responder{
		ctx:          ctx,
		client:       client,
		conn:         conn,
		presentation: defaultPresentation,
		suffix:       " - ",
	}
	responder.build()
	return
}

// Prompt returns the prompt string
func (resp *Responder) Prompt() (p string) {
	dot := "."
	dotDir := strings.Join(resp.route, dot)
	p += dot + dotDir + resp.suffix
	return
}

// Build the responder
func (resp *Responder) build() {
	var err error
	dRequest := &GetDeputyRequest{Name: "Andy"}
	resp.deputy, err = resp.client.GetDeputy(resp.ctx, dRequest)
	if err != nil {
		glog.Fatalf("Error when calling GetDeputy: %v", err)
	}

	cRequest := &ListCollectionsRequest{}
	cResponse, err := resp.client.ListCollections(resp.ctx, cRequest)
	if err != nil {
		glog.Fatalf("Error when calling ListCollections: %v", err)
	}
	resp.deputy.Collections = cResponse.Collections

	resp.deputy.BuildSkillset()
	resp.deputy.BuildGallery()
	resp.deputy.BuildDictionary()

	resp.buildRunners()
	return
}

// parsing information
var (
	ErrEmpty              = errors.New("no input")
	ErrSkillNotFound      = errors.New("command not found")
	ErrFlagNotFound       = errors.New("flag not found")
	ErrFormatNotFound     = errors.New("format value undefined")
	ErrProjectionNotFound = errors.New("projection value undefined")
	ErrConfirmNotFound    = errors.New("auto confirm value undefined")
	ErrExit               = errors.New("exit")
	ErrFeatureNotFound    = errors.New("feature not found")
)

// Parse -
func (resp *Responder) Parse(input string) (runner *Runner, err error) {
	runner = &Runner{}
	s := strings.TrimSpace(input)
	if len(s) < 1 {
		err = ErrEmpty
		return
	}

	args := strings.Fields(s)
	verb := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	task, ok := Task_value[verb]
	if !ok {
		err = info.Inform(err, ErrSkillNotFound, verb)
		return
	}
	runner.presentation, runner.steps, err = resp.parseFlags(args)
	if err != nil {
		return
	}
	runner.taskRunner = resp.runners[task]
	return
}

// parseFlags scans input arguments for presentation flags
// returns presentation
// errors are flagged for invalid flags or values
func (resp *Responder) parseFlags(input []string) (presentation *Presentation, steps []string, err error) {
	// assume no flags all routes
	steps = make([]string, 0, len(input))
	presentation = &Presentation{
		Format:     resp.presentation.Format,
		Projection: resp.presentation.Projection,
		Confirm:    resp.presentation.Confirm,
	}

	var (
		prefix     = "-"                // flag indicator
		separators = []string{"=", ":"} // assignment operators
		pair       = []string{}         // expression pair
	)

	for _, token := range input {
		if !strings.HasPrefix(token, prefix) {
			// not a flag add to route
			steps = append(steps, token)
			continue
		}

		token = strings.TrimPrefix(token, prefix)
		for _, a := range separators {
			pair = strings.Split(token, a)
			if len(pair) > 1 {
				break
			}
		}

		if len(pair) < 2 {
			// no assigment ignore
			continue
		}

		key, value := pair[0], pair[1]
		flag, ok := Feature_Flag_value[key]
		if !ok {
			err = info.Inform(err, ErrFlagNotFound, fmt.Sprintf("flag: %v", key))
			return
		}

		switch Feature_Flag(flag) {
		case Feature_f:
			presentation.Format, err = resp.parseFormat(value)
		case Feature_d:
			presentation.Projection, err = resp.parseProjection(value)
		case Feature_auto:
			presentation.Confirm, err = resp.parseConfirm(value)
		}
		if err != nil {
			return
		}
	}

	return
}

func getKeys(m map[string]int32) (keys []string) {
	keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func (resp *Responder) parseFormat(value string) (format Format, err error) {
	f, ok := Format_value[value]
	if !ok {
		err = info.Inform(err, ErrFormatNotFound,
			fmt.Sprintf(":%v expecting %v", value, getKeys(Format_value)))
		return
	}
	format = Format(f)
	return
}

func (resp *Responder) parseConfirm(value string) (confirm Auto, err error) {
	a, ok := Auto_value[value]
	if !ok {
		err = info.Inform(err, ErrConfirmNotFound,
			fmt.Sprintf(":%v expecting %v", value, getKeys(Auto_value)))
		return
	}
	confirm = Auto(a)
	return
}

func (resp *Responder) parseProjection(value string) (projections []Projection, err error) {
	tokens := strings.Split(value, ",")
	for _, token := range tokens {
		p, ok := Projection_value[token]
		if !ok {
			err = info.Inform(err, ErrProjectionNotFound,
				fmt.Sprintf(":%v expecting %v", value, getKeys(Projection_value)))
			return
		}
		projections = append(projections, Projection(p))
	}
	return
}
