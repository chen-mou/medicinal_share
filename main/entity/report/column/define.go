package column

import (
	"errors"
	"medicinal_share/main/middleware"
	"reflect"
)

type Column interface {
	Verify(any) error       //验证方法
	Builder(map[string]any) //根据数据库的define构建对象
	GetName() string        //获取当前列的名字
}

type BaseColumn struct {
	Name      string
	Type      string
	Value     any `json:"value,omitempty"`
	ValueType reflect.Kind
}

func (b BaseColumn) GetName() string {
	return b.Name
}

func (b BaseColumn) Verify(val any) error {
	t := reflect.TypeOf(val)
	if t.Kind() != b.ValueType {
		return errors.New("数据类型有误")
	}
	return nil
}

func (b *BaseColumn) Builder(data map[string]any) {
	defer func() {
		err := recover()
		if err != nil {
			panic(middleware.NewCustomErr(middleware.ERROR, "解析失败"))
		}
	}()
	b.ValueType = typeEquals(containOrPanic(data, "value_type"), reflect.Uint).(reflect.Kind)
}

type TextColumn struct {
	BaseColumn
}

func (TextColumn) Verify(value any) error {
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Map, reflect.Array, reflect.Struct, reflect.Slice:
		return errors.New("只能为基本数据类型")
	}
	return nil
}

type ListColumn struct {
	BaseColumn
	Values []string
}

func (l *ListColumn) Builder(m map[string]any) {
	l.BaseColumn.Builder(m)
	l.Values = typeEquals(containOrPanic(m, "values"), reflect.Slice).([]string)
	l.ValueType = reflect.Array
}

func (l ListColumn) Verify(value any) error {
	if err := l.BaseColumn.Verify(value); err != nil {
		return err
	}
	v := value.(string)
	for _, val := range l.Values {
		if v == val {
			return nil
		}
	}
	return errors.New("值不存在")
}

type RangeColumn struct {
	Max int
	Min int
	BaseColumn
}

func containOrPanic(m map[string]any, key string) any {
	if v, ok := m[key]; ok {
		return v
	}
	panic("键不存在")
}

func typeEquals(v any, typ reflect.Kind) any {
	t := reflect.TypeOf(v)
	if t.Kind() != typ {
		panic("类型有误")
	}
	return v
}
