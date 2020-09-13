package api

// Reducer is the interface implemented by an object
// that can reduce itself to a summary or brief
type Reducer interface {
	Reduce(...Reduction) interface{}
}

// Reduction level
type Reduction int

// Reduction levels
const (
	// no reduction
	Full Reduction = iota
	// less detail than full
	Summary
	// name and value
	Brief
	ReductionMax
	ReductionNotFound Reduction = -1
)

var reductionKeys []string
var reductionValues map[string]Reduction

func init() {
	reductionKeys = make([]string, ReductionMax)
	reductionKeys[Full] = "full"
	reductionKeys[Summary] = "summary"
	reductionKeys[Brief] = "brief"

	reductionValues = make(map[string]Reduction, ReductionMax)
	for i, k := range reductionKeys {
		reductionValues[k] = Reduction(i)
	}
}

func (r Reduction) String() (key string) {
	if r == ReductionNotFound || r >= ReductionMax {
		return
	}
	key = reductionKeys[r]
	return
}

// NewReduction finds the key and returns a reduction level
// ReductionNotFound is returned when not found
func NewReduction(key string) (reduction Reduction) {
	reduction = ReductionNotFound
	r, ok := reductionValues[key]
	if ok {
		reduction = r
	}
	return
}
