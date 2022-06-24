// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

import (
	"encoding/binary"
	"fmt"
	"log"
)

type decoder struct {
	r   *rbuff
	err error
}

func (dec *decoder) read(p []byte) (int, error) {
	if dec.err != nil {
		return 0, dec.err
	}
	var n int
	n, dec.err = dec.r.Read(p)
	return n, dec.err
}

func (dec *decoder) u16() uint16 {
	if dec.err != nil {
		return 0
	}
	beg := dec.r.c
	dec.r.c += 2
	return binary.LittleEndian.Uint16(dec.r.p[beg:dec.r.c])
}

func (dec *decoder) u32() uint32 {
	if dec.err != nil {
		return 0
	}
	beg := dec.r.c
	dec.r.c += 4
	return binary.LittleEndian.Uint32(dec.r.p[beg:dec.r.c])
}

func (dec *decoder) u64() uint64 {
	if dec.err != nil {
		return 0
	}
	beg := dec.r.c
	dec.r.c += 8
	return binary.LittleEndian.Uint64(dec.r.p[beg:dec.r.c])
}

func (dec *decoder) decodeString() string {
	if dec.err != nil {
		return ""
	}

	n := dec.u32()
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	_, _ = dec.read(buf)
	if dec.err != nil {
		return ""
	}
	return string(buf)
}

func (dec *decoder) decodeHeader() (hdr header) {
	if dec.err != nil {
		return hdr
	}

	hdr.vers = dec.u16()
	hdr.minv = dec.u16()
	log.Printf("hdr: vers=%d, min=%d", hdr.vers, hdr.minv)
	hdr.flags = dec.decodeFlags()
	hdr.release = dec.u32()
	hdr.name = dec.decodeString()
	hdr.descr = dec.decodeString()
	hdr.library = dec.decodeString()
	hdr.fields = dec.decodeFields()
	hdr.cols = dec.decodeCols()
	hdr.aliases = dec.decodeAliases()
	hdr.extra = dec.decodeExtraTypeInfo()

	hdr.crc32 = dec.u32()
	return hdr
}

func (dec *decoder) decodeFooter() (ftr footer) {
	if dec.err != nil {
		return ftr
	}

	ftr.vers = dec.u16()
	ftr.minv = dec.u16()
	log.Printf("ftr: vers=%d, min=%d", ftr.vers, ftr.minv)
	ftr.flags = dec.decodeFlags()
	ftr.hdr = dec.u32()

	ftr.xhdrs = dec.decodeExtHeaders()
	ftr.colGroups = dec.decodeColGroups()
	ftr.clInfos = dec.decodeClusterInfos()
	ftr.clGroups = dec.decodeClusterGroups()
	ftr.mdBlocks = dec.decodeMetaDataBlocks()

	ftr.crc32 = dec.u32()
	return ftr
}

func (dec *decoder) decodeFlags() (flags []uint64) {
	if dec.err != nil {
		return flags
	}

	for {
		f := dec.u64()
		flags = append(flags, f>>1)
		if int64(f) >= 0 {
			break
		}
	}

	return flags
}

func (dec *decoder) decodeFields() (fields []fieldDescr) {
	if dec.err != nil {
		return fields
	}

	size, nfields := dec.decodeFrameHeader()
	if dec.err != nil {
		return fields
	}
	log.Printf("size:%d, n=%d", size, nfields)

	fields = make([]fieldDescr, nfields)
	// 0-field.
	for i := range fields {
		fsz, fnel := dec.decodeFrameHeader()
		log.Printf("field[%d]: %d, %d", i, fsz, fnel)
		f := &fields[i]
		f.vers = dec.u32()
		f.typv = dec.u32()
		f.pfid = dec.u32()
		f.role = fieldRole(dec.u16())
		f.flag = fieldFlag(dec.u16())
		if f.flag&fieldFlagRepetitive != 0 {
			f.nrep = dec.u64()
		}
		f.fname = dec.decodeString()
		f.tname = dec.decodeString()
		f.alias = dec.decodeString()
		f.descr = dec.decodeString()
		log.Printf("field[%d]: %+v", i, f)
	}
	return fields
}

