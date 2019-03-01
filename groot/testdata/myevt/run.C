#include "myevt.h"

Event* makeEvent(int i) {
	Event *evt = new Event;
	Event &e = *evt;

	e.Beg = TString::Format("beg-%03d", i);
	e.I16 = i;
	e.I32 = i;
	e.I64 = i;
	e.U16 = i;
	e.U32 = i;
	e.U64 = i;
	e.F32 = float(i);
	e.F64 = double(i);
	e.Str = TString::Format("evt-%03d", i);
	e.P3.Px = i-1;
	e.P3.Py = double(i);
	e.P3.Pz = i-1;
	e.ObjStr = new TObjString(TString::Format("obj-%03d", i));

	for (int ii = 0; ii != ARRAYSZ; ii++) {
		e.ArrayI16[ii] = i;
		e.ArrayI32[ii] = i;
		e.ArrayI64[ii] = i;
		e.ArrayU16[ii] = i;
		e.ArrayU32[ii] = i;
		e.ArrayU64[ii] = i;
		e.ArrayF32[ii] = float(i);
		e.ArrayF64[ii] = double(i);
		e.ArrayP3s[ii] = {ii-1,double(ii),ii-1};
		e.ArrayObjStr[ii] = TObjString(TString::Format("obj-%03d", ii));
	}

	e.N = int32_t(i) % 10;
	e.SliceI16 = (int16_t*)malloc(sizeof(int16_t)*e.N);
	e.SliceI32 = (int32_t*)malloc(sizeof(int32_t)*e.N);
	e.SliceI64 = (int64_t*)malloc(sizeof(int64_t)*e.N);
	e.SliceU16 = (uint16_t*)malloc(sizeof(uint16_t)*e.N);
	e.SliceU32 = (uint32_t*)malloc(sizeof(uint32_t)*e.N);
	e.SliceU64 = (uint64_t*)malloc(sizeof(uint64_t)*e.N);
	e.SliceF32 = (float*)malloc(sizeof(float)*e.N);
	e.SliceF64 = (double*)malloc(sizeof(double)*e.N);
//	e.SliceP3s = (::P3*)malloc(sizeof(::P3)*e.N);
//	e.SliceStr = (TObjString*)malloc(sizeof(TObjString)*e.N);

	for (int ii = 0; ii != e.N; ii++) {
		e.SliceI16[ii] = i;
		e.SliceI32[ii] = i;
		e.SliceI64[ii] = i;
		e.SliceU16[ii] = i;
		e.SliceU32[ii] = i;
		e.SliceU64[ii] = i;
		e.SliceF32[ii] = float(i);
		e.SliceF64[ii] = double(i);
//		e.SliceP3s[ii] = {i-1,double(i),i-1};
//		e.SliceStr[ii] = TObjString(TString::Format("objstr-%03d", ii));
	}

	e.StdStr = std::string(TString::Format("std-%03d", i));
	e.StlVecI16.resize(e.N);
	e.StlVecI32.resize(e.N);
	e.StlVecI64.resize(e.N);
	e.StlVecU16.resize(e.N);
	e.StlVecU32.resize(e.N);
	e.StlVecU64.resize(e.N);
	e.StlVecF32.resize(e.N);
	e.StlVecF64.resize(e.N);
	e.StlVecStr.resize(e.N);
//	e.StlVecP3s.resize(e.N);
	for (int ii =0; ii != e.N; ii++) {
		e.StlVecI16[ii] = i;
		e.StlVecI32[ii] = i;
		e.StlVecI64[ii] = i;
		e.StlVecU16[ii] = i;
		e.StlVecU32[ii] = i;
		e.StlVecU64[ii] = i;
		e.StlVecF32[ii] = float(i);
		e.StlVecF64[ii] = double(i);
		e.StlVecStr[ii] = std::string(TString::Format("vec-%03d", i));
//		e.StlVecP3s[ii] = {i-1,double(i),i-1};
	}
	e.End = TString::Format("end-%03d", i);

	return evt;
}
void run() {
	gROOT->ProcessLine(".L myevt.h++");
	auto f = TFile::Open("o.root", "RECREATE");
	auto e = makeEvent(2);
	f->WriteObjectAny(e, "Event", "evt");
	f->Write();
	f->Close();
}
