package src

/**
PickPeer() 方法用于根据传入的 key 选择相应节点
*/
type PeerPicker interface {
	PickPeer(key string) (p PeerGetter, ok bool)
}

/**
PeerGetter 的 Get() 方法用于从对应 group 查找缓存值
*/
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
