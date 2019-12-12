// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Line implements the Plotter interface, drawing a line.
type Line struct {
	// XYs is a copy of the points for this line.
	plotter.XYs

	// StepStyle is the kind of the step line.
	StepStyle plotter.StepKind

	// LineStyle is the style of the line connecting the points.
	// Use zero width to disable lines.
	draw.LineStyle

	// FillColor is the color to fill the area below the plot.
	// Use nil to disable the filling. This is the default.
	FillColor color.Color

	// LogY allows rendering with a log-scaled Y axis.
	// When enabled, function values returning 0 will be discarded from
	// the final plot.
	LogY bool
}

// NewLine returns a Line that uses the default line style and
// does not draw glyphs.
func NewLine(xys plotter.XYer) (*Line, error) {
	data, err := plotter.CopyXYs(xys)
	if err != nil {
		return nil, err
	}
	return &Line{
		XYs:       data,
		LineStyle: plotter.DefaultLineStyle,
	}, nil
}

// Plot draws the Line, implementing the plot.Plotter interface.
func (pts *Line) Plot(c draw.Canvas, plt *plot.Plot) {
	trX, trY := plt.Transforms(&c)
	ps := make([]vg.Point, len(pts.XYs))

	for i, p := range pts.XYs {
		ps[i].X = trX(p.X)
		ps[i].Y = trY(p.Y)
	}

	if pts.FillColor != nil && len(ps) > 0 {
		minY := trY(plt.Y.Min)
		fillPoly := []vg.Point{{X: ps[0].X, Y: minY}}
		switch pts.StepStyle {
		case plotter.PreStep:
			fillPoly = append(fillPoly, ps[1:]...)
		case plotter.PostStep:
			fillPoly = append(fillPoly, ps[:len(ps)-1]...)
		default:
			fillPoly = append(fillPoly, ps...)
		}
		fillPoly = append(fillPoly, vg.Point{X: ps[len(ps)-1].X, Y: minY})
		fillPoly = c.ClipPolygonXY(fillPoly)
		if len(fillPoly) > 0 {
			c.SetColor(pts.FillColor)
			var pa vg.Path
			prev := fillPoly[0]
			pa.Move(prev)
			for _, pt := range fillPoly[1:] {
				switch pts.StepStyle {
				case plotter.NoStep:
					pa.Line(pt)
				case plotter.PreStep:
					pa.Line(vg.Point{X: prev.X, Y: pt.Y})
					pa.Line(pt)
				case plotter.MidStep:
					pa.Line(vg.Point{X: (prev.X + pt.X) / 2, Y: prev.Y})
					pa.Line(vg.Point{X: (prev.X + pt.X) / 2, Y: pt.Y})
					pa.Line(pt)
				case plotter.PostStep:
					pa.Line(vg.Point{X: pt.X, Y: prev.Y})
					pa.Line(pt)
				}
				prev = pt
			}
			pa.Close()
			c.Fill(pa)
		}
	}

	lines := c.ClipLinesXY(ps)
	if pts.LineStyle.Width != 0 && len(lines) != 0 {
		c.SetLineStyle(pts.LineStyle)
		for _, l := range lines {
			if len(l) == 0 {
				continue
			}
			var p vg.Path
			prev := l[0]
			p.Move(prev)
			for _, pt := range l[1:] {
				switch pts.StepStyle {
				case plotter.PreStep:
					p.Line(vg.Point{X: prev.X, Y: pt.Y})
				case plotter.MidStep:
					p.Line(vg.Point{X: (prev.X + pt.X) / 2, Y: prev.Y})
					p.Line(vg.Point{X: (prev.X + pt.X) / 2, Y: pt.Y})
				case plotter.PostStep:
					p.Line(vg.Point{X: pt.X, Y: prev.Y})
				}
				p.Line(pt)
				prev = pt
			}
			c.Stroke(p)
		}
	}
}

// DataRange returns the minimum and maximum
// x and y values, implementing the plot.DataRanger interface.
func (pts *Line) DataRange() (xmin, xmax, ymin, ymax float64) {
	return plotter.XYRange(pts)
}

