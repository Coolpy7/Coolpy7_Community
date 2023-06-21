// Copyright 2022 LiDonghai Email:5241871@qq.com Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multimap

import (
	"container/list"
	"github.com/Coolpy7/Coolpy7_Community/packet"
	"sync"
)

type MultiMap struct {
	data map[string]*list.List
	sync.RWMutex
}

func NewMultiMap() *MultiMap {
	mm := &MultiMap{
		data: make(map[string]*list.List),
	}
	return mm
}

func (m *MultiMap) Push(mk string, pkt packet.Generic) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.data[mk]
	if !ok {
		v = &list.List{}
		m.data[mk] = v
	}
	v.PushBack(pkt)
}

func (m *MultiMap) Pop(mk string) (*list.Element, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.data[mk]
	if ok {
		vv := v.Front()
		if vv != nil {
			return vv, true
		}
	}
	return nil, false
}

func (m *MultiMap) GetBody(mk string) (*list.List, bool) {
	m.RLock()
	defer m.RUnlock()
	v, ok := m.data[mk]
	return v, ok
}

func (m *MultiMap) Del(mk string, item *list.Element) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.data[mk]
	if ok {
		v.Remove(item)
	}
}

func (m *MultiMap) Clear(mk string) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, mk)
}

func (m *MultiMap) Reset() {
	m.Lock()
	defer m.Unlock()
	m.data = make(map[string]*list.List)
}
