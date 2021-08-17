// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

// ReadVar describes a variable to be read out of a tree.
type ReadVar struct {
	Name  string      // name of the branch to read
	Leaf  string      // name of the leaf to read
	Value interface{} // pointer to the value to fill

	count string // name of the leaf-count, if any
	leaf  Leaf   // leaf to which this read-var is bound
}

// NewReadVars returns the complete set of ReadVars to read all the data
// contained in the provided Tree.
func NewReadVars(t Tree) []ReadVar {
	var vars []ReadVar
	for _, b := range t.Branches() {
		for _, leaf := range b.Leaves() {
			ptr := newValue(leaf)
			log.Printf("leaf[%s]: %T count=%v", leaf.Name(), leaf, leaf.LeafCount() != nil)
			bname := b.Name()
			lname := leaf.Name()
			if l, ok := leaf.(*tleafElement); ok && strings.HasSuffix(lname, "_") {
				b := l.branch.(*tbranchElement)
				log.Printf("->b= bcount1=%v, clones=%q", b.bcount1, b.clones)
				ptr = rcont.NewClonesArrayOf(b.clones)
			}
			cnt := ""
			if leaf.LeafCount() != nil {
				cnt = leaf.LeafCount().Name()
			}
			vars = append(vars, ReadVar{
				Name:  bname,
				Leaf:  lname,
				Value: ptr,
				count: cnt,
				leaf:  leaf,
			})
		}
	}
	return vars
}

// Deref returns the value pointed at by this read-var.
func (rv ReadVar) Deref() interface{} {
	return reflect.ValueOf(rv.Value).Elem().Interface()
}

// ReadVarsFromStruct returns a list of ReadVars bound to the exported fields
// of the provided pointer to a struct value.
//
// ReadVarsFromStruct panicks if the provided value is not a pointer to
// a struct value.
func ReadVarsFromStruct(ptr interface{}) []ReadVar {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("rtree: expect a pointer value, got %T", ptr))
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("rtree: expect a pointer to struct value, got %T", ptr))
	}

	var (
		rt     = rv.Type()
		rvars  = make([]ReadVar, 0, rt.NumField())
		reDims = regexp.MustCompile(`\w*?\[(\w*)\]+?`)
	)

	split := func(s string) (string, []string) {
		n := s
		if i := strings.Index(s, "["); i > 0 {
			n = s[:i]
		}

		out := reDims.FindAllStringSubmatch(s, -1)
		if len(out) == 0 {
			return n, nil
		}

		dims := make([]string, len(out))
		for i := range out {
			dims[i] = out[i][1]
		}
		return n, dims
	}

	for i := 0; i < rt.NumField(); i++ {
		var (
			ft = rt.Field(i)
			fv = rv.Field(i)
		)
		if ft.Name != strings.Title(ft.Name) {
			// not exported. ignore.
			continue
		}
		rvar := ReadVar{
			Name:  nameOf(ft),
			Value: fv.Addr().Interface(),
		}

		if strings.Contains(rvar.Name, "[") {
			switch ft.Type.Kind() {
			case reflect.Slice:
				sli, dims := split(rvar.Name)
				if len(dims) > 1 {
					panic(fmt.Errorf("rtree: invalid number of slice-dimensions for field %q: %q", ft.Name, rvar.Name))
				}
				rvar.Name = sli
				rvar.count = dims[0]

			case reflect.Array:
				arr, dims := split(rvar.Name)
				if len(dims) > 3 {
					panic(fmt.Errorf("rtree: invalid number of array-dimension for field %q: %q", ft.Name, rvar.Name))
				}
				rvar.Name = arr
			default:
				panic(fmt.Errorf("rtree: invalid field type for %q, or invalid struct-tag %q: %T", ft.Name, rvar.Name, fv.Interface()))
			}
		}
		switch ft.Type.Kind() {
		case reflect.Int, reflect.Uint, reflect.UnsafePointer, reflect.Uintptr, reflect.Chan, reflect.Interface:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T", ft.Name, fv.Interface()))
		case reflect.Map:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T (not yet supported)", ft.Name, fv.Interface()))
		}

		rvar.Leaf = rvar.Name
		rvars = append(rvars, rvar)
	}
	return rvars
}

