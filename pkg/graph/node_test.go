package graph

import (
	"reflect"
	"sort"
	"strconv"
	"sync"
	"testing"
)

func TestNewGraph(t *testing.T) {
	tests := []struct {
		name string
		want *Graph
	}{
		{
			name: "Graph with empty nodes",
			want: &Graph{
				Nodes: make(map[string]Hashset),
				mutex: new(sync.RWMutex),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGraph(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraph_AddOrReplaceNode(t *testing.T) {
	t.Run("Key does not exist and gets added", func(t *testing.T) {
		g := &Graph{
			Nodes: make(map[string]Hashset),
			mutex: new(sync.RWMutex),
		}
		g.AddOrReplaceNode("a", Hashset{"b": true})
		if g.Nodes["a"]["b"] != true {
			t.Errorf("AddOrReplaceNode() = %v, want %v", g.Nodes["a"]["b"], true)
		}
	})
	t.Run("Key already exists and value gets replaced", func(t *testing.T) {
		g := &Graph{
			Nodes: make(map[string]Hashset),
			mutex: new(sync.RWMutex),
		}
		g.AddOrReplaceNode("a", Hashset{"b": true})
		g.AddOrReplaceNode("a", Hashset{"c": true})
		if g.Nodes["a"]["c"] != true {
			t.Errorf("AddOrReplaceNode() = %v, want %v", g.Nodes["a"]["c"], true)
		}
		if g.Nodes["a"]["b"] != false {
			t.Errorf("AddOrReplaceNode() = %v, want %v", g.Nodes["a"]["b"], false)
		}
	})
}

func TestGraph_PruneNode(t *testing.T) {
	t.Run("Node already exists and gets removed", func(t *testing.T) {
		g := &Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true},
			},
			mutex: new(sync.RWMutex),
		}
		g.PruneNode("a")
		if _, ok := g.Nodes["a"]; ok {
			t.Errorf("PruneNode() expected to remove the Node but it still exists")
		}
	})
	t.Run("Node does not exist and continues to not exist", func(t *testing.T) {
		g := &Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true},
			},
			mutex: new(sync.RWMutex),
		}
		g.PruneNode("random")
		if _, ok := g.Nodes["random"]; ok {
			t.Errorf("PruneNode() expected to remove the Node but it exists")
		}
	})
	t.Run("Node exists and there are edges to the node which are removed too", func(t *testing.T) {
		g := &Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true, "d": true},
				"b": {},
			},
			mutex: new(sync.RWMutex),
		}
		g.PruneNode("b")
		if _, ok := g.Nodes["b"]; ok {
			t.Errorf("PruneNode() expected to remove the Node but it still exists")
		}
		if _, ok := g.Nodes["a"]["b"]; ok {
			t.Errorf("PruneNode() expected to remove the edge to Node but it still exists")
		}
		if _, ok := g.Nodes["a"]["d"]; !ok {
			t.Errorf("PruneNode() expected to not remove the edge to other nodes but removed")
		}
	})
}

func TestGraph_GetPruneCandidates(t *testing.T) {
	t.Run("Pruning candidates exist and are returned", func(t *testing.T) {
		graph, pruneCandidates := createHugeGraphAndPruneCandidates()
		expected := pruneCandidates
		sort.Strings(expected)
		got := graph.GetPruneCandidates()
		sort.Strings(got)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("GetPruneCandidates() = %v, want %v", got, expected)
		}
	})

	t.Run("Pruning candidates don't exist so empty slice is returned", func(t *testing.T) {
		g := Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true, "c": true},
				"b": {"c": true, "d": true},
			},
			mutex: new(sync.RWMutex),
		}
		var expected []string
		if got := g.GetPruneCandidates(); !reflect.DeepEqual(got, expected) {
			t.Errorf("GetPruneCandidates() = %v, want %v", got, expected)
		}
	})
}

func TestGraph_PruneNodes(t *testing.T) {
	t.Run("Nodes exists and there are edges to the node which are removed too", func(t *testing.T) {
		graph, pruneCandidates := createHugeGraphAndPruneCandidates()
		graph.PruneNodes(pruneCandidates)
		for _, node := range pruneCandidates {
			if _, ok := graph.Nodes[node]; ok {
				t.Errorf("PruneNodes() expected to remove the Node but it still exists")
			}
			for _, v := range graph.Nodes {
				if _, ok := v[node]; ok {
					t.Errorf("PruneNodes() expected to remove the edge to Node but it still exists")
				}
			}
		}
	})
	t.Run("Nodes don't exist and there are no edges to the node", func(t *testing.T) {
		graph, _ := createHugeGraphAndPruneCandidates()
		originalNodesSize := len(graph.Nodes)
		pruneNodes := []string{"random1", "random2"}
		graph.PruneNodes(pruneNodes)
		if len(graph.Nodes) != originalNodesSize {
			t.Errorf("PruneNodes() expected to not remove any nodes but removed")
		}
	})
}

