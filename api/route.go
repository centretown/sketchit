package api

import "strings"

// Routes are lookup keys to the Dictionary, where Models are
// stored. The Models define the elements of each Collection
// including data types, constraints and descriptions.
//
// Collections include data documents, the data model and the
// the responder api. The api provides access to the sketchit
// services, sketches, devices and tookits.
//
// Each route is a sequence of Labels separated by periods '.'
// or forward slashes '/'. Labels are unique within a
// Collection, should have meaning and must exclude separators.
//
// Route structure:
// <root>. single separator
//   <collection>. collection label
//     <parent>. parent label route to ancestors
//       <label>.
//         <model>...
// <model>
//   <items>... values, arrays, maps and models
var routeSeparator = "."

// buildRoute takes an existing route (could be empty) and
// applies a list of routing steps to calculate a new route.
// Steps may be single labels separated by spaces or periods.
// Leading periods are back steps. Each following label is a
// step forward.
func buildRoute(routeIn []string, stepsIn ...string) (route []string) {
	route = make([]string, len(routeIn))
	copy(route, routeIn)
	if len(stepsIn) < 1 {
		// copy existing route and return it if no new steps
		return
	}

	// copy and replace /'s with .'s
	steps := make([]string, len(stepsIn))
	for i := range steps {
		steps[i] = strings.Replace(stepsIn[i], "/", routeSeparator, -1)
	}

	leadStep := steps[0]
	backSteps := getLeadingDots(leadStep)
	routeLength := len(route)
	if routeLength < backSteps {
		// discard excess back steps
		backSteps = routeLength
	}

	// take steps back
	if backSteps > 0 {
		// reset or remove first step
		if backSteps == len(leadStep) {
			steps = steps[1:]
		} else {
			steps[0] = leadStep[backSteps:]
		}
		route = route[:routeLength-backSteps]
	}

	// separate remaining arguments from their dots
	var tokens []string
	for _, step := range steps {
		tokens = append(tokens,
			strings.Split(step, routeSeparator)...)
	}

	// take steps forward
	for _, token := range tokens {
		// ignore zero length tokens resulting from multiple dot separators
		if len(token) > 0 {
			route = append(route, token)
		}
	}
	return
}

func getLeadingDots(s string) (count int) {
	for _, ch := range s {
		if ch != '.' {
			return
		}
		count++
	}
	return
}
