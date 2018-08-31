// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbimpl

import "sort"

// indices for the 2D-binning overflows
const (
	bngNW int = 1 + iota
	bngN
	bngNE
	bngE
	bngSE
	bngS
	bngSW
	bngW
)

type Binning2D struct {
	XXX_bins     []Bin2D
	XXX_dist     dist2D
	XXX_outflows [8]dist2D
	XXX_xrange   Range
	XXX_yrange   Range
	XXX_nx       int
	XXX_ny       int
	XXX_xedges   []Bin1D
	XXX_yedges   []Bin1D
}

func newBinning2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) Binning2D {
	if xlow >= xhigh {
		panic(errInvalidXAxis)
	}
	if ylow >= yhigh {
		panic(errInvalidYAxis)
	}
	if nx <= 0 {
		panic(errEmptyXAxis)
	}
	if ny <= 0 {
		panic(errEmptyYAxis)
	}
	bng := Binning2D{
		XXX_bins:   make([]Bin2D, nx*ny),
		XXX_xrange: Range{Min: xlow, Max: xhigh},
		XXX_yrange: Range{Min: ylow, Max: yhigh},
		XXX_nx:     nx,
		XXX_ny:     ny,
		XXX_xedges: make([]Bin1D, nx),
		XXX_yedges: make([]Bin1D, ny),
	}
	xwidth := bng.XXX_xrange.Width() / float64(bng.XXX_nx)
	ywidth := bng.XXX_yrange.Width() / float64(bng.XXX_ny)
	xmin := bng.XXX_xrange.Min
	ymin := bng.XXX_yrange.Min
	for ix := range bng.XXX_xedges {
		xbin := &bng.XXX_xedges[ix]
		xbin.XXX_xrange.Min = xmin + float64(ix)*xwidth
		xbin.XXX_xrange.Max = xmin + float64(ix+1)*xwidth
		for iy := range bng.XXX_yedges {
			ybin := &bng.XXX_yedges[iy]
			ybin.XXX_xrange.Min = ymin + float64(iy)*ywidth
			ybin.XXX_xrange.Max = ymin + float64(iy+1)*ywidth
			i := iy*nx + ix
			bin := &bng.XXX_bins[i]
			bin.XXX_xrange.Min = xbin.XXX_xrange.Min
			bin.XXX_xrange.Max = xbin.XXX_xrange.Max
			bin.XXX_yrange.Min = ybin.XXX_xrange.Min
			bin.XXX_yrange.Max = ybin.XXX_xrange.Max
		}
	}
	return bng
}

func newBinning2DFromEdges(xedges, yedges []float64) Binning2D {
	if len(xedges) <= 1 {
		panic(errShortXAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(xedges)) {
		panic(errNotSortedXAxis)
	}
	if len(yedges) <= 1 {
		panic(errShortYAxis)
	}
	if !sort.IsSorted(sort.Float64Slice(yedges)) {
		panic(errNotSortedYAxis)
	}
	var (
		nx    = len(xedges) - 1
		ny    = len(yedges) - 1
		xlow  = xedges[0]
		xhigh = xedges[nx]
		ylow  = yedges[0]
		yhigh = yedges[ny]
	)
	bng := Binning2D{
		XXX_bins:   make([]Bin2D, nx*ny),
		XXX_xrange: Range{Min: xlow, Max: xhigh},
		XXX_yrange: Range{Min: ylow, Max: yhigh},
		XXX_nx:     nx,
		XXX_ny:     ny,
		XXX_xedges: make([]Bin1D, nx),
		XXX_yedges: make([]Bin1D, ny),
	}
	for ix, xmin := range xedges[:nx] {
		xmax := xedges[ix+1]
		if xmin == xmax {
			panic(errDupEdgesXAxis)
		}
		bng.XXX_xedges[ix].XXX_xrange.Min = xmin
		bng.XXX_xedges[ix].XXX_xrange.Max = xmax
		for iy, ymin := range yedges[:ny] {
			ymax := yedges[iy+1]
			if ymin == ymax {
				panic(errDupEdgesYAxis)
			}
			i := iy*nx + ix
			bin := &bng.XXX_bins[i]
			bin.XXX_xrange.Min = xmin
			bin.XXX_xrange.Max = xmax
			bin.XXX_yrange.Min = ymin
			bin.XXX_yrange.Max = ymax
		}
	}
	for iy, ymin := range yedges[:ny] {
		ymax := yedges[iy+1]
		bng.XXX_yedges[iy].XXX_xrange.Min = ymin
		bng.XXX_yedges[iy].XXX_xrange.Max = ymax
	}
	return bng
}

func (bng *Binning2D) entries() int64 {
	return bng.XXX_dist.Entries()
}

func (bng *Binning2D) effEntries() float64 {
	return bng.XXX_dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *Binning2D) xMin() float64 {
	return bng.XXX_xrange.Min
}

// xMax returns the high edge of the X-axis
func (bng *Binning2D) xMax() float64 {
	return bng.XXX_xrange.Max
}

// yMin returns the low edge of the Y-axis
func (bng *Binning2D) yMin() float64 {
	return bng.XXX_yrange.Min
}

// yMax returns the high edge of the Y-axis
func (bng *Binning2D) yMax() float64 {
	return bng.XXX_yrange.Max
}

func (bng *Binning2D) fill(x, y, w float64) {
	idx := bng.coordToIndex(x, y)
	bng.XXX_dist.fill(x, y, w)
	if idx == len(bng.XXX_bins) {
		// GAP bin
		return
	}
	if idx < 0 {
		bng.XXX_outflows[-idx-1].fill(x, y, w)
		return
	}
	bng.XXX_bins[idx].fill(x, y, w)
}

func (bng *Binning2D) coordToIndex(x, y float64) int {
	ix := Bin1Ds(bng.XXX_xedges).IndexOf(x)
	iy := Bin1Ds(bng.XXX_yedges).IndexOf(y)

	switch {
	case ix == bng.XXX_nx && iy == bng.XXX_ny: // GAP
		return len(bng.XXX_bins)
	case ix == OverflowBin && iy == OverflowBin:
		return -bngNE
	case ix == OverflowBin && iy == UnderflowBin:
		return -bngSE
	case ix == UnderflowBin && iy == UnderflowBin:
		return -bngSW
	case ix == UnderflowBin && iy == OverflowBin:
		return -bngNW
	case ix == OverflowBin:
		return -bngE
	case ix == UnderflowBin:
		return -bngW
	case iy == OverflowBin:
		return -bngN
	case iy == UnderflowBin:
		return -bngS
	}
	return iy*bng.XXX_nx + ix
}

// Bins returns the slice of bins for this binning.
func (bng *Binning2D) Bins() []Bin2D {
	return bng.XXX_bins
}
