// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
)

type Schema struct {
	Name   string
	Fields []Field

	StreamerInfo rbytes.StreamerInfo
}

type Field struct {
	Name    string
	Type    reflect.Type
	Element rbytes.StreamerElement
}

func SchemaFrom(si rbytes.StreamerInfo, sictx rbytes.StreamerInfoContext) (*Schema, error) {
	return newSchema(si.Name(), sictx, si, si.Elements())
}

func NewSchema(name string, sictx rbytes.StreamerInfoContext, elmts []rbytes.StreamerElement) (*Schema, error) {
	return newSchema(name, sictx, nil, elmts)
}

func newSchema(name string, sictx rbytes.StreamerInfoContext, si rbytes.StreamerInfo, elmts []rbytes.StreamerElement) (schema *Schema, err error) {
	schema = &Schema{
		Name:   name,
		Fields: make([]Field, 0, len(elmts)),
	}

	defer func() {
		e := recover()
		if e != nil {
			switch e := e.(type) {
			case error:
				err = fmt.Errorf("rdict: could not create schema: %w", e)
				schema = nil
			default:
				err = fmt.Errorf("rdict: could not create schema: %v", e)
				schema = nil
			}
		}
	}()

	for _, se := range elmts {
		rt := genTypeFromSE(sictx, se)
		schema.Fields = append(schema.Fields, Field{
			Name:    se.Name(),
			Type:    rt,
			Element: se,
		})
	}

	return schema, nil
}
