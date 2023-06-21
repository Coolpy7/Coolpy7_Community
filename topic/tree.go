package topic

import (
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	children map[string]*Node
	values   []interface{}
}

func newNode() *Node {
	return &Node{
		children: make(map[string]*Node),
	}
}

func (n *Node) removeValue(value interface{}) {
	for i, v := range n.values {
		if v == value {
			n.values[i] = n.values[len(n.values)-1]
			n.values[len(n.values)-1] = nil
			n.values = n.values[:len(n.values)-1]
			break
		}
	}
}

func (n *Node) clearValues() {
	n.values = []interface{}{}
}

func (n *Node) string(level int) string {
	str := ""
	if level != 0 {
		str = fmt.Sprintf("%d", len(n.values))
	}

	for key, node := range n.children {
		str += fmt.Sprintf("\n| %s'%s' => %s", strings.Repeat(" ", level*2), key, node.string(level+1))
	}
	return str
}

type Tree struct {
	separator    string
	wildcardOne  string
	wildcardSome string
	root         *Node
	mutex        sync.RWMutex
}

func NewTree(separator, wildcardOne, wildcardSome string) *Tree {
	return &Tree{
		separator:    separator,
		wildcardOne:  wildcardOne,
		wildcardSome: wildcardSome,
		root:         newNode(),
	}
}

func NewStandardTree() *Tree {
	return NewTree("/", "+", "#")
}

func (t *Tree) Add(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.add(value, topic, t.root)
}

func (t *Tree) add(value interface{}, topic string, node *Node) {
	if topic == topicEnd {
		for _, v := range node.values {
			if v == value {
				return
			}
		}
		node.values = append(node.values, value)
		return
	}
	segment := topicSegment(topic, t.separator)
	child, ok := node.children[segment]
	if !ok {
		child = newNode()
		node.children[segment] = child
	}
	t.add(value, topicShorten(topic, t.separator), child)
}

func (t *Tree) Set(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.set(value, topic, t.root)
}

func (t *Tree) set(value interface{}, topic string, node *Node) {
	if topic == topicEnd {
		node.values = []interface{}{value}
		return
	}
	segment := topicSegment(topic, t.separator)
	child, ok := node.children[segment]
	if !ok {
		child = newNode()
		node.children[segment] = child
	}
	t.set(value, topicShorten(topic, t.separator), child)
}

func (t *Tree) Get(topic string) []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.get(topic, t.root)
}

func (t *Tree) get(topic string, node *Node) []interface{} {
	if topic == topicEnd {
		return node.values
	}
	segment := topicSegment(topic, t.separator)
	child, ok := node.children[segment]
	if !ok {
		return nil
	}
	return t.get(topicShorten(topic, t.separator), child)
}

func (t *Tree) Remove(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.remove(value, topic, t.root)
}

func (t *Tree) Empty(topic string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.remove(nil, topic, t.root)
}

func (t *Tree) remove(value interface{}, topic string, node *Node) bool {
	if topic == topicEnd {
		if value == nil {
			node.clearValues()
		} else {
			node.removeValue(value)
		}
		return len(node.values) == 0 && len(node.children) == 0
	}
	segment := topicSegment(topic, t.separator)
	child, ok := node.children[segment]
	if !ok {
		return false
	}
	if t.remove(value, topicShorten(topic, t.separator), child) {
		delete(node.children, segment)
	}
	return len(node.values) == 0 && len(node.children) == 0
}

func (t *Tree) Clear(value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.clear(value, t.root)
}

func (t *Tree) clear(value interface{}, node *Node) bool {
	node.removeValue(value)
	for segment, child := range node.children {
		if t.clear(value, child) {
			delete(node.children, segment)
		}
	}

	return len(node.values) == 0 && len(node.children) == 0
}

func (t *Tree) Match(topic string) []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	var list []interface{}
	t.match(topic, t.root, func(values []interface{}) bool {
		list = append(list, values...)
		return true
	})
	return t.clean(list)
}

func (t *Tree) MatchFirst(topic string) interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	var value interface{}
	t.match(topic, t.root, func(values []interface{}) bool {
		value = values[0]
		return false
	})
	return value
}

func (t *Tree) match(topic string, node *Node, fn func([]interface{}) bool) {
	if child, ok := node.children[t.wildcardSome]; ok && len(child.values) > 0 {
		if !fn(child.values) {
			return
		}
	}
	if topic == topicEnd {
		if len(node.values) > 0 {
			fn(node.values)
		}
		return
	}
	if child, ok := node.children[t.wildcardOne]; ok {
		t.match(topicShorten(topic, t.separator), child, fn)
	}
	segment := topicSegment(topic, t.separator)
	if segment != t.wildcardOne && segment != t.wildcardSome {
		if child, ok := node.children[segment]; ok {
			t.match(topicShorten(topic, t.separator), child, fn)
		}
	}
}

func (t *Tree) Search(topic string) []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	var list []interface{}
	t.search(topic, t.root, func(values []interface{}) bool {
		list = append(list, values...)
		return true
	})
	return t.clean(list)
}

func (t *Tree) SearchFirst(topic string) interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	var value interface{}
	t.search(topic, t.root, func(values []interface{}) bool {
		value = values[0]
		return false
	})
	return value
}

func (t *Tree) search(topic string, node *Node, fn func([]interface{}) bool) {
	if topic == topicEnd {
		if len(node.values) > 0 {
			fn(node.values)
		}
		return
	}
	segment := topicSegment(topic, t.separator)
	if segment == t.wildcardSome {
		if len(node.values) > 0 {
			if !fn(node.values) {
				return
			}
		}
		for _, child := range node.children {
			t.search(topic, child, fn)
		}
	}
	if segment == t.wildcardOne {
		if len(node.values) > 0 {
			if !fn(node.values) {
				return
			}
		}

		for _, child := range node.children {
			t.search(topicShorten(topic, t.separator), child, fn)
		}
	}
	if segment != t.wildcardOne && segment != t.wildcardSome {
		if child, ok := node.children[segment]; ok {
			t.search(topicShorten(topic, t.separator), child, fn)
		}
	}
}

func (t *Tree) clean(values []interface{}) []interface{} {
	seen := make(map[interface{}]struct{}, len(values))
	j := 0
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		values[j] = v
		j++
	}
	seen = nil
	return values[:j]
}

func (t *Tree) Count() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.count(t.root)
}

func (t *Tree) count(node *Node) int {
	total := 0
	for _, child := range node.children {
		total += t.count(child)
	}
	return total + len(node.values)
}

func (t *Tree) All() []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.clean(t.all([]interface{}{}, t.root))
}

func (t *Tree) all(result []interface{}, node *Node) []interface{} {
	for _, child := range node.children {
		result = t.all(result, child)
	}
	return append(result, node.values...)
}

func (t *Tree) Reset() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.root = newNode()
}

func (t *Tree) String() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return fmt.Sprintf("topic.Tree:%s", t.root.string(0))
}

var topicEnd = "\x00"

func topicShorten(topic, separator string) string {
	i := strings.Index(topic, separator)
	if i >= 0 {
		return topic[i+1:]
	}
	return topicEnd
}

func topicSegment(topic, separator string) string {
	i := strings.Index(topic, separator)
	if i >= 0 {
		return topic[:i]
	}
	return topic
}
