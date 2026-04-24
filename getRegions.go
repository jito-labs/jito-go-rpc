package jitorpc

type Region string

var (
	regions = []Region{
		"Mainnet",
		"Amsterdam",
		"Dublin",
		"Frankfurt",
		"London",
		"New York",
		"Salt Lake City",
		"Singapore",
		"Tokyo",
	}
)

func GetRegions() []Region {
	return regions
}
