package api

// Presenter deputy's Print/Marshall functions use this interface
// to project and format list items
type Presenter interface {
	Present(presentation *Presentation) string
}
