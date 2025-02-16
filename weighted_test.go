/*
平滑加权轮询算法的实现
*/

package loadbalance

import (
	"fmt"
	"sync"
	"testing"
)

type Node struct {
	name       string
	weight     int
	currWeight int
}

func (n *Node) Invoke() {}

type Balancer struct {
	nodes []*Node
	mux   sync.Mutex
	t     *testing.T
}

func TestSmoothWRR(t *testing.T) {
	// 模拟单个服务的三个节点
	nodes := []*Node{
		{
			name:       "A",
			weight:     10,
			currWeight: 10,
		},
		{
			name:       "B",
			weight:     20,
			currWeight: 20,
		},
		{
			name:       "C",
			weight:     30,
			currWeight: 30,
		},
	}

	bl := &Balancer{
		nodes: nodes,
		t:     t,
	}

	for i := 1; i <= 6; i++ {
		t.Log(fmt.Sprintf("第 %d 个请求发送前,nodes %v", i, convert(nodes)))
		target := bl.pick()
		// 模拟 rpc 调用
		target.Invoke()
		t.Log(fmt.Sprintf("第 %d 个请求发送后,nodes %v", i, convert(nodes)))
	}
}

func (b *Balancer) pick() *Node {
	b.mux.Lock()
	defer b.mux.Unlock()
	total := 0
	// 计算总权重
	for _, n := range b.nodes {
		total += n.weight
	}
	// 计算当前权重
	for _, n := range b.nodes {
		n.currWeight = n.currWeight + n.weight
	}

	// 挑选节点
	var target *Node
	for _, n := range b.nodes {
		if target == nil {
			target = n
		} else {
			if target.currWeight < n.currWeight {
				target = n
			}
		}
	}

	b.t.Log("选中的节点的当前权重", target)
	target.currWeight = target.currWeight - total
	b.t.Log("选中的节点减去总权重后的权重", target)
	return target
}

// 简洁写法
//func (b *Balancer) pick() *Node {
//	b.mux.Lock()
//	defer b.mux.Unlock()
//	total := 0
//	// 计算总权重 当前权重 挑选节点
//	target := b.nodes[0]
//	for _, n := range b.nodes {
//		total += n.weight
//		n.currWeight = n.currWeight + n.weight
//		if target == nil || target.currWeight < n.currWeight {
//			target = n
//		}
//	}
//
//	target.currWeight = target.currWeight - total
//
//	b.t.Log("选中的节点的当前权重", target)
//	target.currWeight = target.currWeight - total
//	b.t.Log("选中的节点减去总权重后的权重", target)
//	return target
//}

func convert(src []*Node) []Node {
	dst := make([]Node, 0, len(src))
	for _, n := range src {
		dst = append(dst, Node{
			name:       n.name,
			weight:     n.weight,
			currWeight: n.currWeight,
		})
	}

	return dst
}
