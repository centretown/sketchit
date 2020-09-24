package api

// BuildSkillset map
func (dep *Deputy) BuildSkillset() {
	for _, skill := range dep.Skills {
		dep.Skillset[int32(skill.Task)] = skill
	}
}

// BuildGallery map
func (dep *Deputy) BuildGallery() {
	for _, feature := range dep.Features {
		dep.Gallery[int32(feature.Flag)] = feature
	}
}

// BuildDictionary -
func (dep *Deputy) BuildDictionary() {
	dep.Dictionary = make(map[string]*Model)

	var f = func(s *Model, t Traveler) {
		t.Push(s)
		dep.Dictionary[t.String()] = s
		t.Pop()
	}

	var keyMaker = KeyMaker{
		stack: make([]string, 0, KeyStackDepth),
	}

	for _, c := range dep.Collections {
		c.Model.Travel(&keyMaker, f)
	}
	return
}

// Present Presenter interface contract
func (dep *Deputy) Present(presentation *Presentation) (s string) {
	for _, coll := range dep.Collections {
		s += coll.Present(presentation)
	}
	return
}

// cmd := strings.TrimSuffix(input, "\n")
// if input == "" {
// 	err = ErrEmpty
// 	return
// }

// // not empty so at least one
// args := strings.Fields(cmd)
// verb := args[0]
// if len(args) > 1 {
// 	args = args[1:]
// } else {
// 	args = []string{}
// }

// // skill, ok := resp.deputy.Skillset[verb]
// handler, ok := resp.handlers[verb]
// if !ok {
// 	err = info.Inform(err, ErrSkillNotFound, verb)
// 	return
// }
