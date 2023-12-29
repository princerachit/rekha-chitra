package graph

import "sync"

// Dependents declare a type which is a map of string to bool
// This is used to represent the dependents of a node
type Dependents map[string]bool

// Graph is a struct that represents a graph which is a map of string to Dependents
// It uses a sync.RWMutex to make all operations thread safe by acquiring a lock
type Graph struct {
	Nodes map[string]Dependents
	mutex *sync.RWMutex
}

// NewGraph returns a new Graph with empty nodes
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]Dependents),
		mutex: new(sync.RWMutex),
	}
}

// AddOrReplaceNode adds or replaces a node in the graph
func (g *Graph) AddOrReplaceNode(k string, v Dependents) {
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
// It internally calls PruneNode for each node in the array
func (g *Graph) PruneNodes(nodes []string) {
	for _, node := range nodes {
		g.PruneNode(node)
	}
}

// GetPruneCandidates returns an array of nodes that have no children and can be pruned from the graph safely
// It does not modify the graph, but it acquires a read lock on entire graph which blocks all other operations
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
