package geecache

import "DisCache/consistenthash"

/*
 根据传入的key值去寻找相应的节点
*/

type PeerPicker interface {
	PickerPeer(key string) (peer PeerGetter, ok bool)
}

/**
去对应的group缓存中寻找缓存值
*/

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultreplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath} // 每个节点的网址
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _PeerPicker = (*HTTPPool)(nil)
