package dnsovertor

// 前缀树用于匹配DOMAIN-SUFFIX域名后缀
// TrieNode 定义前缀树的节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Trie 定义前缀树
type Trie struct {
	root *TrieNode
}

// NewTrieNode 创建一个新的前缀树节点
func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

// NewTrie 创建一个新的前缀树
func NewTrie() *Trie {
	return &Trie{root: NewTrieNode()}
}

// Insert 向前缀树中插入一个字符串
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// Search 搜索前缀树中的字符串
func (t *Trie) Search(word string) bool {
	node := t.root
	for _, ch := range word {
		if _, ok := node.children[ch]; !ok {
			return false
		}
		node = node.children[ch]
	}
	return node.isEnd
}
