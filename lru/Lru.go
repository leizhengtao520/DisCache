package lru

import "container/list"

type Cache struct {
	maxBytes  int64                         //最大缓存大小
	nbytes    int64                         //当前缓存大小
	ll        *list.List                    //双向链表 维护lru
	cache     map[string]*list.Element      //缓存内容
	OnEvicted func(key string, value Value) //
}
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//cache   map[string]*list.Element

func (c *Cache) Len() int {
	return c.ll.Len()
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { //如果加入的数据已存在
		c.ll.MoveToFront(ele) //队尾
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value}) //队尾
		c.cache[key] = ele
		c.nbytes += int64(value.Len()) + int64(len(key))
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 头部元素
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry) //强转
		delete(c.cache, kv.key)
		c.nbytes -= int64(kv.value.Len()) + int64(len(kv.key))
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
