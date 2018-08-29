// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "math"

// Dist0D is a 0-dim distribution.
type Dist0D struct {
	N     int64   // number of entries
	SumW  float64 // sum of weights
	sumW2 float64 // sum of squared weights
}

// Rank returns the number of dimensions of the distribution.
func (*Dist0D) Rank() int {
	return 1
}

// Entries returns the number of entries in the distribution.
func (d *Dist0D) Entries() int64 {
	return d.N
}

// EffEntries returns the number of weighted entries, such as:
//  (\sum w)^2 / \sum w^2
func (d *Dist0D) EffEntries() float64 {
	if d.sumW2 == 0 {
		return 0
	}
	return d.SumW * d.SumW / d.sumW2
}

// XXX_SumW returns the sum of weights of the distribution.
func (d *Dist0D) XXX_SumW() float64 {
	return d.SumW
}

// XXX_SumW2 returns the sum of squared weights of the distribution.
func (d *Dist0D) XXX_SumW2() float64 {
	return d.sumW2
}

// errW returns the absolute error on sumW()
func (d *Dist0D) errW() float64 {
	return math.Sqrt(d.XXX_SumW2())
}

// relErrW returns the relative error on sumW()
func (d *Dist0D) relErrW() float64 {
	// FIXME(sbinet) check for low stats ?
	return d.errW() / d.XXX_SumW()
}

func (d *Dist0D) fill(w float64) {
	d.N++
	d.SumW += w
	d.sumW2 += w * w
}

func (d *Dist0D) scaleW(f float64) {
	d.SumW *= f
	d.sumW2 *= f * f
}

// dist1D is a 1-dim distribution.
type dist1D struct {
	dist   Dist0D  // weight moments
	sumWX  float64 // 1st order weighted x moment
	sumWX2 float64 // 2nd order weighted x moment
}

// Rank returns the number of dimensions of the distribution.
func (*dist1D) Rank() int {
	return 1
}

// Entries returns the number of entries in the distribution.
func (d *dist1D) Entries() int64 {
	return d.dist.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *dist1D) EffEntries() float64 {
	return d.dist.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *dist1D) SumW() float64 {
	return d.dist.XXX_SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist1D) SumW2() float64 {
	return d.dist.XXX_SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *dist1D) SumWX() float64 {
	return d.sumWX
}

// SumWX2 returns the 2nd order weighted x moment
func (d *dist1D) SumWX2() float64 {
	return d.sumWX2
}

// errW returns the absolute error on sumW()
func (d *dist1D) errW() float64 {
	return d.dist.errW()
}

// relErrW returns the relative error on sumW()
func (d *dist1D) relErrW() float64 {
	return d.dist.relErrW()
}

// mean returns the weighted mean of the distribution
func (d *dist1D) mean() float64 {
	// FIXME(sbinet): check for low stats?
	return d.sumWX / d.SumW()
}

// variance returns the weighted variance of the distribution, defined as:
//  sig2 = ( \sum(wx^2) * \sum(w) - \sum(wx)^2 ) / ( \sum(w)^2 - \sum(w^2) )
// see: https://en.wikipedia.org/wiki/Weighted_arithmetic_mean
func (d *dist1D) variance() float64 {
	// FIXME(sbinet): check for low stats?
	sumw := d.SumW()
	num := d.sumWX2*sumw - d.sumWX*d.sumWX
	den := sumw*sumw - d.SumW2()
	v := num / den
	return math.Abs(v)
}

// stdDev returns the weighted standard deviation of the distribution
func (d *dist1D) stdDev() float64 {
	return math.Sqrt(d.variance())
}

// stdErr returns the weighted standard error of the distribution
func (d *dist1D) stdErr() float64 {
	// FIXME(sbinet): check for low stats?
	// TODO(sbinet): unbiased should check that Neff>1 and divide by N-1?
	return math.Sqrt(d.variance() / d.EffEntries())
}

