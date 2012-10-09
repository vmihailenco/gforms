package gforms

import (
	"reflect"
	"sync"
)

var (
	tconstrMap = newTypeConstructorMap()
)

type constructor func() interface{}

type typeConstructorMap struct {
	l sync.RWMutex
	m map[reflect.Type]constructor
}

func newTypeConstructorMap() *typeConstructorMap {
	return &typeConstructorMap{
		m: make(map[reflect.Type]constructor),
	}
}

func (m *typeConstructorMap) Register(field Field, constr constructor) {
	typ := reflect.ValueOf(field).Type()
	m.l.Lock()
	m.m[typ] = constr
	m.l.Unlock()
}

func (m *typeConstructorMap) Constructor(typ reflect.Type) constructor {
	m.l.RLock()
	constr := m.m[typ]
	m.l.RUnlock()
	return constr
}

//------------------------------------------------------------------------------

func Register(field Field, constr constructor) {
	tconstrMap.Register(field, constr)
}
