package gee

import "strings"

/**
前缀树处理路由(动态路由)
*/
type node struct {
	pattern  string  // 完全匹配后的路径
	part     string  // 匹配的一部分
	children []*node // 子节点
	isWild   bool    // 是否是模糊匹配 true开启模糊查找
}

/**
插入元素形成前缀树
*/
func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) { // 判断到了最后一个元素赋值直接返回
		n.pattern = pattern
		return
	}
	part := parts[height]

	child := n.matchChild(part)
	if child == nil { // 不存在则插入
		child = &node{
			part:   part,
			isWild: part[0] == '*' || part[0] == ':',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

/**
搜索符合条件的node
*/
func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

/**
匹配所有子元素
*/
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.isWild || child.part == part {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

/**
查询子元素内有没有匹配的
*/
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}

	return nil
}
