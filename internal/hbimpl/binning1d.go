// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbimpl

import (
	"errors"
	"sort"
)

// Indices for the under- and over-flow 1-dim bins.
const (
	UnderflowBin = -1
	OverflowBin  = -2
)

var (
	errInvalidXAxis   = errors.New("hbook: invalid X-axis limits")
	errEmptyXAxis     = errors.New("hbook: X-axis with zero bins")
	errShortXAxis     = errors.New("hbook: too few 1-dim X-bins")
	errOverlapXAxis   = errors.New("hbook: invalid X-binning (overlap)")
	errNotSortedXAxis = errors.New("hbook: X-edges slice not sorted")
	errDupEdgesXAxis  = errors.New("hbook: duplicates in X-edge values")

	errInvalidYAxis   = errors.New("hbook: invalid Y-axis limits")
	errEmptyYAxis     = errors.New("hbook: Y-axis with zero bins")
	errShortYAxis     = errors.New("hbook: too few 1-dim Y-bins")
	errOverlapYAxis   = errors.New("hbook: invalid Y-binning (overlap)")
	errNotSortedYAxis = errors.New("hbook: Y-edges slice not sorted")
	errDupEdgesYAxis  = errors.New("hbook: duplicates in Y-edge values")
)

// Binning1D is a 1-dim binning of the x-axis.
type Binning1D struct {
	XXX_bins     []Bin1D
	XXX_dist     Dist1D
	XXX_outflows [2]Dist1D
	XXX_xrange   Range
}

func newBinning1D(n int, xmin, xmax float64) Binning1D {
	if xmin >= xmax {
		panic(errInvalidXAxis)
	}
	if n <= 0 {
		panic(errEmptyXAxis)
	}
	bng := Binning1D{
		XXX_bins:   make([]Bin1D, n),
		XXX_xrange: Range{Min: xmin, Max: xmax},
	}
	width := bng.XXX_xrange.Width() / float64(n)
	for i := range bng.XXX_bins {
		bin := &bng.XXX_bins[i]
		bin.XXX_xrange.Min = xmin + float64(i)*width
		bin.XXX_xrange.Max = xmin + float64(i+1)*width
	}
	return bng
}

func newBinning1DFromBins(xbins []Range) Binning1D {
	if len(xbins) < 1 {
		panic(errShortXAxis)
	}
	n := len(xbins)
	bng := Binning1D{
		XXX_bins: make([]Bin1D, n),
	}
	for i, xbin := range xbins {
		bin := &bng.XXX_bins[i]
		bin.XXX_xrange = xbin
	}
	sort.Sort(Bin1Ds(bng.XXX_bins))
	for i := 0; i < len(bng.XXX_bins)-1; i++ {
		b0 := bng.XXX_bins[i]
		b1 := bng.XXX_bins[i+1]
		if b0.XXX_xrange.Max > b1.XXX_xrange.Min {
			panic(errOverlapXAxis)
		}
	}
	bng.XXX_xrange = Range{Min: bng.XXX_bins[0].XMin(), Max: bng.XXX_bins[n-1].XMax()}
	return bng
}

func newBinning1DFromEdges(edges []float64) Binning1D {
	if len(edges) <= 1 {
		panic(errShortXAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(edges)) {
		panic(errNotSortedXAxis)
	}
	n := len(edges) - 1
	bng := Binning1D{
		XXX_bins:   make([]Bin1D, n),
		XXX_xrange: Range{Min: edges[0], Max: edges[n]},
	}
	for i := range bng.XXX_bins {
		bin := &bng.XXX_bins[i]
		xmin := edges[i]
		xmax := edges[i+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bin.XXX_xrange.Min = xmin
		bin.XXX_xrange.Max = xmax
	}
	return bng
}

func (bng *Binning1D) entries() int64 {
	return bng.XXX_dist.Entries()
}

func (bng *Binning1D) effEntries() float64 {
	return bng.XXX_dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *Binning1D) xMin() float64 {
	return bng.XXX_xrange.Min
}

// xMax returns the high edge of the X-axis
func (bng *Binning1D) xMax() float64 {
	return bng.XXX_xrange.Max
}

func (bng *Binning1D) fill(x, w float64) {
	idx := bng.coordToIndex(x)
	bng.XXX_dist.fill(x, w)
	if idx < 0 {
		bng.XXX_outflows[-idx-1].fill(x, w)
		return
	}
	if idx == len(bng.XXX_bins) {
		// gap bin.
		return
	}
	bng.XXX_bins[idx].fill(x, w)
}

// coordToIndex returns the bin index corresponding to the coordinate x.
func (bng *Binning1D) coordToIndex(x float64) int {
	switch {
	case x < bng.XXX_xrange.Min:
		return UnderflowBin
	case x >= bng.XXX_xrange.Max:
		return OverflowBin
	}
	return Bin1Ds(bng.XXX_bins).IndexOf(x)
}

func (bng *Binning1D) scaleW(f float64) {
	bng.XXX_dist.scaleW(f)
	bng.XXX_outflows[0].scaleW(f)
	bng.XXX_outflows[1].scaleW(f)
	for i := range bng.XXX_bins {
		bin := &bng.XXX_bins[i]
		bin.scaleW(f)
	}
}

// Bins returns the slice of bins for this binning.
func (bng *Binning1D) Bins() []Bin1D {
	return bng.XXX_bins
}

func (bng *Binning1D) Underflow() *Dist1D {
	return &bng.XXX_outflows[0]
}

func (bng *Binning1D) Overflow() *Dist1D {
	return &bng.XXX_outflows[1]
}
