package schema

type Tier struct {
	Level    Level
	Interval [2]int
	Val      [2]int
}

type Level struct {
	Name  string
	Color string
}
