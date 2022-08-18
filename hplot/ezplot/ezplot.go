// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ezplot provides a set of high-level functions to ease
// plotting data.
package ezplot // import "go-hep.org/x/hep/hplot/ezplot"

import (
	"fmt"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

type E struct {
	Name string
	XYs  plotter.XYer
}

type item struct {
	name  string
	value plot.Thumbnailer
}

func Scatters(p *hplot.Plot, vs ...E) error {
	var (
		ps    []plot.Plotter
		items []item
	)

	for i, v := range vs {
		s, err := hplot.NewScatter(v.XYs)
		if err != nil {
			return fmt.Errorf("ezplot: could not create scatter: %w", err)
		}
		s.Color = plotutil.Color(i)
		s.Shape = plotutil.Shape(i)

		ps = append(ps, s)
		if v.Name != "" {
			items = append(items, item{name: v.Name, value: s})
		}
	}

	p.Add(ps...)
	for _, v := range items {
		p.Legend.Add(v.name, v.value)
	}

	return nil
}
