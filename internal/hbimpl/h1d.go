// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbimpl

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
)

// H1D is a 1-dim histogram with weighted entries.
type H1D struct {
	XXX_bng Binning1D
	XXX_ann Annotation
}

// NewH1D returns a 1-dim histogram with n bins between xmin and xmax.
func NewH1D(n int, xmin, xmax float64) *H1D {
	return &H1D{
		XXX_bng: newBinning1D(n, xmin, xmax),
		XXX_ann: make(Annotation),
	}
}

// NewH1DFromEdges returns a 1-dim histogram given a slice of edges.
// The number of bins is thus len(edges)-1.
// It panics if the length of edges is <= 1.
// It panics if the edges are not sorted.
// It panics if there are duplicate edge values.
func NewH1DFromEdges(edges []float64) *H1D {
	return &H1D{
		XXX_bng: newBinning1DFromEdges(edges),
		XXX_ann: make(Annotation),
	}
}

// NewH1DFromBins returns a 1-dim histogram given a slice of 1-dim bins.
// It panics if the length of bins is < 1.
// It panics if the bins overlap.
func NewH1DFromBins(bins ...Range) *H1D {
	return &H1D{
		XXX_bng: newBinning1DFromBins(bins),
		XXX_ann: make(Annotation),
	}
}

