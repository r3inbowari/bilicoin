package bilicoin

// @author lycblank
// @change r3inbowari

import "fmt"

type Options struct {
	graph string
}

type Option func(opts *Options)

type ProgressBar struct {
	totalValue int64
	currValue  int64
	graph      string
	rate       string
}

func NewProgressBar(totalValue int64, options ...Option) *ProgressBar {
	opts := Options{}
	for _, opt := range options {
		opt(&opts)
	}
	if opts.graph == "" {
		opts.graph = "█"
	}
	bar := &ProgressBar{
		totalValue: totalValue,
		graph:      opts.graph,
	}
	return bar
}

func (bar *ProgressBar) Play(value int64) {
	val := float64(bar.totalValue) / 50
	prePercent := int32(float64(bar.currValue) / val)
	nowPercent := int32(float64(value) / val)
	for i := prePercent + 1; i <= nowPercent; i++ {
		bar.rate += bar.graph
	}
	bar.currValue = value
	fmt.Printf("\r[INFO] [%-50s]%0.2f%%   	%8d/%d", bar.rate, float64(bar.currValue)/float64(bar.totalValue)*100,
		bar.currValue, bar.totalValue)
}

func (bar *ProgressBar) Finish() {
	val := float64(bar.totalValue) / 50
	prePercent := int32(float64(bar.currValue) / val)
	for i := prePercent + 1; i <= 50; i++ {
		bar.rate += bar.graph
	}
	bar.currValue = bar.totalValue
	fmt.Printf("\r[INFO] [%-50s]%0.2f%%   	%8d/%d", bar.rate, float64(bar.currValue)/float64(bar.totalValue)*100,
		bar.currValue, bar.totalValue)
	fmt.Println()
}

func (bar *ProgressBar) Stop() {
	fmt.Println()
}
