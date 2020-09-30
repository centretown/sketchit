package api

import "testing"

func TestRoute(t *testing.T) {
	// empty
	route := buildRoute([]string{})
	if len(route) != 0 {
		t.Fatalf("route should be empty %v", route)
	}
	// one route
	route = buildRoute([]string{"devices"})
	if len(route) != 1 {
		t.Fatalf("route should have one entry %v", route)
	}
	// same route except as step
	otherRoute := buildRoute([]string{}, "devices")
	if len(otherRoute) != 1 {
		t.Fatalf("route should have one entry %v", otherRoute)
	}
	// both routes should be the same
	if !compareRoutes(route, otherRoute) {
		t.Fatalf("routes should be equal %v %v", route, otherRoute)
	}
	// make sure both routes are what we expect
	if !compareRoutes(route, []string{"devices"}) {
		t.Fatalf("routes should be equal %v %v", route, []string{"devices"})
	}
	// append to a route - step forward
	route = buildRoute([]string{"devices"}, "sector1")
	if !compareRoutes(route, []string{"devices", "sector1"}) {
		t.Fatalf("routes should be equal %v %v", route, []string{"devices", "sector1"})
	}
	// check compareRoutes for consistency
	if compareRoutes(route, []string{"devices", "sector2"}) {
		t.Fatalf("routes should not be equal %v %v", route, []string{"devices", "sector1"})
	}
	// step forward
	route = buildRoute(route, "esp32-01")
	if !compareRoutes(route,
		[]string{"devices", "sector1", "esp32-01"}) {

		t.Fatalf("routes should be equal %v %v", route,
			[]string{"devices", "sector1", "esp32-01"})
	}
	// step back
	route = buildRoute(route, ".")
	if !compareRoutes(route,
		[]string{"devices", "sector1"}) {

		t.Fatalf("routes should be equal %v %v", route,
			[]string{"devices", "sector1"})
	}
	// step forward
	route = buildRoute(route, "esp32-01")
	// two steps back
	route = buildRoute(route, "..")
	if !compareRoutes(route, []string{"devices"}) {
		t.Fatalf("routes should be equal %v %v", route, []string{"devices"})
	}
	// step back and step forward to new collections
	route = buildRoute(route, ".sketches")
	if !compareRoutes(route, []string{"sketches"}) {
		t.Fatalf("routes should be equal %v %v", route, []string{"sketches"})
	}
	// two steps forward
	route = buildRoute(route, "esp32kit", "blink")
	if !compareRoutes(route,
		[]string{"sketches", "esp32kit", "blink"}) {
		t.Fatalf("routes should be equal %v %v", route,
			[]string{"sketches", "esp32kit", "blink"})
	}
	// alternate route
	otherRoute = buildRoute([]string{}, "sketches.esp32kit.blink")
	// both routes should be the same
	if !compareRoutes(route, otherRoute) {
		t.Fatalf("routes should be equal %v %v", route, otherRoute)
	}
	// back to empty
	route = buildRoute(route, ".......")
	if len(route) != 0 {
		t.Fatalf("route should be empty %v", route)
	}
	// alternate route prefixed with a '.'
	route = buildRoute(route, ".sketches.esp32kit.blink")
	// both routes should be the same
	if !compareRoutes(route, otherRoute) {
		t.Fatalf("routes should be equal %v %v", route, otherRoute)
	}

	// two back, one forward
	route = buildRoute(route, "..nanokit")
	if !compareRoutes(route,
		[]string{"sketches", "nanokit"}) {
		t.Fatalf("routes should be equal %v %v", route,
			[]string{"sketches", "nanokit"})
	}

	// two forward, extra dot ignored
	route = buildRoute(route, "nano02..setup")
	if !compareRoutes(route,
		[]string{"sketches", "nanokit", "nano02", "setup"}) {
		t.Fatalf("routes should be equal %v %v", route,
			[]string{"sketches", "nanokit", "nano02", "setup"})
	}

	t.Log(route)
}

func compareRoutes(route []string, otherRoute []string) (status bool) {
	if len(route) != len(otherRoute) {
		return
	}
	for i, s := range route {
		if s != otherRoute[i] {
			return
		}
	}
	status = true
	return
}
