// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdtype

import "reflect"

type Schema struct {
	fields []Field        // fields for this schema description
	index  map[string]int // field-name to field-index
}

func NewSchema(fields []Field) Schema {
	sc := Schema{
		fields: make([]Field, len(fields)),
		index:  make(map[string]int, len(fields)),
	}
	for i, f := range fields {
		sc.fields[i] = f
		sc.index[f.Name] = i
	}
	return sc
}

func (sc Schema) Fields() []Field   { return sc.fields }
func (sc Schema) Field(i int) Field { return sc.fields[i] }

func (sc Schema) FieldByName(n string) (Field, bool) {
	i, ok := sc.index[n]
	if !ok {
		return Field{}, ok
	}
	return sc.fields[i], ok
}

type Field struct {
	Name string
	Type reflect.Type
}