func nameOf(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("groot")
	if ok {
		return tag
	}
	return field.Name
}

func bindRVarsTo(t Tree, rvars []ReadVar) []ReadVar {
	var (
		ors    = make([]ReadVar, 0, len(rvars))
		counts []ReadVar

		flatten func(b Branch, rvar ReadVar) []ReadVar
	)

	flatten = func(br Branch, rvar ReadVar) []ReadVar {
		nsub := len(br.Branches())
		subs := make([]ReadVar, 0, nsub)
		log.Printf(">>> flatten(%q, %q)... nsub=%d", br.Name(), rvar.Name, nsub)
		rv := reflect.ValueOf(rvar.Value).Elem()
		get := func(name string) int {
			rt := rv.Type()
			switch rt.Kind() {
			case reflect.Struct:
				for i := 0; i < rt.NumField(); i++ {
					ft := rt.Field(i)
					nn := nameOf(ft)
					if nn == name {
						// exact match.
						return i
					}
					// try to remove any [xyz][range].
					// do it after exact match not to shortcut arrays
					if idx := strings.Index(nn, "["); idx > 0 {
						nn = string(nn[:idx])
					}
					if nn == name {
						return i
					}
				}
			case reflect.Map:
				log.Printf("get-map(%s)...", name)
				switch name {
				case "first": // C++ std::map<K,V> -> K
					return 0
				case "second": // C++ std::map<K,V> -> V
					return 1
				}
			default:
				panic(fmt.Errorf("rtree: unhandled reflect kind=%v", rt.Kind()))
			}
			return -1
		}

		for _, sub := range br.Branches() {
			log.Printf("sub[%s.%s]...", br.Name(), sub.Name())
			bn := sub.Name()
			if strings.Contains(bn, ".") {
				toks := strings.Split(bn, ".")
				bn = toks[len(toks)-1]
			}
			j := get(bn)
			if j < 0 {
				continue
			}
			var (
				subrv ReadVar
				bname = sub.Name()
				lname = sub.Name()
			)
			if prefix := br.Name() + "."; strings.HasPrefix(bname, prefix) {
				bname = string(bname[len(prefix):])
			}
			if idx := strings.Index(bname, "["); idx > 0 {
				bname = string(bname[:idx])
			}
			if idx := strings.Index(lname, "["); idx > 0 {
				lname = string(lname[:idx])
			}

			var (
				leaf  = sub.Leaf(lname)
				count = ""
			)
			if leaf != nil {
				if lc := leaf.LeafCount(); lc != nil {
					count = lc.Name()
				}
			}
			switch rv.Kind() {
			case reflect.Struct:
				fv := rv.Field(j)
				subrv = ReadVar{
					Name:  rvar.Name + "." + bname,
					Leaf:  lname,
					Value: fv.Addr().Interface(),
					leaf:  leaf,
					count: count,
				}
				if count != "" && count == rvar.Name+"_" {
					// add count leaf.
					counts = append(counts, ReadVar{
						Name:  rvar.Name + ".count__groot",
						Leaf:  count,
						Value: reflect.New(leaf.LeafCount().Type()).Interface(),
						leaf:  t.Leaf(count),
					})
				}

			case reflect.Map:
				var (
					kt reflect.Type
					rt = rv.Type()
				)
				switch {
				case strings.HasSuffix(bname, "first"):
					rt = rt.Key()
					kt = rt
				case strings.HasSuffix(bname, "second"):
					rt = rt.Elem()
				}

				log.Printf("sub: %T, name=%q, rvar=%q --> %q (count=%q, %v)",
					rv.Interface(), sub.Name(), rvar.Name, rvar.Name+"."+bname,
					count, kt,
				)
				subrv = ReadVar{
					Name:  rvar.Name + "." + bname,
					Leaf:  lname,
					Value: reflect.New(rt).Interface(),
					leaf:  leaf,
					count: count,
				}

				if kt != nil && count != "" {
					// add count leaf.
					counts = append(counts, ReadVar{
						Name:  rvar.Name + ".count",
						Leaf:  count,
						Value: reflect.New(leaf.LeafCount().Type()).Interface(),
						leaf:  t.Leaf(count),
					})
				}

			default:
				panic(fmt.Errorf("rtree: unhandled reflect kind=%v", rv.Kind()))
			}
			switch len(sub.Branches()) {
			case 0:
				subs = append(subs, subrv)
			default:
				subs = append(subs, flatten(sub, subrv)...)
			}
		}
		return subs
	}

	for i := range rvars {
		var (
			rvar = &rvars[i]
			br   = t.Branch(rvar.Name)
			leaf = br.Leaf(rvar.Leaf)
			nsub = len(br.Branches())
		)
		switch nsub {
		case 0:
			rvar.leaf = leaf
			ors = append(ors, *rvar)
		default:
			ors = append(ors, flatten(br, *rvar)...)
		}
	}

	log.Printf("tree[%s] leaves=%d", t.Name(), len(t.Leaves()))
	for i, leaf := range t.Leaves() {
		log.Printf("leaf[%d]: %q (%T)", i, leaf.Name(), leaf)
	}

	if len(counts) > 0 {
		ors = append(counts, ors...)
	}
	return ors
}

