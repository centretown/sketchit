package api

// Dictionary is keyed to data model sch
// string keys are dot separated
type Dictionary map[string]*Schema

// DictionaryNew -
func DictionaryNew(collections []*Collection) (dict Dictionary) {
	dict = make(Dictionary)

	var f = func(s *Schema, t Traveler) {
		t.Push(s)
		dict[t.String()] = s
		t.Pop()
	}

	var km = KeyMaker{
		stack: make([]string, 0, maxStackDepth),
	}

	for _, c := range collections {
		c.Schema.Travel(&km, f)
	}
	return
}
