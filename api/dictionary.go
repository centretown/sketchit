package api

// Dictionary is keyed to data model sch
// string keys are dot separated
// type Dictionary map[string]*Schema

// DictionaryNew -
func DictionaryNew(collections []*Collection) (dictionary *Dictionary) {
	dictionary = &Dictionary{}
	dictionary.DictionaryMap = make(map[string]*Schema)
	dictionary.Collections = make([]*Collection, len(collections))
	dictionary.Collections = append(dictionary.Collections, collections...)

	var f = func(s *Schema, t Traveler) {
		t.Push(s)
		dictionary.DictionaryMap[t.String()] = s
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

// Reduce interface support
func (dictionary *Dictionary) Reduce(projection ...Projection) (out interface{}) {

	var f = func(s *Schema, t Traveler) {
		t.Push(s)
		s.Reduce(projection...)
		t.Pop()
	}

	var km = KeyMaker{
		stack: make([]string, 0, maxStackDepth),
	}

	for _, c := range dictionary.Collections {
		c.Schema.Travel(&km, f)
	}
	return
}

func (coll *Collection) Reduce() {

}
