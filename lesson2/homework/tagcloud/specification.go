package tagcloud

import (
	"sort"
)

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	dict map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{map[string]int{}}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (selfPointer *TagCloud) AddTag(tag string) {
	selfPointer.dict[tag] += 1
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (selfPointer *TagCloud) TopN(n int) []TagStat {
	var result []TagStat
	for key, value := range selfPointer.dict {
		result = append(result, TagStat{
			Tag:             key,
			OccurrenceCount: value,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].OccurrenceCount > result[j].OccurrenceCount
	})
	if n >= len(selfPointer.dict) && len(selfPointer.dict) > 0 {
		return result
	} else if len(selfPointer.dict) > 0 {
		return result[:n]
	} else {
		return nil
	}
}
