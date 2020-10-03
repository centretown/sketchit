package api

import "testing"

func TestAccessor(t *testing.T) {

	route := buildRoute([]string{"devices"})
	a, _ := NewAccessor(route...)
	t.Log(a.Parent)

	route = buildRoute(route, "home-iot")
	a, _ = NewAccessor(route...)
	t.Log(a.Parent)

	route = buildRoute(route, "..sketches.ESP32")
	a, _ = NewAccessor(route...)
	t.Log(a.Parent)

	route = buildRoute(route, "//devices.home-iot")
	a, _ = NewAccessor(route...)
	t.Log(a.Parent)

	route = buildRoute(route, "//sketches/ESP32")
	a, _ = NewAccessor(route...)
	t.Log(a.Parent)
}
