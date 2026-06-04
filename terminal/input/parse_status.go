package input

type parseStatus uint8

const (
	parseNoMatch parseStatus = iota
	parseNeedMore
	parseDone
)