func (dec *decoder) decodeCols() (cols []colDescr) {
	if dec.err != nil {
		return nil
	}

	size, n := dec.decodeFrameHeader()
	if dec.err != nil {
		return nil
	}
	log.Printf("cols: size:%d, n=%d", size, n)

	cols = make([]colDescr, n)
	for i := range cols {
		csz, cnel := dec.decodeFrameHeader()
		log.Printf("col[%d]: %d, %d", i, csz, cnel)
		col := &cols[i]
		col.kind = colKind(dec.u16())
		col.bits = dec.u16()
		col.fieldID = dec.u32()
		col.flags = colFlag(dec.u32())
		log.Printf("col[%d]: %+v", i, col)
	}

	return cols
}

func (dec *decoder) decodeAliases() (cols []colAlias) {
	if dec.err != nil {
		return nil
	}

	size, n := dec.decodeFrameHeader()
	if dec.err != nil {
		return nil
	}
	log.Printf("aliases: size:%d, n=%d", size, n)

	cols = make([]colAlias, n)
	for i := range cols {
		csz, cnel := dec.decodeFrameHeader()
		log.Printf("alias[%d]: %d, %d", i, csz, cnel)
		col := &cols[i]
		col.physID = dec.u32()
		col.fieldID = dec.u32()
		log.Printf("alias[%d]: %+v", i, col)
	}

	return cols
}

func (dec *decoder) decodeExtraTypeInfo() (cols []colExtra) {
	if dec.err != nil {
		return nil
	}

	size, n := dec.decodeFrameHeader()
	if dec.err != nil {
		return nil
	}
	log.Printf("extra: size:%d, n=%d", size, n)

	cols = make([]colExtra, n)
	for i := range cols {
		csz, cnel := dec.decodeFrameHeader()
		log.Printf("extra[%d]: %d, %d", i, csz, cnel)
		col := &cols[i]
		col.typeFrom = dec.u32()
		col.typeTo = dec.u32()
		col.contentID = dec.u32()
		col.typeName = dec.decodeString()
		log.Printf("extra[%d]: %+v", i, col)
	}

	return cols
}

func (dec *decoder) decodeExtHeaders() []extHeader {
	if dec.err != nil {
		return nil
	}

	_, n := dec.decodeFrameHeader()
	log.Printf("ftr: nxhdr=%d", n)
	if n > 0 {
		dec.err = fmt.Errorf("rntup: extension headers are not supported")
		return nil
	}

	return nil
}

func (dec *decoder) decodeColGroups() []colGroup {
	if dec.err != nil {
		return nil
	}

	_, n := dec.decodeFrameHeader()
	log.Printf("ftr: col-groups=%d", n)
	if n > 0 {
		dec.err = fmt.Errorf("rntup: sharded clusters are not supported")
		return nil
	}

	return nil
}

func (dec *decoder) decodeClusterInfos() (cls []clusterInfo) {
	if dec.err != nil {
		return nil
	}

	size, n := dec.decodeFrameHeader()
	if dec.err != nil {
		return nil
	}
	log.Printf("ftr: size=%d cluster-summaries=%d", size, n)

	cls = make([]clusterInfo, n)
	for i := range cls {
		csz, cnel := dec.decodeFrameHeader()
		log.Printf("cluster[%d]: %d, %d", i, csz, cnel)
		clus := &cls[i]
		clus.firstEntry = dec.u64()
		nentries := int64(dec.u64())
		switch {
		case nentries < 0:
			clus.nentries = uint64(-nentries)
			clus.colGrpID = int32(dec.u32())
		default:
			clus.nentries = uint64(nentries)
			clus.colGrpID = -1
		}
		log.Printf("cluster[%d]: %+v", i, clus)
	}

	return cls
}

