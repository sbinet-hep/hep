// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ezplot_test

import (
	"log"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/ezplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func ExampleScatters() {
	var (
		xs = []float64{1, 2, 3, 4, 5}
		ys = []float64{2, 4, 6, 8, 10}
	)

	p := hplot.New()
	p.Title.Text = "easy plot"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	err := ezplot.Scatters(p, ezplot.E{"f(x) = 2*x", hplot.ZipXY(xs, ys)})
	if err != nil {
		log.Fatalf("could not create scatters: %+v", err)
	}

	err = p.Save(20*vg.Centimeter, -1, "testdata/scatter.png")
	if err != nil {
		log.Fatalf("could not save scatter plot: %+v", err)
	}

	// Output:
}

func ExampleScatters_plotutil() {
	var (
		xs = []float64{1, 2, 3, 4, 5}
		ys = []float64{2, 4, 6, 8, 10}
	)

	p := plot.New()
	p.Title.Text = "easy plot"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"

	err := plotutil.AddScatters(p, "f(x) = 2*x", hplot.ZipXY(xs, ys))
	if err != nil {
		log.Fatalf("could not create scatters: %+v", err)
	}

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, "testdata/scatter_pl.png")
	if err != nil {
		log.Fatalf("could not save scatter plot: %+v", err)
	}

	// Output:
}
