package graph

import "sync"

// Hashset declare a type which is a map of string to bool
type Hashset map[string]bool

type Graph struct {
	Nodes map[string]Hashset
	mutex *sync.RWMutex
}

// NewGraph returns a new Graph with empty nodes
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]Hashset),
		mutex: new(sync.RWMutex),
	}
}

// AddOrReplaceNode adds or replaces a node in the graph
func (g *Graph) AddOrReplaceNode(k string, v Hashset) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Nodes[k] = v
}

// PruneNode removes a node from the graph. It also removes all edges
func (g *Graph) PruneNode(k string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	delete(g.Nodes, k)
	for _, v := range g.Nodes {
		delete(v, k)
	}
}

// PruneNodes removes an array of nodes from the graph. It also removes all edges.
func (g *Graph) PruneNodes(nodes []string) {
	for _, node := range nodes {
		g.PruneNode(node)
	}
}

// GetPruneCandidates returns an array of nodes that have no children
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
