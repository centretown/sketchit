package cmdr

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/centretown/sketchit/api"
	"github.com/centretown/sketchit/auth"
	"golang.org/x/net/context"
)

var snakeoil = "../" + auth.SnakeOil

// TestDictionary -
func testDictionary(t *testing.T) {
	a := &auth.Authentication{
		Login:    "testing",
		Password: "test",
	}

	conn, err := auth.Connect(snakeoil, a)
	if err != nil {
		t.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewSketchitClient(conn)
	ctx := context.Background()
	response, err := client.SayHello(ctx, &api.PingMessage{Greeting: "foo"})
	if err != nil {
		t.Fatalf("Error when calling SayHello: %s", err)
	}
	t.Logf("Response from server: %s", response.Greeting)

	creq := &api.ListCollectionsRequest{}
	cres, err := client.ListCollections(ctx, creq)
	if err != nil {
		t.Fatalf("Error when calling ListCollections: %s", err)
	}

	dictionary := api.DictionaryNew(cres.Collections)

	sch := dictionary[".sketches"]
	s, err := Marshal(sch, api.YAML, api.Full)
	if err != nil {
		t.Fatalf("marshal: %s, yaml, .sketches", err)
	}
	t.Log(s)

	sch = dictionary[".devices"]
	if err != nil {
		t.Fatalf("marshal: %s, yaml, .devices", err)
	}
	t.Log(s)

	sch = dictionary[".sketches"]
	s, err = Marshal(sch, api.YAML, api.Brief)
	if err != nil {
		t.Fatalf("marshal: %s, json, .sketches", err)
	}
	t.Log(s)

	sch = dictionary[".devices"]
	s, err = Marshal(sch, api.XML, api.Summary)
	if err != nil {
		t.Fatalf("marshal: %s, xml, .devices", err)
	}
	t.Log(s)
}

func TestSkills(t *testing.T) {

	a := &auth.Authentication{
		Login:    "testing",
		Password: "test",
	}

	conn, err := auth.Connect(snakeoil, a)
	if err != nil {
		t.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewSketchitClient(conn)
	ctx := context.Background()
	cmdr := New(ctx, conn, client)
	cmdr.Build()

	// testHello(t, cmdr)
	// testHelp(t, cmdr)
	// testFlags(t, cmdr)
	testReducer(t, cmdr)
}

func testHello(t *testing.T, cmdr *Deputy) {
	fv := cmdr.Flags.Values()
	name := "hello"
	c, ok := cmdr.Aliases[name]
	if !ok {
		t.Fatalf("%s command not found.", name)
	}
	s, err := c.Run(fv, "sketchit")
	if err != nil {
		t.Fatalf("hello reported: %v", err)
	}
	t.Log(s)
}

func testHelp(t *testing.T, cmdr *Deputy) {
	fv := cmdr.Flags.Values()
	name := "help"
	var err error
	var errFunc = func() {
		if err != nil {
			t.Fatalf("%s reported: %v", name, err)
		}
	}
	defer errFunc()

	c, ok := cmdr.Aliases[name]
	if !ok {
		t.Fatalf("%s command not found.", name)
	}

	s, err := c.Run(fv)
	if err != nil {
		return
	}

	s, err = c.Run(fv, "all")
	if err != nil {
		return
	}
	t.Log(s)

	s, err = c.Run(fv, "list")
	if err != nil {
		return
	}
	pre := "Topic: list"
	if !strings.HasPrefix(s, pre) {
		err = fmt.Errorf("got '%v' expected '%v'", s, pre)
		return
	}
	t.Log(s)

	s, err = c.Run(fv, "foo")
	if err == nil {
		err = fmt.Errorf("uncaught unknown command: %v", "foo")
		return
	}
	err = nil
	t.Log(s)
}

func testFlags(t *testing.T, cmdr *Deputy) {
	var err error
	var s string

	cmd, flags, args, err := cmdr.Parse("flags -f=json -d:summary")
	if err != nil {
		t.Fatal(err)
	}

	cmd, flags, args, err = cmdr.Parse("flags -f=xmls -d:brief")
	if err == nil {
		t.Fatal("invalid -f flag value accepted")
	}
	t.Log(err)

	cmd, flags, args, err = cmdr.Parse("flags -f=xml -d:brieff")
	if err == nil {
		t.Fatal("invalid -d flag value accepted")
	}
	t.Log(err)

	cmd, flags, args, err = cmdr.Parse("flags -g=xml -d:brief")
	if err == nil {
		t.Fatal("invalid -d flags accepted")
	}
	t.Log(err)

	cmd, flags, args, err = cmdr.Parse("flags -f=xml -d:brief")
	if err != nil {
		t.Fatal(err)
	}

	cmd, flags, args, err = cmdr.Parse("list -f=json -d:brief cottage")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("list command", args, flags)

	s, err = cmd.Run(flags, args...)
	if err != nil {
		return
	}
	t.Log(s)

	// if devKind == reflect.Slice {
	// 	t.Log("its a slice")
	// 	elemType := devType.Elem()
	// 	t.Log(devType.Elem())
	// 	if elemType.Implements(reflect.Type()) {
	// 		t.Fatal("device not a reducer")
	// 	} else {
	// 		t.Log("device is a reducer")
	// 	}
	// }

}
func readTag(e interface{}, name string) (tag reflect.StructTag, ok bool) {
	t := reflect.TypeOf(e)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	f, ok := t.FieldByName(name)
	tag = f.Tag
	return
}

func testReducer(t *testing.T, cmdr *Deputy) {
	device := &api.Device{}
	devices := make([]*api.Device, 0)
	devices = append(devices, device)
	d := &api.DevicesReducer{Devices: devices}
	tag, _ := readTag(d, "Devices")
	t.Log("tag", tag)
	// t, ok := devices.(api.DeviceSlice)

	// var o interface{}
	// o = devices
	// ovalue := reflect.ValueOf(o)
	// ovalue.Interface()
	// otype := ovalue.Type()
	// okind := otype.Kind()
	// oelem := otype.Elem()
	// reducerType := reflect.TypeOf((*api.Reducer)(nil)).Elem()
	// ok := oelem.Implements(reducerType)
	// t.Log("implements reducer", ok, oelem)
	// t.Log("reducerType", reducerType)
	// t.Logf("v=%v t=%v, k=%v, e=%v ek=%v",
	// 	ovalue, otype, okind, oelem, oelem.Kind())
	// if ok {
	// 	m, ok := oelem.MethodByName("SetReduction")
	// 	t.Logf("method:%v,ok:%v", m, ok)
	// 	if ok {
	// 		t.Log("SetReduction")
	// 		f := m.Func
	// 		levels := []api.Reduction{api.Full, api.Brief}
	// 		lv := reflect.ValueOf(levels)
	// 		rlv := []reflect.Value{lv}
	// 		t.Logf("func=%v lv%v, rlv%v", f, lv, rlv)
	// 		f.Call(rlv)
	// 	}
	// }

}
