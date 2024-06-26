package controller

import (
	"github.com/kaytu-io/kaytu/pkg/plugin/proto/src/golang"
)

type Optimizations[T golang.OptimizationItem | golang.ChartOptimizationItem] struct {
	itemsChan chan *T
	items     []*T

	summaryChan chan string
	summary     string

	selectedItem *T

	reEvaluateFunc func(id string, items []*golang.PreferenceItem)
	initializing   bool
}

func NewOptimizations[T golang.OptimizationItem | golang.ChartOptimizationItem]() *Optimizations[T] {
	o := Optimizations[T]{
		itemsChan:    make(chan *T, 1000),
		initializing: true,
		summaryChan: make(chan string),
	}
	go o.Process()
	go o.SummaryProcess()
	return &o
}

func (o *Optimizations[T]) Process() {
	for newItem := range o.itemsChan {
		if o.initializing {
			o.initializing = false
		}
		updated := false
		for idx, i := range o.items {
			switch castedNewItem := any(newItem).(type) {
			case *golang.OptimizationItem:
				castedI := any(i).(*golang.OptimizationItem)
				if castedNewItem.Id == castedI.Id {
					o.items[idx] = newItem
					updated = true
					break
				}
			case *golang.ChartOptimizationItem:
				castedI := any(i).(*golang.ChartOptimizationItem)
				if castedNewItem.GetOverviewChartRow().GetRowId() == castedI.GetOverviewChartRow().GetRowId() {
					o.items[idx] = newItem
					updated = true
					break
				}
			}
		}
		if !updated {
			o.items = append(o.items, newItem)
		}
	}
}

func (o *Optimizations[T]) SummaryProcess() {
	for msg := range o.summaryChan {
		o.summary = msg
	}
}

func (o *Optimizations[T]) SendItem(item *T) {
	o.itemsChan <- item
}

func (o *Optimizations[T]) Items() []*T {
	return o.items
}

func (o *Optimizations[T]) SetReEvaluateFunc(f func(id string, items []*golang.PreferenceItem)) {
	o.reEvaluateFunc = f
}

func (o *Optimizations[T]) SelectItem(i *T) {
	o.selectedItem = i
}

func (o *Optimizations[T]) SelectedItem() *T {
	return o.selectedItem
}

func (o *Optimizations[T]) ReEvaluate(id string, preferences []*golang.PreferenceItem) {
	o.reEvaluateFunc(id, preferences)
}

func (o *Optimizations[T]) GetInitialization() bool {
	return o.initializing
}

func (o *Optimizations[T]) SetResultSummary(msg string) {
	o.summaryChan <- msg
}

func (o *Optimizations[T]) GetResultSummary() string {
	return o.summary
}
