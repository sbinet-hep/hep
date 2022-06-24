// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
)

func TestNTuple(t *testing.T) {
	for _, tc := range []struct {
		want rtests.ROOTer
	}{
		{
			want: &NTuple{1, 2, span{1, 2, 3}, span{4, 5, 6}, 7},
		},
	} {
		t.Run("", func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal: %+v", err)
			}

			rt := reflect.Indirect(reflect.ValueOf(tc.want)).Type()
			got := reflect.New(rt).Interface().(rtests.ROOTer)
			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)

			err = got.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal: %+v", err)
			}

			if got, want := got, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid r/w round-trip:\ngot= %#v\nwant=%#v", got, want)
			}
		})
	}
}

func TestReadNTuple(t *testing.T) {
	f, err := riofs.Open("../../testdata/ntpl001_staff.root")
	if err != nil {
		t.Fatalf("could not open file: +%v", err)
	}
	defer f.Close()

	obj, err := f.Get("Staff")
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	nt, ok := obj.(*NTuple)
	if !ok {
		t.Fatalf("%q not an NTuple: %T", "Staff", obj)
	}

	want := NTuple{
		vers: 0x0,
		size: 0x30,
		header: span{
			seek:   854,
			nbytes: 376,
			length: 880,
		},
		footer: span{
			seek:   72449,
			nbytes: 86,
			length: 104,
		},
		reserved: 0,
	}

	if got, want := *nt, want; got != want {
		t.Fatalf("error:\ngot= %#v\nwant=%#v", got, want)
	}

	if got, want := nt.String(), want.String(); got != want {
		t.Fatalf("error:\ngot= %v\nwant=%v", got, want)
	}

	raw, err := os.ReadFile("../../testdata/ntpl001_staff.root")
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	var (
		hdr header
		ftr footer
	)

	{
		dst := make([]byte, want.header.length)
		err = rcompress.Decompress(dst, bytes.NewReader(raw[want.header.seek:want.header.seek+uint64(want.header.nbytes)]))
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
		dec := &decoder{r: &rbuff{p: dst}}
		hdr = dec.decodeHeader()
	}

	{
		dst := make([]byte, want.footer.length)
		err = rcompress.Decompress(dst, bytes.NewReader(raw[want.footer.seek:want.footer.seek+uint64(want.footer.nbytes)]))
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
		dec := &decoder{r: &rbuff{p: dst}}
		ftr = dec.decodeFooter()
	}

	dataOf := func(ftr footer) string {
		pos := ftr.clGroups[0].pages.loc.pos
		zsz := ftr.clGroups[0].pages.loc.storage
		buf := make([]byte, ftr.clGroups[0].pages.size)
		err := rcompress.Decompress(buf, bytes.NewReader(raw[pos:pos+uint64(zsz)]))
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
		dec := &decoder{r: &rbuff{p: buf}}
		ps := dec.decodePagelist()
		var out []string
		for i, cl := range ps.clusters {
			for j, col := range cl.columns {
				for k, p := range col.pages {
					src := raw[p.loc.pos : p.loc.pos+uint64(p.loc.storage)]
					dst := make([]byte, p.nelem*uint32(hdr.cols[j].bits/8))
					t.Logf("col-compr[%d,%d,%d]: 0o%b, %q (src=%d, dst=%d)", i, j, k, col.compr, src[:4], len(src), len(dst))
					switch {
					case len(src) == len(dst):
						dst = src
					default:
						err := rcompress.Decompress(dst, bytes.NewReader(src))
						if err != nil {
							t.Logf("decompress[%d][%d][%d]: [%d:%d+%d] len=%d, src=%d, dst=%d", i, j, k, p.loc.pos, p.loc.pos, p.loc.storage, len(raw), len(src), len(dst))
							t.Fatalf("error: cluster[%d].col[%d].page[%d]: %+v", i, j, k, err)
						}
					}
					fid := hdr.cols[j].fieldID
					out = append(out, fmt.Sprintf("cluster[%d,%d,%d]: %s\n", i, j, k, hdr.fields[fid].fname))
					out = append(out, hex.Dump(dst[:128]))
				}
			}
		}
		return strings.Join(out, "\n")
	}

	t.Fatalf("error:\nheader:\n%+v\nfooter:\n%+v\ndata:\n%+v", hdr, ftr, dataOf(ftr))
}
