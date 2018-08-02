// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kalman_test

import (
	"math"
	"math/rand"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/kalman"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// func Example() {
// 	n := 3 // number of states
// 	m := 1 // number of measurements
//
// 	dt := 1.0 / 30.0 // time step
//
// 	a := mat.NewDense(n, n, []float64{1, dt, 0, 0, 1, dt, 0, 0, 1})
// 	c := mat.NewDense(m, n, []float64{1, 0, 0})
//
// 	q := mat.NewDense(n, n, []float64{0.05, 0.05, 0, 0.05, 0.05, 0, 0, 0, 0})
// 	r := mat.NewDense(m, m, []float64{5})
// 	p := mat.NewDense(n, n, []float64{0.1, 0.1, 0.1, 0.1, 10000, 10, 0.1, 10, 100})
//
// 	kf := kalman.New(dt, c, q, r, p)
//
// 	// measurements
// 	data := []float64{
// 		1.04202710058, 1.10726790452, 1.2913511148, 1.48485250951, 1.72825901034,
// 		1.74216489744, 2.11672039768, 2.14529225112, 2.16029641405, 2.21269371128,
// 		2.57709350237, 2.6682215744, 2.51641839428, 2.76034056782, 2.88131780617,
// 		2.88373786518, 2.9448468727, 2.82866600131, 3.0006601946, 3.12920591669,
// 		2.858361783, 2.83808170354, 2.68975330958, 2.66533185589, 2.81613499531,
// 		2.81003612051, 2.88321849354, 2.69789264832, 2.4342229249, 2.23464791825,
// 		2.30278776224, 2.02069770395, 1.94393985809, 1.82498398739, 1.52526230354,
// 		1.86967808173, 1.18073207847, 1.10729605087, 0.916168349913, 0.678547664519,
// 		0.562381751596, 0.355468474885, -0.155607486619, -0.287198661013, -0.602973173813,
// 	}
//
// 	x0 := mat.NewVecDense(n, []float64{data[0], 0, -9.81})
//
// 	kf.Init(0, x0)
//
// 	y := mat.NewVecDense(m, nil)
// 	t := 0.0
// 	for i, v := range data {
// 		t += dt
// 		y.SetVec(i, v)
// 		err := kf.Update(y, dt, a)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
//
// 	fmt.Printf("state=%v\n", kf.State())
//
// 	// Output:
// 	// asd
// }

func Example() {

	sstd := 0.000001
	ostd := 0.1

	// trend model
	kf := kalman.NewKF(&kalman.Settings{
		F: mat.NewDense(2, 2, []float64{2, -1, 1, 0}),
		G: mat.NewDense(2, 1, []float64{1, 0}),
		Q: mat.NewDense(1, 1, []float64{sstd}),
		H: mat.NewDense(1, 2, []float64{1, 0}),
		R: mat.NewDense(1, 1, []float64{ostd}),
	})

	n := 10000
	sys := mat.NewDense(1, n, nil)
	x, dx := 0.0, 0.01
	xs := make([]float64, 0, n)
	ys := make([]float64, 0, n)

	for i := 0; i < n; i++ {
		y := math.Sin(x) + 0.1*(rand.NormFloat64()-0.5)
		sys.Set(0, i, y)
		x += dx

		xs = append(xs, x)
		ys = append(ys, y)
	}

	out, err := kf.Filter(nil, sys)
	if err != nil {
		panic(err)
	}
	oys := mat.Row(nil, 0, out)

	p := hplot.New()
	err = plotutil.AddLinePoints(p.Plot,
		"Original", hplot.ZipXY(xs, ys),
		"Filtered", hplot.ZipXY(xs, oys),
	)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Centimeter, -1, "sample.png"); err != nil {
		panic(err)
	}

	// Output:
}

func zipXY(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))

	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}

	return pts
}
