// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbimpl

import "math"

// Dist0D is a 0-dim distribution.
type Dist0D struct {
	XXX_n     int64   // number of entries
	XXX_sumW  float64 // sum of weights
	XXX_sumW2 float64 // sum of squared weights
}

// Rank returns the number of dimensions of the distribution.
func (*Dist0D) Rank() int {
	return 1
}

// Entries returns the number of entries in the distribution.
func (d *Dist0D) Entries() int64 {
	return d.XXX_n
}

// EffEntries returns the number of weighted entries, such as:
//  (\sum w)^2 / \sum w^2
func (d *Dist0D) EffEntries() float64 {
	if d.XXX_sumW2 == 0 {
		return 0
	}
	return d.XXX_sumW * d.XXX_sumW / d.XXX_sumW2
}

// SumW returns the sum of weights of the distribution.
func (d *Dist0D) SumW() float64 {
	return d.XXX_sumW
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *Dist0D) SumW2() float64 {
	return d.XXX_sumW2
}

// errW returns the absolute error on sumW()
func (d *Dist0D) errW() float64 {
	return math.Sqrt(d.SumW2())
}

// relErrW returns the relative error on sumW()
func (d *Dist0D) relErrW() float64 {
	// FIXME(sbinet) check for low stats ?
	return d.errW() / d.SumW()
}

func (d *Dist0D) fill(w float64) {
	d.XXX_n++
	d.XXX_sumW += w
	d.XXX_sumW2 += w * w
}

func (d *Dist0D) scaleW(f float64) {
	d.XXX_sumW *= f
	d.XXX_sumW2 *= f * f
}

// Dist1D is a 1-dim distribution.
type Dist1D struct {
	dist   Dist0D  // weight moments
	sumWX  float64 // 1st order weighted x moment
	sumWX2 float64 // 2nd order weighted x moment
}

// Rank returns the number of dimensions of the distribution.
func (*Dist1D) Rank() int {
	return 1
}

// Entries returns the number of entries in the distribution.
func (d *Dist1D) Entries() int64 {
	return d.dist.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *Dist1D) EffEntries() float64 {
	return d.dist.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *Dist1D) SumW() float64 {
	return d.dist.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *Dist1D) SumW2() float64 {
	return d.dist.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *Dist1D) SumWX() float64 {
	return d.sumWX
}

// SumWX2 returns the 2nd order weighted x moment
func (d *Dist1D) SumWX2() float64 {
	return d.sumWX2
}

// errW returns the absolute error on sumW()
func (d *Dist1D) errW() float64 {
	return d.dist.errW()
}

// relErrW returns the relative error on sumW()
func (d *Dist1D) relErrW() float64 {
	return d.dist.relErrW()
}

// mean returns the weighted mean of the distribution
func (d *Dist1D) mean() float64 {
	// FIXME(sbinet): check for low stats?
	return d.sumWX / d.SumW()
}

// variance returns the weighted variance of the distribution, defined as:
//  sig2 = ( \sum(wx^2) * \sum(w) - \sum(wx)^2 ) / ( \sum(w)^2 - \sum(w^2) )
// see: https://en.wikipedia.org/wiki/Weighted_arithmetic_mean
func (d *Dist1D) variance() float64 {
	// FIXME(sbinet): check for low stats?
	sumw := d.SumW()
	num := d.sumWX2*sumw - d.sumWX*d.sumWX
	den := sumw*sumw - d.SumW2()
	v := num / den
	return math.Abs(v)
}

// stdDev returns the weighted standard deviation of the distribution
func (d *Dist1D) stdDev() float64 {
	return math.Sqrt(d.variance())
}

// stdErr returns the weighted standard error of the distribution
func (d *Dist1D) stdErr() float64 {
	// FIXME(sbinet): check for low stats?
	// TODO(sbinet): unbiased should check that Neff>1 and divide by N-1?
	return math.Sqrt(d.variance() / d.EffEntries())
}