// Name returns the name of this histogram, if any
func (h *H1D) Name() string {
	v, ok := h.XXX_ann["name"]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

// Annotation returns the annotations attached to this histogram
func (h *H1D) Annotation() Annotation {
	return h.XXX_ann
}

// Rank returns the number of dimensions for this histogram
func (h *H1D) Rank() int {
	return 1
}

// Entries returns the number of entries in this histogram
func (h *H1D) Entries() int64 {
	return h.XXX_bng.entries()
}

// EffEntries returns the number of effective entries in this histogram
func (h *H1D) EffEntries() float64 {
	return h.XXX_bng.effEntries()
}

// Binning returns the binning of this histogram
func (h *H1D) Binning() *Binning1D {
	return &h.XXX_bng
}

// SumW returns the sum of weights in this histogram.
// Overflows are included in the computation.
func (h *H1D) SumW() float64 {
	return h.XXX_bng.XXX_dist.SumW()
}

// SumW2 returns the sum of squared weights in this histogram.
// Overflows are included in the computation.
func (h *H1D) SumW2() float64 {
	return h.XXX_bng.XXX_dist.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (h *H1D) SumWX() float64 {
	return h.XXX_bng.XXX_dist.SumWX()
}

// SumWX2 returns the 2nd order weighted x moment
func (h *H1D) SumWX2() float64 {
	return h.XXX_bng.XXX_dist.SumWX2()
}

// XMean returns the mean X.
// Overflows are included in the computation.
func (h *H1D) XMean() float64 {
	return h.XXX_bng.XXX_dist.mean()
}

// XVariance returns the variance in X.
// Overflows are included in the computation.
func (h *H1D) XVariance() float64 {
	return h.XXX_bng.XXX_dist.variance()
}

// XStdDev returns the standard deviation in X.
// Overflows are included in the computation.
func (h *H1D) XStdDev() float64 {
	return h.XXX_bng.XXX_dist.stdDev()
}

// XStdErr returns the standard error in X.
// Overflows are included in the computation.
func (h *H1D) XStdErr() float64 {
	return h.XXX_bng.XXX_dist.stdErr()
}

// XRMS returns the XRMS in X.
// Overflows are included in the computation.
func (h *H1D) XRMS() float64 {
	return h.XXX_bng.XXX_dist.rms()
}

// Fill fills this histogram with x and weight w.
func (h *H1D) Fill(x, w float64) {
	h.XXX_bng.fill(x, w)
}

// XMin returns the low edge of the X-axis of this histogram.
func (h *H1D) XMin() float64 {
	return h.XXX_bng.xMin()
}

// XMax returns the high edge of the X-axis of this histogram.
func (h *H1D) XMax() float64 {
	return h.XXX_bng.xMax()
}

// Scale scales the content of each bin by the given factor.
func (h *H1D) Scale(factor float64) {
	h.XXX_bng.scaleW(factor)
}

// Integral computes the integral of the histogram.
//
// The number of parameters can be 0 or 2.
// If 0, overflows are included.
// If 2, the first parameter must be the lower bound of the range in which
// the integral is computed and the second one the upper range.
//
// If the lower bound is math.Inf(-1) then the underflow bin is included.
// If the upper bound is math.Inf(+1) then the overflow bin is included.
//
// Examples:
//
//    // integral of all in-range bins + overflows
//    v := h.Integral()
//
//    // integral of all in-range bins, underflow and overflow bins included.
//    v := h.Integral(math.Inf(-1), math.Inf(+1))
//
//    // integrall of all in-range bins, overflow bin included
//    v := h.Integral(h.Binning().LowerEdge(), math.Inf(+1))
//
//    // integrall of all bins for which the lower edge is in [0.5, 5.5)
//    v := h.Integral(0.5, 5.5)
func (h *H1D) Integral(args ...float64) float64 {
	min, max := 0., 0.
	switch len(args) {
	case 0:
		return h.SumW()
	case 2:
		min = args[0]
		max = args[1]
		if min > max {
			panic("hbook: min > max")
		}
	default:
		panic("hbook: invalid number of arguments. expected 0 or 2.")
	}

	integral := 0.0
	for _, bin := range h.XXX_bng.XXX_bins {
		v := bin.XXX_xrange.Min
		if min <= v && v < max {
			integral += bin.SumW()
		}
	}
	if math.IsInf(min, -1) {
		integral += h.XXX_bng.XXX_outflows[0].SumW()
	}
	if math.IsInf(max, +1) {
		integral += h.XXX_bng.XXX_outflows[1].SumW()
	}
	return integral
}

// Value returns the content of the idx-th bin.
//
// Value implements gonum/plot/plotter.Valuer
func (h *H1D) Value(i int) float64 {
	return h.XXX_bng.XXX_bins[i].SumW()
}

// Len returns the number of bins for this histogram
//
// Len implements gonum/plot/plotter.Valuer
func (h *H1D) Len() int {
	return len(h.XXX_bng.XXX_bins)
}

// XY returns the x,y values for the i-th bin
//
// XY implements gonum/plot/plotter.XYer
func (h *H1D) XY(i int) (float64, float64) {
	bin := h.XXX_bng.XXX_bins[i]
	x := bin.XXX_xrange.Min
	y := bin.SumW()
	return x, y
}

// DataRange implements the gonum/plot.DataRanger interface
func (h *H1D) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = h.XMin()
	xmax = h.XMax()
	ymin = +math.MaxFloat64
	ymax = -math.MaxFloat64
	for _, b := range h.XXX_bng.XXX_bins {
		v := b.SumW()
		ymax = math.Max(ymax, v)
		ymin = math.Min(ymin, v)
	}
	return xmin, xmax, ymin, ymax
}

// // RioMarshal implements rio.RioMarshaler
// func (h *H1D) RioMarshal(w io.Writer) error {
// 	data, err := h.MarshalBinary()
// 	if err != nil {
// 		return err
// 	}
// 	var buf [8]byte
// 	binary.LittleEndian.PutUint64(buf[:], uint64(len(data)))
// 	_, err = w.Write(buf[:])
// 	if err != nil {
// 		return err
// 	}
// 	_, err = w.Write(data)
// 	return err
// }
//
// // RioUnmarshal implements rio.RioUnmarshaler
// func (h *H1D) RioUnmarshal(r io.Reader) error {
// 	buf := make([]byte, 8)
// 	_, err := io.ReadFull(r, buf)
// 	if err != nil {
// 		return err
// 	}
// 	n := int64(binary.LittleEndian.Uint64(buf))
// 	buf = make([]byte, int(n))
// 	_, err = io.ReadFull(r, buf)
// 	if err != nil {
// 		return err
// 	}
// 	return h.UnmarshalBinary(buf)
// }
//
// // RioVersion implements rio.RioStreamer
// func (h *H1D) RioVersion() rio.Version {
// 	return 0
// }

// annToYODA creates a new Annotation with fields compatible with YODA
func (h *H1D) annToYODA() Annotation {
	ann := make(Annotation, len(h.XXX_ann))
	ann["Type"] = "Histo1D"
	ann["Path"] = "/" + h.Name()
	ann["Title"] = ""
	for k, v := range h.XXX_ann {
		if k == "name" {
			continue
		}
		ann[k] = v
	}
	return ann
}

// annFromYODA creates a new Annotation from YODA compatible fields
func (h *H1D) annFromYODA(ann Annotation) {
	if len(h.XXX_ann) == 0 {
		h.XXX_ann = make(Annotation, len(ann))
	}
	for k, v := range ann {
		switch k {
		case "Type":
			// noop
		case "Path":
			h.XXX_ann["name"] = string(v.(string)[1:]) // skip leading '/'
		default:
			h.XXX_ann[k] = v
		}
	}
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H1D) MarshalYODA() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := h.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_HISTO1D %s\n", ann["Path"])
	data, err := ann.MarshalYODA()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	fmt.Fprintf(buf, "# Mean: %e\n", h.XMean())
	fmt.Fprintf(buf, "# Area: %e\n", h.Integral())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")
	d := h.XXX_bng.XXX_dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)
	d = h.XXX_bng.XXX_outflows[0]
	fmt.Fprintf(
		buf,
		"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)

	d = h.XXX_bng.XXX_outflows[1]
	fmt.Fprintf(
		buf,
		"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")
	for _, bin := range h.XXX_bng.XXX_bins {
		d := bin.XXX_dist
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			bin.XXX_xrange.Min, bin.XXX_xrange.Max,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H1D) UnmarshalYODA(data []byte) error {
	r := bytes.NewBuffer(data)
	_, err := readYODAHeader(r, "BEGIN YODA_HISTO1D")
	if err != nil {
		return err
	}
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# Mean:"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid H1D-YODA data")
	}
	err = ann.UnmarshalYODA(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %v\nhbook: %q", err, string(r.Bytes()[:pos+1]))
	}
	h.annFromYODA(ann)
	r.Next(pos)

	var ctx struct {
		total bool
		under bool
		over  bool
		bins  bool
	}

	// sets of xlow values, to infer number of bins in X.
	xset := make(map[float64]int)

	var (
		dist   Dist1D
		oflows [2]Dist1D
		bins   []Bin1D
		xmin   = math.Inf(+1)
		xmax   = math.Inf(-1)
	)
	s := bufio.NewScanner(r)
scanLoop:
	for s.Scan() {
		buf := s.Bytes()
		if len(buf) == 0 || buf[0] == '#' {
			continue
		}
		rbuf := bytes.NewReader(buf)
		switch {
		case bytes.HasPrefix(buf, []byte("END YODA_HISTO1D")):
			break scanLoop
		case !ctx.total && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.total = true
			d := &dist
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.XXX_sumW, &d.dist.XXX_sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.XXX_n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
		case !ctx.under && bytes.HasPrefix(buf, []byte("Underflow\t")):
			ctx.under = true
			d := &oflows[0]
			_, err = fmt.Fscanf(
				rbuf,
				"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.XXX_sumW, &d.dist.XXX_sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.XXX_n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
		case !ctx.over && bytes.HasPrefix(buf, []byte("Overflow\t")):
			ctx.over = true
			d := &oflows[1]
			_, err = fmt.Fscanf(
				rbuf,
				"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.XXX_sumW, &d.dist.XXX_sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.XXX_n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			ctx.bins = true
		case ctx.bins:
			var bin Bin1D
			d := &bin.XXX_dist
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&bin.XXX_xrange.Min, &bin.XXX_xrange.Max,
				&d.dist.XXX_sumW, &d.dist.XXX_sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.XXX_n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			xset[bin.XXX_xrange.Min] = 1
			xmin = math.Min(xmin, bin.XXX_xrange.Min)
			xmax = math.Max(xmax, bin.XXX_xrange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid H1D-YODA data: %q", string(buf))
		}
	}
	h.XXX_bng = Binning1D{
		XXX_bins:     bins,
		XXX_dist:     dist,
		XXX_outflows: oflows,
		XXX_xrange:   Range{xmin, xmax},
	}
	return err
}

// check various interfaces
var _ Object = (*H1D)(nil)
var _ Histogram = (*H1D)(nil)

// serialization interfaces
// var _ rio.Marshaler = (*H1D)(nil)
// var _ rio.Unmarshaler = (*H1D)(nil)
// var _ rio.Streamer = (*H1D)(nil)

func init() {
	gob.Register((*H1D)(nil))
}
