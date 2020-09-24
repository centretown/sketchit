package api

// Projector -
type Projector interface {
	Full() (o interface{})
	Summary() (o interface{})
	Brief() (o interface{})
}

// Project -
func Project(o interface{}, p ...Projection) interface{} {
	pr, ok := o.(Projector)
	if ok {
		switch p[0] {
		case Projection_full:
			return pr.Full()
		case Projection_summary:
			return pr.Summary()
		case Projection_brief:
			return pr.Brief()
		}
	}
	return o
}
