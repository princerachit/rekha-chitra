package main

import "sync"

// Hashset declare a type which is a map of string to bool
type Hashset map[string]bool

type Graph struct {
	Nodes map[string]Hashset
	mutex *sync.RWMutex
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]Hashset),
		mutex: new(sync.RWMutex),
	}
}

func (g *Graph) AddOrReplaceNode(k string, v Hashset) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Nodes[k] = v
}

func (g *Graph) RemoveNode(k string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.Nodes, k)
	for _, v := range g.Nodes {
		delete(v, k)
	}
}

func (g *Graph) GetPruneCandidates() []string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	var pruneCandidates []string
	for k, v := range g.Nodes {
		if len(v) == 0 {
			pruneCandidates = append(pruneCandidates, k)
		}
	}
	return pruneCandidates
}
