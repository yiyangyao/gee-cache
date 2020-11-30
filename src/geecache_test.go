package src

import (
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	// 借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
