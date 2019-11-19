// Copyright 2019 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dflow exposes functions and types to represent a data-flow dependency graph.
package dflow // import "go-hep.org/x/hep/fwk/internal/dflow"

import (
	"reflect"
	"sort"
	"strings"

	"golang.org/x/xerrors"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type Graph struct {
	dg    *simple.DirectedGraph
	ids   map[string]int64
	nodes map[string]*Node
	edges map[string]*Edge
}

type Node struct {
	id   int64
	name string
	in   map[string]reflect.Type
	out  map[string]reflect.Type
}

func (n *Node) Name() string                 { return n.name }
func (n *Node) ID() int64                    { return n.id }
func (n *Node) In() map[string]reflect.Type  { return n.in }
func (n *Node) Out() map[string]reflect.Type { return n.out }

func (n *Node) AddIn(name string, typ reflect.Type) error {
	if _, dup := n.in[name]; dup {
		return xerrors.Errorf("component [%s] already declared in-port with name [%s]",
			n.name, name,
		)
	}
	n.in[name] = typ
	return nil
}

type Edge struct {
	Name string
	From []int64
	To   []int64
}

func New() *Graph {
	return &Graph{
		dg:    simple.NewDirectedGraph(),
		ids:   make(map[string]int64),
		nodes: make(map[string]*Node),
		edges: make(map[string]*Edge),
	}
}

func (g *Graph) Has(name string) bool {
	_, ok := g.nodes[name]
	return ok
}

func (g *Graph) Add(name string, in []string, out []string) error {
	if _, dup := g.nodes[name]; dup {
		return xerrors.Errorf("duplicate node %q", name)
	}

	if dups := dups(in); len(dups) > 0 {
		return xerrors.Errorf("duplicate inputs for node %q: %v", name, dups)
	}
	if dups := dups(out); len(dups) > 0 {
		return xerrors.Errorf("duplicate outputs for node %q: %v", name, dups)
	}

	n := &Node{
		name: name,
		id:   int64(len(g.nodes) + 1), // id must not be zero
		in:   make(map[string]reflect.Type, len(in)),
		out:  make(map[string]reflect.Type, len(out)),
	}
	for _, v := range in {
		n.in[v] = nil // FIXME
		e, ok := g.edges[v]
		if !ok {
			e = &Edge{Name: v}
			g.edges[v] = e
		}
		e.To = append(e.To, n.id)
	}
	for _, v := range out {
		n.out[v] = nil // FIXME
		e, ok := g.edges[v]
		if !ok {
			e = &Edge{Name: v}
			g.edges[v] = e
		}
		e.From = append(e.From, n.id)
	}

	g.ids[name] = n.id
	g.nodes[name] = n
	g.dg.AddNode(n)

	return nil
}

func (g *Graph) build() (*simple.DirectedGraph, error) {
	names := make([]string, 0, len(g.nodes))
	for name := range g.nodes {
		names = append(names, name)
	}
	sort.Strings(names)

	// make sure all inputs of nodes are available as outputs of another node
	// detect whether an output is labeled as such by only 1 node.
	out := make(map[string]string) // outport-name -> node-name
	for _, name := range names {
		node := g.nodes[name]
		for k := range node.out {
			n, dup := out[k]
			if dup {
				return nil, xerrors.Errorf("node %q already declared %q as its output (dup-node=%q)", n, k, name)
			}
			out[k] = name
		}
	}

	for _, name := range names {
		node := g.nodes[name]
		for k := range node.in {
			_, ok := out[k]
			if !ok {
				return nil, xerrors.Errorf("node %q declared %q as input but NO KNOWN produced for it", name, k)
			}
		}
	}

	for _, edge := range g.edges {
		for _, from := range edge.From {
			for _, to := range edge.To {
				var (
					from = g.dg.Node(from)
					to   = g.dg.Node(to)
				)
				g.dg.SetEdge(simple.Edge{F: from, T: to})
			}
		}
	}

	dg := g.dg

	g.dg = simple.NewDirectedGraph()
	for _, n := range g.nodes {
		g.dg.AddNode(n)
	}

	return dg, nil
}

func (g *Graph) Analyze() error {
	dg, err := g.build()
	if err != nil {
		return xerrors.Errorf("could not build graph for analysis: %w", err)
	}

	return g.check(dg)
}

func (g *Graph) check(dg *simple.DirectedGraph) error {
	sccs := topo.TarjanSCC(dg)
	for _, c := range sccs {
		if len(c) == 1 {
			continue
		}
		cycle := make([]string, 0, len(c))
		for _, n := range c {
			cycle = append(cycle, n.(*Node).name)
		}
		return xerrors.Errorf("cycle detected: %v", strings.Join(cycle, " -> "))
	}
	return nil
}

func dups(vs []string) []string {
	var (
		dups []string
		set  = make(map[string]struct{}, len(vs))
	)
	for _, v := range vs {
		if _, dup := set[v]; dup {
			dups = append(dups, v)
			continue
		}
		set[v] = struct{}{}
	}
	return dups
}

var (
	_ graph.Node = (*Node)(nil)
)
