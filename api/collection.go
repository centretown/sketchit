package api

// Present Presenter interface contract
func (coll *Collection) Present(presentation *Presentation) string {
	return coll.Model.Present(presentation)
}
