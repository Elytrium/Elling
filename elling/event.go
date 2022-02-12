package elling

import (
	"reflect"
)

var Dispatchers map[reflect.Type][]Dispatcher

type Dispatcher struct {
	Method reflect.Value
}

func (d *Dispatcher) Dispatch(event interface{}) {
	d.Method.Call([]reflect.Value{reflect.ValueOf(event)})
}

func RegisterListener(d interface{}) {
	dispatcherType := reflect.ValueOf(d)
	dispatcherMethods := dispatcherType.NumMethod()

	for i := 0; i < dispatcherMethods; i++ {
		method := dispatcherType.Method(i)
		methodType := method.Type()

		if methodType.NumIn() == 1 {
			Dispatchers[methodType.In(0)] = append(Dispatchers[methodType.In(0)], Dispatcher{Method: method})
		}
	}
}

func DispatchEvent(event interface{}) {
	for _, dispatcher := range Dispatchers[reflect.TypeOf(event)] {
		dispatcher.Dispatch(event)
	}
}
