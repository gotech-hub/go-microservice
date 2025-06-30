package mongodb

import (
	"os"
	"strings"
)

var (
	mappingRepositoryRegion = map[string]map[string][]string{
		"DEFAULT": {
			"config_games":    {"VN::loyalty"},
			"config_currency": {"VN::loyalty"},
		},
	}
	countryMappingRegion = map[string]string{
		"VN":  "VN",
		"TW":  "SEA",
		"HK":  "SEA",
		"SG":  "SEA",
		"MY":  "SEA",
		"ID":  "SEA",
		"TH":  "SEA",
		"PH":  "SEA",
		"SEA": "SEA",
	}
)

func GetMappingRepositoryRegion(collectionName string) []string {
	env := strings.ToUpper(os.Getenv("ENV"))
	v, ok := mappingRepositoryRegion[env][collectionName]
	if !ok {
		return mappingRepositoryRegion["DEFAULT"][collectionName]
	}
	return v
}

func GetRegionCountry(country string) string {
	return countryMappingRegion[country]
}
