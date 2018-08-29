// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"

	"go-hep.org/x/hep/hbook"
)

func newH1D() *H1D {
	return &H1D{
		rvers: 2, // FIXME(sbinet): harmonize versions
		th1:   *newH1(),
	}
}

func newH1() *th1 {
	return &th1{
		rvers:     7, // FIXME(sbinet): harmonize versions
		tnamed:    *newNamed("", ""),
		attline:   *newAttLine(),
		attfill:   *newAttFill(),
		attmarker: *newAttMarker(),
		xaxis:     *newAxis("xaxis"),
		yaxis:     *newAxis("yaxis"),
		zaxis:     *newAxis("zaxis"),
		funcs:     *newList(""),
	}
}

func NewH1DFrom(h *hbook.H1D) *H1D {
	var (
		hroot = newH1D()
		nbins = h.Len()
		edges = make([]float64, 0, nbins+1)
		sumw  = make([]float64, 0, nbins+2)
		sumw2 = make([]float64, 0, nbins+2)
		uflow = h.Binning().Underflow()
		oflow = h.Binning().Overflow()
	)

	sumw = append(sumw, uflow.SumW())
	sumw2 = append(sumw2, uflow.SumW2())

	for i, bin := range h.Binning().Bins() {
		if i == 0 {
			edges = append(edges, bin.XMin())
		}
		edges = append(edges, bin.XMax())
		sumw = append(sumw, bin.SumW())
		sumw2 = append(sumw2, bin.SumW2())
	}
	sumw = append(sumw, oflow.SumW())
	sumw2 = append(sumw2, oflow.SumW2())

	hroot.th1.name = h.Name()
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.title = v.(string)
	}
	hroot.th1.entries = float64(h.Entries())
	hroot.th1.tsumw = h.SumW()
	hroot.th1.tsumw2 = h.SumW2()
	hroot.th1.tsumwx = h.SumWX()
	hroot.th1.tsumwx2 = h.SumWX2()
	hroot.th1.ncells = len(edges)
	hroot.th1.xaxis.xbins.Data = edges
	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()
	hroot.arr.Data = sumw
	hroot.sumw2.Data = sumw2

	return hroot
}

func (h1d *H1D) UnmarshalYODA(raw []byte) error {
	var h hbook.H1D
	err := h.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h1d = *NewH1DFrom(&h)
	return nil
}

func newH2D() *H2D {
	return &H2D{
		rvers: 3, // FIXME(sbinet): harmonize versions
		th2:   *newH2(),
	}
}

func newH2() *th2 {
	return &th2{
		rvers: 4, // FIXME(sbinet): harmonize versions
		th1:   *newH1(),
	}
}

func NewH2DFrom(h *hbook.H2D) *H2D {
	var (
		hroot  = newH2D()
		bins   = h.Binning().XXX_GetBins()
		xedges = make([]float64, 0, h.Bng.Nx+1)
		yedges = make([]float64, 0, h.Bng.Ny+1)
		sumw   = make([]float64, 0)
		sumw2  = make([]float64, 0)
	)

	hroot.th2.th1.entries = float64(h.Entries())
	hroot.th2.th1.tsumw = h.SumW()
	hroot.th2.th1.tsumw2 = h.SumW2()
	hroot.th2.th1.tsumwx = h.SumWX()
	hroot.th2.th1.tsumwx2 = h.SumWX2()
	hroot.th2.tsumwy = h.SumWY()
	hroot.th2.tsumwy2 = h.SumWY2()
	hroot.th2.tsumwxy = h.SumWXY()

	hroot.th2.th1.ncells = len(bins)

	hroot.th2.th1.xaxis.nbins = h.Bng.Nx
	hroot.th2.th1.xaxis.xmin = h.XMin()
	hroot.th2.th1.xaxis.xmax = h.XMax()

	hroot.th2.th1.yaxis.nbins = h.Bng.Ny
	hroot.th2.th1.yaxis.xmin = h.YMin()
	hroot.th2.th1.yaxis.xmax = h.YMax()

	ncells := (h.Bng.Nx + 2) * (h.Bng.Ny + 2)
	hroot.arr.Data = make([]float64, ncells)
	hroot.th2.sumw2.Data = make([]float64, ncells)
	hroot.th1.xaxis.xbins.Data = make([]float64, h.Bng.Nx+1)
	hroot.th1.yaxis.xbins.Data = make([]float64, h.Bng.Ny+1)

	for ix := 0; ix < h.Bng.Nx; ix++ {
		for iy := 0; iy < h.Bng.Ny; iy++ {
			i := iy*h.Bng.Nx + ix
			bin := bins[i]
			if ix == 0 {
				xedges = append(xedges, bin.XMin())
				yedges = append(yedges, bin.YMin())
			}
			hroot.setDist2D(ix+1, iy+1, bin)
		}
	}

	hroot.th2.th1.name = h.Name()
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th2.th1.title = v.(string)
	}
	hroot.th2.th1.xaxis.xbins.Data = xedges
	hroot.th2.th1.yaxis.xbins.Data = yedges
	hroot.arr.Data = sumw
	hroot.sumw2.Data = sumw2

	return hroot
}

func (h *H2D) setDist2D(ix, iy int, v hbook.Bin2D) {
	i := h.bin(ix, iy)
	if i >= len(h.arr.Data) {
		panic(fmt.Errorf("err: i=%d, len=%d %d %d", i, len(h.arr.Data), ix, iy))
	}
	h.arr.Data[i] = v.Dist.SumW()
	h.th1.sumw2.Data[i] = v.Dist.SumW2()
	h.th1.xaxis.xbins.Data[ix] = v.XRange.Max
	h.th1.yaxis.xbins.Data[iy] = v.YRange.Max
}