func (dec *decoder) decodeClusterGroups() (cls []clusterGroup) {
	if dec.err != nil {
		return nil
	}

	size, n := dec.decodeFrameHeader()
	if dec.err != nil {
		return nil
	}
	log.Printf("ftr: size=%d cluster-groups=%d", size, n)

	cls = make([]clusterGroup, n)
	for i := range cls {
		csz, cnel := dec.decodeFrameHeader()
		log.Printf("cluster[%d]: %d, %d", i, csz, cnel)
		clus := &cls[i]
		clus.n = dec.u32()
		clus.pages = dec.decodeEnvelopeLink()
		log.Printf("cluster[%d]: %+v", i, clus)
	}

	return cls
}

func (dec *decoder) decodeMetaDataBlocks() []metaDataBlock {
	if dec.err != nil {
		return nil
	}

	_, n := dec.decodeFrameHeader()
	log.Printf("ftr: n-md=%d", n)
	if n > 0 {
		dec.err = fmt.Errorf("rntup: metadata blocks are not supported")
		return nil
	}

	return nil
}

func (dec *decoder) decodeFrameHeader() (size, nelm uint32) {
	if dec.err != nil {
		return
	}

	size = dec.u32()
	sz := int32(size)
	switch {
	case sz >= 0:
		// record frame
		nelm = 1
		if size < 4 {
			dec.err = fmt.Errorf("rntup: corrupt record frame")
			return
		}
	default:
		// list frame
		nelm = dec.u32()
		nelm &= (2 << 28) - 1
		size = uint32(-sz)
		if size < 2*4 {
			dec.err = fmt.Errorf("rntup: corrupt list frame")
		}
	}

	return
}

func (dec *decoder) decodeEnvelopeLink() (lnk envelopeLink) {
	if dec.err != nil {
		return lnk
	}

	lnk.size = dec.u32()
	lnk.loc = dec.decodeLocator()

	return lnk
}

func (dec *decoder) decodeLocator() (loc locator) {
	if dec.err != nil {
		return loc
	}

	hdr := int32(dec.u32())
	switch {
	case hdr < 0:
		hdr = -hdr
		typ := uint8(hdr >> 24)
		if typ != 0x02 {
			dec.err = fmt.Errorf("rntup: invalid locator type=0x%x", typ)
			return loc
		}
		loc.pos = 0
		loc.storage = 0
		raw := make([]byte, uint32(hdr)&0x00FFFFFF)
		_, err := dec.read(raw)
		if err != nil {
			return loc
		}
		loc.url = string(raw)
	default:
		loc.pos = dec.u64()
		loc.storage = uint32(hdr)
		loc.url = ""
	}

	return loc
}

func (dec *decoder) decodePagelist() (ps pageList) {
	if dec.err != nil {
		return ps
	}

	ps.vers = dec.u16()
	ps.minv = dec.u16()

	size, n := dec.decodeFrameHeader()
	log.Printf("ps: cluster-descr: size=%d, n=%d", size, n)
	ps.clusters = make([]clusterDescr, n)
	for i := range ps.clusters {
		cluster := &ps.clusters[i]
		size, n := dec.decodeFrameHeader()
		log.Printf("ps: cluster[%d], cols: size=%d, n=%d", i, size, n)
		cluster.columns = make([]columnDescr, n)
		for j := range cluster.columns {
			size, n := dec.decodeFrameHeader()
			col := &cluster.columns[j]
			log.Printf("ps: cluster[%d].col[%d]: pages: size=%d, n=%d", i, j, size, n)
			col.pages = make([]pageDescr, n)
			for k := range col.pages {
				p := &col.pages[k]
				p.nelem = dec.u32()
				p.loc = dec.decodeLocator()
			}
			col.offset = dec.u64()
			col.compr = dec.u32()
		}
	}

	ps.crc32 = dec.u32()
	return ps
}
