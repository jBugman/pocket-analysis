package analysis

import "sort"

type Counter map[string]int

type CounterItem struct {
	Key   string
	Count int
}

func (this Counter) Add(key string) {
	this[key] = this[key] + 1
}

func (this Counter) Items() (result CounterItems) {
	for k, v := range this {
		result = append(result, CounterItem{Key: k, Count: v})
	}
	sort.Sort(sort.Reverse(result))
	return
}

/* sort.Interface implementation */
type CounterItems []CounterItem

func (a CounterItems) Len() int           { return len(a) }
func (a CounterItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CounterItems) Less(i, j int) bool { return a[i].Count < a[j].Count }
