// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbimpl

// Bin2D models a bin in a 2-dim space.
type Bin2D struct {
	XXX_xrange Range
	XXX_yrange Range
	XXX_dist   dist2D
}

// Rank returns the number of dimensions for this bin.
func (Bin2D) Rank() int { return 2 }

func (b *Bin2D) scaleW(f float64) {
	b.XXX_dist.scaleW(f)
}

func (b *Bin2D) fill(x, y, w float64) {
	b.XXX_dist.fill(x, y, w)
}

// Entries returns the number of entries in this bin.
func (b *Bin2D) Entries() int64 {
	return b.XXX_dist.Entries()
}

// EffEntries returns the effective number of entries \f$ = (\sum w)^2 / \sum w^2 \f$
func (b *Bin2D) EffEntries() float64 {
	return b.XXX_dist.EffEntries()
}

// SumW returns the sum of weights in this bin.
func (b *Bin2D) SumW() float64 {
	return b.XXX_dist.SumW()
}

// SumW2 returns the sum of squared weights in this bin.
func (b *Bin2D) SumW2() float64 {
	return b.XXX_dist.SumW2()
}

// XEdges returns the [low,high] edges of this bin.
func (b *Bin2D) XEdges() Range {
	return b.XXX_xrange
}

// YEdges returns the [low,high] edges of this bin.
func (b *Bin2D) YEdges() Range {
	return b.XXX_yrange
}

// XMin returns the lower limit of the bin (inclusive).
func (b *Bin2D) XMin() float64 {
	return b.XXX_xrange.Min
}

// YMin returns the lower limit of the bin (inclusive).
func (b *Bin2D) YMin() float64 {
	return b.XXX_yrange.Min
}

// XMax returns the upper limit of the bin (exclusive).
func (b *Bin2D) XMax() float64 {
	return b.XXX_xrange.Max
}

// YMax returns the upper limit of the bin (exclusive).
func (b *Bin2D) YMax() float64 {
	return b.XXX_yrange.Max
}

// XMid returns the geometric center of the bin.
// i.e.: 0.5*(high+low)
func (b *Bin2D) XMid() float64 {
	return 0.5 * (b.XXX_xrange.Min + b.XXX_xrange.Max)
}

// YMid returns the geometric center of the bin.
// i.e.: 0.5*(high+low)
func (b *Bin2D) YMid() float64 {
	return 0.5 * (b.XXX_yrange.Min + b.XXX_yrange.Max)
}

// XYMid returns the (x,y) coordinates of the geometric center of the bin.
// i.e.: 0.5*(high+low)
func (b *Bin2D) XYMid() (float64, float64) {
	return b.XMid(), b.YMid()
}

// XWidth returns the (signed) width of the bin
func (b *Bin2D) XWidth() float64 {
	return b.XXX_xrange.Max - b.XXX_xrange.Min
}

// YWidth returns the (signed) width of the bin
func (b *Bin2D) YWidth() float64 {
	return b.XXX_yrange.Max - b.XXX_yrange.Min
}

// XYWidth returns the (signed) (x,y) widths of the bin
func (b *Bin2D) XYWidth() (float64, float64) {
	return b.XWidth(), b.YWidth()
}

// XFocus returns the mean position in the bin, or the midpoint (if the
// sum of weights for this bin is 0).
func (b *Bin2D) XFocus() float64 {
	if b.SumW() == 0 {
		return b.XMid()
	}
	return b.XMean()
}

// YFocus returns the mean position in the bin, or the midpoint (if the
// sum of weights for this bin is 0).
func (b *Bin2D) YFocus() float64 {
	if b.SumW() == 0 {
		return b.YMid()
	}
	return b.YMean()
}

// XYFocus returns the mean position in the bin, or the midpoint (if the
// sum of weights for this bin is 0).
func (b *Bin2D) XYFocus() (float64, float64) {
	if b.SumW() == 0 {
		return b.XMid(), b.YMid()
	}
	return b.XMean(), b.YMean()
}

// XMean returns the mean X.
func (b *Bin2D) XMean() float64 {
	return b.XXX_dist.xMean()
}

// YMean returns the mean Y.
func (b *Bin2D) YMean() float64 {
	return b.XXX_dist.yMean()
}

// XVariance returns the variance in X.
func (b *Bin2D) XVariance() float64 {
	return b.XXX_dist.xVariance()
}

// YVariance returns the variance in Y.
func (b *Bin2D) YVariance() float64 {
	return b.XXX_dist.yVariance()
}

// XStdDev returns the standard deviation in X.
func (b *Bin2D) XStdDev() float64 {
	return b.XXX_dist.xStdDev()
}

// YStdDev returns the standard deviation in Y.
func (b *Bin2D) YStdDev() float64 {
	return b.XXX_dist.yStdDev()
}

// XStdErr returns the standard error in X.
func (b *Bin2D) XStdErr() float64 {
	return b.XXX_dist.xStdErr()
}

// YStdErr returns the standard error in Y.
func (b *Bin2D) YStdErr() float64 {
	return b.XXX_dist.yStdErr()
}

// XRMS returns the RMS in X.
func (b *Bin2D) XRMS() float64 {
	return b.XXX_dist.xRMS()
}

// YRMS returns the RMS in Y.
func (b *Bin2D) YRMS() float64 {
	return b.XXX_dist.yRMS()
}