// rms returns the weighted RMS of the distribution, defined as:
//  rms = \sqrt{\sum{w . x^2} / \sum{w}}
func (d *Dist1D) rms() float64 {
	// FIXME(sbinet): check for low stats?
	meansq := d.sumWX2 / d.SumW()
	return math.Sqrt(meansq)
}

func (d *Dist1D) fill(x, w float64) {
	d.dist.fill(w)
	d.sumWX += w * x
	d.sumWX2 += w * x * x
}

func (d *Dist1D) scaleW(f float64) {
	d.dist.scaleW(f)
	d.sumWX *= f
	d.sumWX2 *= f
}

func (d *Dist1D) scaleX(f float64) {
	d.sumWX *= f
	d.sumWX2 *= f * f
}

// dist2D is a 2-dim distribution.
type dist2D struct {
	x      Dist1D  // x moments
	y      Dist1D  // y moments
	sumWXY float64 // 2nd-order cross-term
}

// Rank returns the number of dimensions of the distribution.
func (*dist2D) Rank() int {
	return 2
}

// Entries returns the number of entries in the distribution.
func (d *dist2D) Entries() int64 {
	return d.x.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *dist2D) EffEntries() float64 {
	return d.x.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *dist2D) SumW() float64 {
	return d.x.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist2D) SumW2() float64 {
	return d.x.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *dist2D) SumWX() float64 {
	return d.x.SumWX()
}

// SumWX2 returns the 2nd order weighted x moment
func (d *dist2D) SumWX2() float64 {
	return d.x.SumWX2()
}

// SumWY returns the 1st order weighted y moment
func (d *dist2D) SumWY() float64 {
	return d.y.SumWX()
}

// SumWY2 returns the 2nd order weighted y moment
func (d *dist2D) SumWY2() float64 {
	return d.y.SumWX2()
}

// errW returns the absolute error on sumW()
func (d *dist2D) errW() float64 {
	return d.x.errW()
}

// relErrW returns the relative error on sumW()
func (d *dist2D) relErrW() float64 {
	return d.x.relErrW()
}

// xMean returns the weighted mean of the distribution
func (d *dist2D) xMean() float64 {
	return d.x.mean()
}

// yMean returns the weighted mean of the distribution
func (d *dist2D) yMean() float64 {
	return d.y.mean()
}

// xVariance returns the weighted variance of the distribution
func (d *dist2D) xVariance() float64 {
	return d.x.variance()
}

// yVariance returns the weighted variance of the distribution
func (d *dist2D) yVariance() float64 {
	return d.y.variance()
}

// xStdDev returns the weighted standard deviation of the distribution
func (d *dist2D) xStdDev() float64 {
	return d.x.stdDev()
}

// yStdDev returns the weighted standard deviation of the distribution
func (d *dist2D) yStdDev() float64 {
	return d.y.stdDev()
}

// xStdErr returns the weighted standard error of the distribution
func (d *dist2D) xStdErr() float64 {
	return d.x.stdErr()
}

// yStdErr returns the weighted standard error of the distribution
func (d *dist2D) yStdErr() float64 {
	return d.y.stdErr()
}

// xRMS returns the weighted RMS of the distribution
func (d *dist2D) xRMS() float64 {
	return d.x.rms()
}

// yRMS returns the weighted RMS of the distribution
func (d *dist2D) yRMS() float64 {
	return d.y.rms()
}

func (d *dist2D) fill(x, y, w float64) {
	d.x.fill(x, w)
	d.y.fill(y, w)
	d.sumWXY += w * x * y
}

func (d *dist2D) scaleW(f float64) {
	d.x.scaleW(f)
	d.y.scaleW(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleX(f float64) {
	d.x.scaleX(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleY(f float64) {
	d.y.scaleX(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleXY(fx, fy float64) {
	d.scaleX(fx)
	d.scaleY(fy)
}
