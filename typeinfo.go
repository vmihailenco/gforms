package gforms

import (
	"reflect"
	"strings"
	"sync"
)

var (
	fieldType = reflect.TypeOf((*Field)(nil)).Elem()
	tinfoMap  = newTypeInfoMap()
)

//------------------------------------------------------------------------------

type fieldFlags int

const (
	fReq fieldFlags = 1 << iota
)

type fieldInfo struct {
	idx    []int
	name   string
	constr constructor
	flags  fieldFlags
}

type typeInfo struct {
	fields []*fieldInfo
}

type typeInfoMap struct {
	l sync.RWMutex
	m map[reflect.Type]*typeInfo
}

func newTypeInfoMap() *typeInfoMap {
	return &typeInfoMap{
		m: make(map[reflect.Type]*typeInfo),
	}
}

func (m *typeInfoMap) TypeInfo(typ reflect.Type) *typeInfo {
	m.l.RLock()
	tinfo, ok := tinfoMap.m[typ]
	m.l.RUnlock()
	if ok {
		return tinfo
	}

	tinfo = &typeInfo{}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue // Private field
		}

		if !f.Type.Implements(fieldType) {
			continue
		}

		tinfo.fields = append(tinfo.fields, m.structFieldInfo(typ, &f))
	}

	tinfoMap.m[typ] = tinfo
	return tinfo
}

func (m *typeInfoMap) structFieldInfo(typ reflect.Type, f *reflect.StructField) *fieldInfo {
	finfo := &fieldInfo{
		idx:    f.Index,
		constr: tconstrMap.Constructor(f.Type),
	}

	tokens := strings.Split(f.Tag.Get("gforms"), ",")
	finfo.name = tokens[0]
	if len(tokens) > 1 {
		for _, flag := range tokens[1:] {
			switch flag {
			case "req":
				finfo.flags |= fReq
			}
		}
	}

	if finfo.name == "" {
		finfo.name = f.Name
	}
	finfo.name = strings.Join(SplitWords(finfo.name), " ")

	return finfo
}
