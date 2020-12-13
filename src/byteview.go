package src

// 抽象了一个只读数据结构 ByteView 用来表示缓存值

type ByteView struct {
	B []byte // b 将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等
}

// 实现 Len() 方法，实现 Value 接口
func (v ByteView) Len() int {
	return len(v.B)
}

// ByteSlice returns a copy of the data as a byte slice
func (v ByteView) ByteSlice() []byte {
	c := make([]byte, len(v.B))
	copy(c, v.B)
	return c
}

func (v ByteView) String() string {
	return string(v.B)
}