// rms returns the weighted RMS of the distribution, defined as:
//  rms = \sqrt{\sum{w . x^2} / \sum{w}}
func (d *dist1D) rms() float64 {
	// FIXME(sbinet): check for low stats?
	meansq := d.sumWX2 / d.SumW()
	return math.Sqrt(meansq)
}

func (d *dist1D) fill(x, w float64) {
	d.dist.fill(w)
	d.sumWX += w * x
	d.sumWX2 += w * x * x
}

func (d *dist1D) scaleW(f float64) {
	d.dist.scaleW(f)
	d.sumWX *= f
	d.sumWX2 *= f
}

func (d *dist1D) scaleX(f float64) {
	d.sumWX *= f
	d.sumWX2 *= f * f
}

// Dist2D is a 2-dim distribution.
type Dist2D struct {
	X      dist1D  // x moments
	Y      dist1D  // y moments
	SumWXY float64 // 2nd-order cross-term
}

// Rank returns the number of dimensions of the distribution.
func (*Dist2D) Rank() int {
	return 2
}

// Entries returns the number of entries in the distribution.
func (d *Dist2D) Entries() int64 {
	return d.X.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *Dist2D) EffEntries() float64 {
	return d.X.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *Dist2D) SumW() float64 {
	return d.X.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *Dist2D) SumW2() float64 {
	return d.X.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *Dist2D) SumWX() float64 {
	return d.X.SumWX()
}

// SumWX2 returns the 2nd order weighted x moment
func (d *Dist2D) SumWX2() float64 {
	return d.X.SumWX2()
}

// SumWY returns the 1st order weighted y moment
func (d *Dist2D) SumWY() float64 {
	return d.Y.SumWX()
}

// SumWY2 returns the 2nd order weighted y moment
func (d *Dist2D) SumWY2() float64 {
	return d.Y.SumWX2()
}

// errW returns the absolute error on sumW()
func (d *Dist2D) errW() float64 {
	return d.X.errW()
}

// relErrW returns the relative error on sumW()
func (d *Dist2D) relErrW() float64 {
	return d.X.relErrW()
}

// xMean returns the weighted mean of the distribution
func (d *Dist2D) xMean() float64 {
	return d.X.mean()
}

// yMean returns the weighted mean of the distribution
func (d *Dist2D) yMean() float64 {
	return d.Y.mean()
}

// xVariance returns the weighted variance of the distribution
func (d *Dist2D) xVariance() float64 {
	return d.X.variance()
}

// yVariance returns the weighted variance of the distribution
func (d *Dist2D) yVariance() float64 {
	return d.Y.variance()
}

// xStdDev returns the weighted standard deviation of the distribution
func (d *Dist2D) xStdDev() float64 {
	return d.X.stdDev()
}

// yStdDev returns the weighted standard deviation of the distribution
func (d *Dist2D) yStdDev() float64 {
	return d.Y.stdDev()
}

// xStdErr returns the weighted standard error of the distribution
func (d *Dist2D) xStdErr() float64 {
	return d.X.stdErr()
}

// yStdErr returns the weighted standard error of the distribution
func (d *Dist2D) yStdErr() float64 {
	return d.Y.stdErr()
}

// xRMS returns the weighted RMS of the distribution
func (d *Dist2D) xRMS() float64 {
	return d.X.rms()
}

// yRMS returns the weighted RMS of the distribution
func (d *Dist2D) yRMS() float64 {
	return d.Y.rms()
}

func (d *Dist2D) fill(x, y, w float64) {
	d.X.fill(x, w)
	d.Y.fill(y, w)
	d.SumWXY += w * x * y
}

func (d *Dist2D) scaleW(f float64) {
	d.X.scaleW(f)
	d.Y.scaleW(f)
	d.SumWXY *= f
}

func (d *Dist2D) scaleX(f float64) {
	d.X.scaleX(f)
	d.SumWXY *= f
}

func (d *Dist2D) scaleY(f float64) {
	d.Y.scaleX(f)
	d.SumWXY *= f
}

func (d *Dist2D) scaleXY(fx, fy float64) {
	d.scaleX(fx)
	d.scaleY(fy)
}
