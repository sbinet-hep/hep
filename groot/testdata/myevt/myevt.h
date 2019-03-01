#ifndef MYEVT_H
#define MYEVT_H 1

#include <vector>
#include <string>

#include "TObjString.h"
#include "TString.h"

const int ARRAYSZ = 10;

struct P3 {
	int32_t Px;
	double  Py;
	int32_t Pz;
};

struct Event {
	TString  Beg;
	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	float    F32;
	double   F64;
	TString  Str;
	::P3      P3;
	TObjString *ObjStr;

	int16_t  ArrayI16[ARRAYSZ];
	int32_t  ArrayI32[ARRAYSZ];
	int64_t  ArrayI64[ARRAYSZ];
	uint16_t ArrayU16[ARRAYSZ];
	uint32_t ArrayU32[ARRAYSZ];
	uint64_t ArrayU64[ARRAYSZ];
	float    ArrayF32[ARRAYSZ];
	double   ArrayF64[ARRAYSZ];
	::P3     ArrayP3s[ARRAYSZ];
	TObjString ArrayObjStr[ARRAYSZ];

	int32_t  N;
	int16_t  *SliceI16;  //[N]
	int32_t  *SliceI32;  //[N]
	int64_t  *SliceI64;  //[N]
	uint16_t *SliceU16;  //[N]
	uint32_t *SliceU32;  //[N]
	uint64_t *SliceU64;  //[N]
	float    *SliceF32;  //[N]
	double   *SliceF64;  //[N]
//	::P3     *SliceP3s;  //[N]
//	TObjString *SliceStr;  //[N]

	std::string StdStr;

	std::vector<int16_t> StlVecI16;
	std::vector<int32_t> StlVecI32;
	std::vector<int64_t> StlVecI64;
	std::vector<uint16_t> StlVecU16;
	std::vector<uint32_t> StlVecU32;
	std::vector<uint64_t> StlVecU64;
	std::vector<float> StlVecF32;
	std::vector<double> StlVecF64;
	std::vector<std::string> StlVecStr;
//	std::vector< ::P3 > StlVecP3s;

	TString End;
};

#endif // MYEVT_H
