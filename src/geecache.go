package src

import (
	"fmt"
	"gee-cache/src/byteview"
	"gee-cache/src/cache"
	"log"
	"sync"
)

// 负责与用户的交互，并且控制缓存值存储和获取的流程，接口型函数

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/**
一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name
*/
type Group struct {
	name      string
	getter    Getter
	mainCache cache.Cache
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.Cache{CacheBytes: cacheBytes},
	}
	groups[name] = g

	return g
}

func GetGroup(name string) *Group {
	mutex.RLock()
	g := groups[name]
	mutex.RUnlock()
	return g
}

func (g *Group) Get(key string) (byteview.ByteView, error) {
	if key == "" {
		return byteview.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value byteview.ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (byteview.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return byteview.ByteView{}, err
	}
	newBytes := make([]byte, len(bytes))
	copy(newBytes, bytes)
	value := byteview.ByteView{B: newBytes}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value byteview.ByteView) {
	g.mainCache.Add(key, value)
}
