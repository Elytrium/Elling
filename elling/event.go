package elling

import (
	"reflect"
)

var ModuleDispatchers = make(map[reflect.Type][]*Dispatcher)

var LocalDispatchers = make(map[reflect.Type][]*Dispatcher)

type Dispatcher struct {
	Method reflect.Value
}

func (d *Dispatcher) Dispatch(event interface{}) {
	d.Method.Call([]reflect.Value{reflect.ValueOf(event)})
}

func RegisterListener(d interface{}) {
	registerListener(d, ModuleDispatchers)
}

func registerLocalListener(d interface{}) {
	registerListener(d, LocalDispatchers)
}

func registerListener(d interface{}, storage map[reflect.Type][]*Dispatcher) {
	dispatcherType := reflect.ValueOf(d)
	dispatcherMethods := dispatcherType.NumMethod()

	for i := 0; i < dispatcherMethods; i++ {
		method := dispatcherType.Method(i)
		methodType := method.Type()

		if methodType.NumIn() == 1 {
			storage[methodType.In(0)] = append(storage[methodType.In(0)], &Dispatcher{Method: method})
		}
	}
}

func DispatchEvent(event interface{}) {
	for _, dispatcher := range LocalDispatchers[reflect.TypeOf(event)] {
		dispatcher.Dispatch(event)
	}

	for _, dispatcher := range ModuleDispatchers[reflect.TypeOf(event)] {
		dispatcher.Dispatch(event)
	}
}
