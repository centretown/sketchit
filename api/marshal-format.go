package api

// MarshalFormat -
type MarshalFormat int

// MarshalFormat constants
const (
	XML MarshalFormat = iota
	JSON
	YAML
	MarshallMax
	FormatNotFound MarshalFormat = -1
)

var formatKeys []string
var formatValues map[string]MarshalFormat

func init() {
	formatKeys = make([]string, MarshallMax)
	formatKeys[XML] = "xml"
	formatKeys[JSON] = "json"
	formatKeys[YAML] = "yaml"

	formatValues = make(map[string]MarshalFormat, MarshallMax)
	for i, k := range reductionKeys {
		formatValues[k] = MarshalFormat(i)
	}
}

func (r MarshalFormat) String() (key string) {
	if r == FormatNotFound {
		return
	}
	key = formatKeys[r]
	return
}

// NewMarshalFormat finds the key and returns a format level
// FormatNotFound is returned when not found
func NewMarshalFormat(key string) (format MarshalFormat) {
	format = FormatNotFound
	r, ok := formatValues[key]
	if ok {
		format = r
	}
	return
}