func newValue(leaf Leaf) interface{} {
	etype := leaf.Type()
	unsigned := leaf.IsUnsigned()

	switch etype.Kind() {
	case reflect.Interface, reflect.Map, reflect.Chan:
		panic(fmt.Errorf("rtree: type %T not supported", reflect.New(etype).Elem().Interface()))
	case reflect.Int8:
		if unsigned {
			etype = reflect.TypeOf(uint8(0))
		}
	case reflect.Int16:
		if unsigned {
			etype = reflect.TypeOf(uint16(0))
		}
	case reflect.Int32:
		if unsigned {
			etype = reflect.TypeOf(uint32(0))
		}
	case reflect.Int64:
		if unsigned {
			etype = reflect.TypeOf(uint64(0))
		}
	case reflect.Float32:
		if _, ok := leaf.(*LeafF16); ok {
			etype = reflect.TypeOf(root.Float16(0))
		}
	case reflect.Float64:
		if _, ok := leaf.(*LeafD32); ok {
			etype = reflect.TypeOf(root.Double32(0))
		}
	}

	switch {
	case leaf.LeafCount() != nil:
		shape := leaf.Shape()
		switch leaf.(type) {
		case *LeafF16, *LeafD32:
			// workaround for https://sft.its.cern.ch/jira/browse/ROOT-10149
			shape = nil
		}
		for i := range shape {
			etype = reflect.ArrayOf(shape[len(shape)-1-i], etype)
		}
		etype = reflect.SliceOf(etype)
	case leaf.Len() > 1:
		shape := leaf.Shape()
		switch leaf.Kind() {
		case reflect.String:
			switch dims := len(shape); dims {
			case 0, 1:
				// interpret as a single string.
			default:
				// FIXME(sbinet): properly handle [N]string (but ROOT doesn't support that.)
				// see: https://root-forum.cern.ch/t/char-t-in-a-branch/5591/2
				// etype = reflect.ArrayOf(leaf.Len(), etype)
				panic(fmt.Errorf("groot/rtree: invalid number of dimensions (%d)", dims))
			}
		default:
			switch leaf.(type) {
			case *LeafF16, *LeafD32:
				// workaround for https://sft.its.cern.ch/jira/browse/ROOT-10149
				shape = []int{leaf.Len()}
			}
			for i := range shape {
				etype = reflect.ArrayOf(shape[len(shape)-1-i], etype)
			}
		}
	}
	return reflect.New(etype).Interface()
}
