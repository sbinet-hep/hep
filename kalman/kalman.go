// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kalman

// biblio:
//  - https://www.uzh.ch/cmsssl/physik/dam/jcr:e705cebf-c99d-4651-9a78-761a9f96c66c/empp15_OS_reco.pdf
//  - http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.57.1034&rep=rep1&type=pdf

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)

var (
	errNotSquare = errors.New("kalman: non-square matrix")
	errMismatch  = errors.New("kalman: dimension mismatch")
)

func newIdentity(n int) *mat.Dense {
	if n <= 0 {
		panic("kalman: invalid matrix identity dimension")
	}
	data := make([]float64, n*n)
	for i := 0; i < n; i++ {
		data[i*(i+1)] = 1
	}
	return mat.NewDense(n, n, data)
}

type Settings struct {
	F *mat.Dense
	G *mat.Dense
	Q *mat.Dense

	H *mat.Dense
	R *mat.Dense
}

type KF struct {
	Settings

	id *mat.Dense // identity matrix

	x *mat.VecDense
	v *mat.Dense
}

func NewKF(settings *Settings) *KF {
	rf, cf := settings.F.Dims()
	rg, cg := settings.G.Dims()
	rq, cq := settings.Q.Dims()
	rh, ch := settings.H.Dims()
	rr, cr := settings.R.Dims()

	switch {
	case rf != cf, rq != cq, rr != cr:
		panic(errNotSquare)
	case rf != rg:
		panic(errMismatch)
	case cg != rq:
		panic(errMismatch)
	case ch != cf:
		panic(errMismatch)
	case rh != rr:
		panic(errMismatch)
	}

	return &KF{
		Settings: *settings,
		id:       newIdentity(rf),
		x:        mat.NewVecDense(cf, nil),
		v:        mat.NewDense(rf, cf, nil),
	}
}

func (kf *KF) Init(x *mat.VecDense, v mat.Matrix) {
	rf, cf := kf.Settings.F.Dims()
	switch x {
	case nil:
		kf.x = mat.NewVecDense(cf, nil)
	default:
		rx, _ := x.Dims()
		if rx != cf {
			panic(errMismatch)
		}
		kf.x = x
	}

	switch v {
	case nil:
		var m mat.Dense
		m.Mul(kf.G, kf.Q)
		kf.v.Mul(&m, kf.G.T())
	default:
		rv, cv := v.Dims()
		switch {
		case rv != cv:
			panic(errNotSquare)
		case rv != rf:
			panic(errMismatch)
		}
		kf.v = mat.DenseCopyOf(v)
	}
}

func (kf *KF) Filter(out, sys *mat.Dense) (*mat.Dense, error) {
	var (
		rh, _  = kf.H.Dims()
		rs, cs = sys.Dims()
	)

	if out == nil {
		out = mat.NewDense(rh, cs, nil)
	}

	var (
		m0, m1, m2     mat.Dense
		di0, di, k, kh mat.Dense
		ke, e          mat.VecDense
		retv           mat.VecDense
		arr            []float64

		ft = kf.F.T()
		gt = kf.G.T()
		ht = kf.H.T()
	)

	for j := 0; j < cs; j++ {
		// x = F . x
		kf.x.MulVec(kf.F, kf.x)

		// v = F . v . F^T + G . Q . G^T
		m0.Mul(kf.F, kf.v)
		kf.v.Mul(&m0, ft)

		m1.Mul(kf.G, kf.Q)
		m2.Mul(&m1, gt)
		kf.v.Add(kf.v, &m2)

		// d = (H . v . H^T + R)^{-1}
		di0.Mul(kf.H, kf.v)
		di.Mul(&di0, ht)
		di.Add(&di, kf.R)
		err := di.Inverse(&di) // FIXME(sbinet): not numerically stable.
		if err != nil {
			return nil, err
		}

		// v . H^T . d^-1
		k.Mul(kf.v, ht)
		k.Mul(&k, &di)

		// e = y - H . x
		y := sys.ColView(j)
		e.MulVec(kf.H, kf.x)
		e.SubVec(y, &e)

		hasNaN := false
		for i := 0; i < rs; i++ {
			v := y.At(i, 0)
			if math.IsNaN(v) {
				sys.Set(i, j, 0)
				hasNaN = true
			}
		}

		if !hasNaN {
			// x = x + K . e
			ke.MulVec(&k, &e)
			kf.x.AddVec(kf.x, &ke)

			// v = (I - K . H) . v
			kh.Mul(&k, kf.H)
			kh.Sub(kf.id, &kh)
			kf.v.Mul(&kh, kf.v)
		}

		retv.MulVec(kf.H, kf.x)
		arr = mat.Col(arr, 0, &retv)
		out.SetCol(j, arr)
	}

	return out, nil
}
