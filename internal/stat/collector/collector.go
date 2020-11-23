package collector

type Collector interface {
	Init(opts ...Option) error
	Collect() error
	Stop()
	ReCollect(opts ...Option) error
	Options() *Options
	String() string
}

func FireCollectors(collectors ...Collector) {
	for _, c := range collectors {
		go c.Collect()
	}
}