// Thumbnail returns the thumbnail for the Line, implementing the plot.Thumbnailer interface.
func (pts *Line) Thumbnail(c *draw.Canvas) {
	if pts.FillColor != nil {
		var topY vg.Length
		if pts.LineStyle.Width == 0 {
			topY = c.Max.Y
		} else {
			topY = (c.Min.Y + c.Max.Y) / 2
		}
		points := []vg.Point{
			{X: c.Min.X, Y: c.Min.Y},
			{X: c.Min.X, Y: topY},
			{X: c.Max.X, Y: topY},
			{X: c.Max.X, Y: c.Min.Y},
		}
		poly := c.ClipPolygonY(points)
		c.FillPolygon(pts.FillColor, poly)
	}

	if pts.LineStyle.Width != 0 {
		y := c.Center().Y
		c.StrokeLine2(pts.LineStyle, c.Min.X, y, c.Max.X, y)
	}
}

// NewLinePoints returns both a Line and a
// Points for the given point data.
func NewLinePoints(xys plotter.XYer) (*Line, *plotter.Scatter, error) {
	s, err := NewScatter(xys)
	if err != nil {
		return nil, nil, err
	}
	l := &Line{
		XYs:       s.XYs,
		LineStyle: plotter.DefaultLineStyle,
	}
	return l, s, nil
}

// VertLine draws a vertical line at X and colors the
// left and right portions of the plot with the provided
// colors.
type VertLine struct {
	X     float64
	Line  draw.LineStyle
	Left  color.Color
	Right color.Color
}

// VLine creates a vertical line at x with the default line style.
func VLine(x float64, left, right color.Color) *VertLine {
	return &VertLine{
		X:     x,
		Line:  plotter.DefaultLineStyle,
		Left:  left,
		Right: right,
	}
}

func (vline *VertLine) Plot(c draw.Canvas, plt *plot.Plot) {
	var (
		trX, _ = plt.Transforms(&c)
		x      = trX(vline.X)
		xmin   = c.Min.X
		xmax   = c.Max.X
		ymin   = c.Min.Y
		ymax   = c.Max.Y
	)

	if vline.Left != nil {
		c.SetColor(vline.Left)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: x, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if vline.Right != nil {
		c.SetColor(vline.Right)
		rect := vg.Rectangle{
			Min: vg.Point{X: x, Y: ymin},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}

	if vline.Line.Width != 0 {
		c.StrokeLine2(vline.Line, x, ymin, x, ymax)
	}
}

// HorizLine draws a horizontal line at Y and colors the
// top and bottom portions of the plot with the provided
// colors.
type HorizLine struct {
	Y      float64
	Line   draw.LineStyle
	Top    color.Color
	Bottom color.Color
}

// HLine creates a horizontal line at y with the default line style.
func HLine(y float64, top, bottom color.Color) *HorizLine {
	return &HorizLine{
		Y:      y,
		Line:   plotter.DefaultLineStyle,
		Top:    top,
		Bottom: bottom,
	}
}

func (hline *HorizLine) Plot(c draw.Canvas, plt *plot.Plot) {
	var (
		_, trY = plt.Transforms(&c)
		y      = trY(hline.Y)
		xmin   = c.Min.X
		xmax   = c.Max.X
		ymin   = c.Min.Y
		ymax   = c.Max.Y
	)

	if hline.Top != nil {
		c.SetColor(hline.Top)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: y},
			Max: vg.Point{X: xmax, Y: ymax},
		}
		c.Fill(rect.Path())
	}
	if hline.Bottom != nil {
		c.SetColor(hline.Bottom)
		rect := vg.Rectangle{
			Min: vg.Point{X: xmin, Y: ymin},
			Max: vg.Point{X: xmax, Y: y},
		}
		c.Fill(rect.Path())
	}

	if hline.Line.Width != 0 {
		c.StrokeLine2(hline.Line, xmin, y, xmax, y)
	}
}

var (
	_ plot.Plotter = (*VertLine)(nil)
	_ plot.Plotter = (*HorizLine)(nil)
)