// ----- Below set of benchmarks are for small graphs-----//
func BenchmarkGraph_AddOrReplaceNode(b *testing.B) {
	g := createBenchmarkGraph()
	for i := 0; i < b.N; i++ {
		g.AddOrReplaceNode("a", Hashset{"b": true})
	}
}

func BenchmarkGraph_PruneNode(b *testing.B) {
	g := createBenchmarkGraph()
	for i := 0; i < b.N; i++ {
		g.PruneNode("empty1")
	}
}

func BenchmarkGraph_GetPruneCandidates(b *testing.B) {
	g := createBenchmarkGraph()
	for i := 0; i < b.N; i++ {
		g.GetPruneCandidates()
	}
}
func createBenchmarkGraph() Graph {
	return Graph{
		Nodes: map[string]Hashset{
			"a":       {"b": true, "c": true, "empty1": true, "empty2": true},
			"b":       {"c": true, "d": true, "empty1": true, "empty2": true},
			"c":       {"d": true, "e": true, "empty1": true, "empty2": true},
			"d":       {"e": true, "f": true, "empty1": true, "empty2": true},
			"e":       {"f": true, "g": true, "empty1": true, "empty2": true},
			"f":       {"g": true, "h": true, "empty1": true, "empty2": true},
			"g":       {"h": true, "i": true, "empty1": true, "empty2": true},
			"h":       {"i": true, "j": true, "empty1": true, "empty2": true},
			"i":       {"j": true, "k": true, "empty1": true, "empty2": true},
			"j":       {"k": true, "l": true, "empty1": true, "empty2": true},
			"k":       {"l": true, "m": true, "empty1": true, "empty2": true},
			"l":       {"m": true, "n": true, "empty1": true, "empty2": true},
			"m":       {"n": true, "o": true, "empty1": true, "empty2": true},
			"n":       {"o": true, "p": true, "empty1": true, "empty2": true},
			"o":       {"p": true, "q": true, "empty1": true, "empty2": true},
			"empty1":  {},
			"empty2":  {},
			"empty3":  {},
			"empty4":  {},
			"empty5":  {},
			"empty6":  {},
			"empty7":  {},
			"empty8":  {},
			"empty9":  {},
			"empty10": {},
		},
		mutex: new(sync.RWMutex),
	}
}

//------ Below set of benchmarks are for huge graphs------//

// Since the benchmark function will account for graph creation
// we have created a benchmark function for the huge
func BenchmarkGraph_CreateHugeGraphAndPruneCandidates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createHugeGraphAndPruneCandidates()
	}
}

func BenchmarkGraph_PruneNodeForHugeGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, p := createHugeGraphAndPruneCandidates()
		g.PruneNode(p[0])
	}
}

func BenchmarkGraph_PruneNodesForHugeGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, p := createHugeGraphAndPruneCandidates()
		g.PruneNodes(p)
	}
}

func BenchmarkGraph_GetPruneCandidatesForHugeGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, _ := createHugeGraphAndPruneCandidates()
		g.GetPruneCandidates()
	}
}

func BenchmarkGraph_AddOrReplaceNodeForHugeGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g, p := createHugeGraphAndPruneCandidates()
		g.AddOrReplaceNode(p[0], Hashset{"random": true})
	}
}

func createHugeGraphAndPruneCandidates() (Graph, []string) {
	m := make(map[string]Hashset, 10000)
	pruneCandidates := make([]string, 500)
	for i := 0; i < 5000; i++ {
		h := make(Hashset, 1000)
		for j := 5000; j < 6000; j++ {
			h[strconv.Itoa(j)] = true
		}
		m[strconv.Itoa(i)] = h
	}
	for i := 5000; i < 5500; i++ {
		m[strconv.Itoa(i)] = Hashset{}
		pruneCandidates[i-5000] = strconv.Itoa(i)
	}
	g := Graph{
		Nodes: m,
		mutex: new(sync.RWMutex),
	}
	return g, pruneCandidates
}
