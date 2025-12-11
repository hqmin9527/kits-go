package lru_cache

import (
	"container/list"
)

// IKey 使用Lru的数据需要实现的接口
// GetKey的结果是可以比较相等的
type IKey interface {
	Key() any
}

var defOnEvict = func(value any) {
	// default evict func do nothing
}

type LruCache struct {
	limit    int                   // 缓存容量
	evicts   *list.List            // 双向链表用于淘汰数据
	elements map[any]*list.Element // 记录缓存数据
	onEvict  func(value any)       // 缓存淘汰时的回调
}

func NewLruCache(limit int, onEvict func(value any)) *LruCache {
	if onEvict == nil {
		onEvict = defOnEvict
	}

	return &LruCache{
		limit:    limit,
		evicts:   list.New(),
		elements: make(map[any]*list.Element, limit),
		onEvict:  onEvict,
	}
}

// Add 添加元素到最新位置，如果缓存命中，也会更新值
func (l *LruCache) Add(values ...IKey) *LruCache {
	for _, v := range values {
		l.addSingle(v)
	}
	return l
}

func (l *LruCache) addSingle(value IKey) *LruCache {
	// 如果在缓存中，就将元素移到最前面
	if elm, ok := l.elements[value.Key()]; ok {
		l.evicts.MoveToFront(elm)
		// 如果缓存中存在，也要更新value
		elm.Value = value
		return l
	}
	// 添加新节点
	elm := l.evicts.PushFront(value)
	l.elements[value.Key()] = elm
	// 检查是否超出容量，如果超出就淘汰末尾节点数据
	if l.evicts.Len() > l.limit {
		l.removeOldest()
	}
	return l
}

// RangeFromLatest 正序遍历，从最近添加的开始访问
func (l *LruCache) RangeFromLatest(f func(value any) bool) {
	elm := l.evicts.Front()
	for elm != nil {
		if !f(elm.Value) {
			break
		}
		elm = elm.Next()
	}
}

// RangeFromEarliest 逆序遍历，从最早添加的开始访问
func (l *LruCache) RangeFromEarliest(f func(value any) bool) {
	elm := l.evicts.Back()
	for elm != nil {
		if !f(elm.Value) {
			break
		}
		elm = elm.Prev()
	}
}

// Evict 手动淘汰
func (l *LruCache) Evict() {
	l.removeOldest()
}

func (l *LruCache) Limit() int {
	return l.limit
}

func (l *LruCache) Len() int {
	return l.evicts.Len()
}

// 淘汰末尾节点
func (l *LruCache) removeOldest() {
	elem := l.evicts.Back() // 获取链表末尾节点
	if elem != nil {
		l.removeElement(elem)
	}
}

// 删除节点操作
func (l *LruCache) removeElement(e *list.Element) {
	l.evicts.Remove(e)
	key := e.Value.(IKey).Key()
	delete(l.elements, key)
	l.onEvict(e.Value)
}
