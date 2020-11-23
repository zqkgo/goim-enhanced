package exporter

type Exporter interface {
	Export() (interface{}, error)
	Options() Options
	String() string
}
