package lru

import (
	"container/list"
)

// 定义LRU缓存，非线程安全
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nbytes    int64                         // 记录当前已使用的内存
	ll        *list.List                    // 双向链表，当访问到某个值时，将其移动到队尾的复杂度是O(1)，在队尾新增一条记录以及删除一条记录的复杂度均为O(1)
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 记录某条记录被移除时的回调函数，可以为nil
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 用于返回该Value所占用的内存大小
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找缓存: 1）字典中找到对应的双向链表的节点；2）将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, ok
	}
	return
}

// 缓存淘汰: 移除最近最少访问的节点，如果回调函数 OnEvicted 不为 nil，则调用回调函数
func (c *Cache) Eliminate() {
	if ele := c.ll.Back(); ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) bool {
	if c.maxBytes <= 0 {
		return false
	}
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.nbytes = c.nbytes + int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
		c.ll.MoveToFront(ele)
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.nbytes = c.nbytes + int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes < c.nbytes {
		c.Eliminate()
	}
	return true
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
