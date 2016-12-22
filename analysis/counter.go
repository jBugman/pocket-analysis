package analysis

import "sort"

type Counter map[string]int

type CounterItem struct {
	Key   string
	Count int
}

type Weights map[string]float64

func (this Counter) Add(key string) {
	this[key] = this[key] + 1
}

func (this Counter) Items() (result CounterItems) {
	return this.ItemsWithThreshold(0)
}

func (this Counter) ItemsWithThreshold(threshold int) (result CounterItems) {
	for k, v := range this {
		if v >= threshold {
			result = append(result, CounterItem{Key: k, Count: v})
		}
	}
	sort.Sort(sort.Reverse(result))
	return
}

func (this Counter) Weights(threshold float64) (result Weights) {
	sum := 0
	for _, v := range this {
		sum = sum + v
	}
	sumF := float64(sum)
	result = Weights{}
	for k, v := range this {
		weight := float64(v) / sumF
		if (threshold < 0 && weight <= -threshold) || (threshold >= 0 && weight >= threshold) {
			result[k] = weight
		}
	}
	return
}

/* sort.Interface implementation */
type CounterItems []CounterItem

func (a CounterItems) Len() int           { return len(a) }
func (a CounterItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CounterItems) Less(i, j int) bool { return a[i].Count < a[j].Count }
