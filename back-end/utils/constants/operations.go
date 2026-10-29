package constants

const (
	Equal        string = "eq"
	NotEqual     string = "neq"
	Greater      string = "gt"
	GreaterEqual string = "gte"
	Less         string = "lt"
	LessEqual    string = "lte"
	Contains     string = "ilike"
	In           string = "in" // value is a '|' separated list: id,in,<a>|<b>|<c>

	Ascendent  string = "ASC"
	Descendent string = "DESC"
)
