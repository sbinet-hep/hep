// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !cross_compile

package main

import (
	"bytes"
	"image/color"
	"math/rand"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	tk "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

const (
	NPOINTS = 100000
)

var plt *tk.TLabelWidget

func main() {
	plt = tk.TLabel(newImg())

	tk.Pack(
		plt,
		tk.TExit(),
		tk.Padx("1m"), tk.Pady("2m"), tk.Ipadx("1m"), tk.Ipady("1m"),
	)
	tk.App.WmTitle("iplot")
	tk.ActivateTheme("azure light")
	tk.App.SetResizable(false, false)

	tk.Bind(tk.App, "<KeyPress-Escape>", tk.ExitHandler())
	tk.Bind(tk.App, "<KeyPress-q>", tk.ExitHandler())
	for _, name := range []string{
		"<KeyPress-KP_Enter>", "<KeyPress-Return>",
		"<KeyPress-space>",
	} {
		tk.Bind(tk.App, name, tk.Command(func() {
			plt.Configure(newImg())
		}))
	}

	tk.App.Wait()
}

func newImg() tk.Opt {
	w, h := hplot.Dims(-1, -1)
	c := vgimg.PngCanvas{vgimg.NewWith(vgimg.UseWH(w, h))}

	newPlot().Draw(draw.New(c))

	// FIXME(sbinet): use image.Image when modernc.org/tk9.0@2042105 is available.
	buf := new(bytes.Buffer)
	_, err := c.WriteTo(buf)
	if err != nil {
		panic(err)
	}

	return tk.Image(tk.NewPhoto(tk.Data(buf.Bytes())))
}

func newPlot() *hplot.Plot {
	// Draw some random values from the standard
	// normal distribution.
	hist1 := hbook.NewH1D(100, -5, +5)
	hist2 := hbook.NewH1D(100, -5, +5)
	for i := 0; i < NPOINTS; i++ {
		v1 := rand.NormFloat64() - 1
		v2 := rand.NormFloat64() + 1
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1 := hplot.NewH1D(hist1)
	h1.Infos.Style = hplot.HInfoSummary
	h1.Color = color.Black
	h1.FillColor = nil

	h2 := hplot.NewH1D(hist2)
	h2.Infos.Style = hplot.HInfoNone
	h2.Color = color.RGBA{255, 0, 0, 255}
	h2.FillColor = nil

	p.Add(h1, h2)

	p.Add(plotter.NewGrid())
	return p
}
