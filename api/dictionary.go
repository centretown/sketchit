package api

// Dictionary is keyed to data model sch
// string keys are dot separated
// type Dictionary map[string]*Model

// DictionaryNew -
func DictionaryNew(collections []*Collection) (dictionary *Dictionary) {
	dictionary = &Dictionary{}
	dictionary.Models = make(map[string]*Model)
	dictionary.Collections = make([]*Collection, len(collections))
	dictionary.Collections = append(dictionary.Collections, collections...)

	var f = func(s *Model, t Traveler) {
		t.Push(s)
		dictionary.Models[t.String()] = s
		t.Pop()
	}

	var km = KeyMaker{
		stack: make([]string, 0, KeyStackDepth),
	}

	for _, c := range collections {
		c.Model.Travel(&km, f)
	}
	return
}

// Present Presenter interface contract
func (dict *Dictionary) Present(presentation *Presentation) (s string) {
	for _, coll := range dict.Collections {
		s += coll.Present(presentation)
	}
	return
}
