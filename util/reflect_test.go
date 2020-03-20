package util

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tianhongw/misc-go/util/assert"
)

func TestReflect(t *testing.T) {
	var f func() int
	{
		ReplaceFuncVar(&f, func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(30)}
		})

		result := f()
		assert.Equal(t, 30, result)
	}

	{
		fv := reflect.ValueOf(&f)
		ReplaceFuncVar(fv, func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(31)}
		})

		result := f()
		assert.Equal(t, 31, result)
	}

	s := &struct {
		F func() int
		I int
	}{}

	{
		ReplaceFuncVar(&s.F, func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(40)}
		})

		result := s.F()
		assert.Equal(t, 40, result)
	}

	{
		sv := reflect.ValueOf(s)
		ReplaceFuncVar(reflect.Indirect(sv).Field(0), func([]reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(41)}
		})

		result := s.F()
		assert.Equal(t, 41, result)
	}

	{
		sfields := StructFieldValues(s, func(_ string, field reflect.Value) bool {
			return field.Kind() == reflect.Int
		})
		assert.Equal(t, 1, len(sfields))
		sfields["I"].Set(reflect.ValueOf(20))
		assert.Equal(t, 20, s.I)
	}

	{
		vf := Func2Value(f)
		ret := vf.Call(nil)
		assert.Equal(t, 31, ret[0].Interface().(int))
	}

	{
		inTypes := FuncInputTypes(testTarget)
		assert.Equal(t, 2, len(inTypes))
		assert.Equal(t, reflect.Int, inTypes[0].Kind())
		assert.Equal(t, reflect.String, inTypes[1].Kind())

		outTypes := FuncOutputTypes(testTarget)
		assert.Equal(t, 1, len(outTypes))
		assert.Equal(t, reflect.Slice, outTypes[0].Kind())
	}

	{
		stringType := TypeByPointer((*string)(nil))
		assert.Equal(t, reflect.ValueOf("").Type(), stringType)

		is := InstanceByType(stringType)
		_, ok := is.(string)
		assert.Equal(t, true, ok)
		isPtr := InstancePtrByType(stringType)
		_, ok = isPtr.(*string)
		assert.Equal(t, true, ok)
	}

	var t2 TestType
	{
		methods := ScanMethods(t2)
		_, ok := methods["M1"]
		assert.Equal(t, true, ok)
		assert.Equal(t, 1, len(methods))

		methods = ScanMethods(&t2)
		_, ok = methods["M2"]
		//assert.Assert(t, len(methods) == 3 && ok, "%v %v", len(methods), ok)
		assert.Equal(t, 3, len(methods))
		assert.Equal(t, true, ok)
	}

	{
		s := "abc"
		sv := reflect.ValueOf(s)
		sptr := InstancePtrByClone(sv)
		sp, ok := sptr.(*string)
		assert.Equal(t, true, ok)
		assert.Equal(t, "abc", *sp)
		*sp = "def"
		sp, ok = sptr.(*string)
		assert.Equal(t, true, ok)
		assert.Equal(t, "abc", s)
		assert.Equal(t, "def", *sp)
	}

	{
		var itf interface{}
		itf = t2
		_, ok := itf.(interface{ M2() })
		assert.Equal(t, true, !ok)
		itf = &itf
		_, ok = itf.(interface{ M2() })
		assert.Equal(t, true, !ok)
		itf = &t2
		_, ok = itf.(interface{ M2() })
		assert.Equal(t, true, ok)
	}

	{
		a := JSON{A: 1}
		b := &JSON{A: 1}
		bytes1, _ := json.Marshal(a)
		bytes2, _ := json.Marshal(b)
		assert.Equal(t, string(bytes1), string(bytes2))
	}

}

type JSON struct {
	A int
}

func testTarget(int, string) []int {
	return nil
}

type TestType struct {
	Base
}

type Base struct {
}

func (t TestType) M1() {

}

func (t TestType) m1() {

}

func (t *TestType) M2() {
}

func (b *Base) OK() {

}
