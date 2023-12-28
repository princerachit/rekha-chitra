package main

import (
	"reflect"
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
			Nodes: make(map[string]Hashset),
			mutex: new(sync.RWMutex),
		}
		g.PruneNode("a")
		if _, ok := g.Nodes["a"]; ok {
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
		g := createTestGraph()
		expected := []string{"c", "d"}
		if got := g.GetPruneCandidates(); !reflect.DeepEqual(got, expected) {
			t.Errorf("GetPruneCandidates() = %v, want %v", got, expected)
		}
	})

	t.Run("Pruning candidates don't exist so empty slice is returned", func(t *testing.T) {
		g := createTestGraph()
		delete(g.Nodes, "c")
		delete(g.Nodes, "d")
		var expected []string
		if got := g.GetPruneCandidates(); !reflect.DeepEqual(got, expected) {
			t.Errorf("GetPruneCandidates() = %v, want %v", got, expected)
		}
	})
}

func TestGraph_PruneNodes(t *testing.T) {
	t.Run("Nodes exists and there are edges to the node which are removed too", func(t *testing.T) {
		g := &Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true, "d": true, "e": true},
				"c": {"b": true, "d": true, "f": true},
				"b": {},
				"d": {},
			},
			mutex: new(sync.RWMutex),
		}
		g.PruneNodes([]string{"b", "d"})
		if _, ok := g.Nodes["b"]; ok {
			t.Errorf("PruneNodes() expected to remove the Node b but it still exists")
		}
		if _, ok := g.Nodes["d"]; ok {
			t.Errorf("PruneNodes() expected to remove the Node d but it still exists")
		}
		if _, ok := g.Nodes["a"]["b"]; ok {
			t.Errorf("PruneNodes() expected to remove the edge to Node b but it still exists")
		}
		if _, ok := g.Nodes["a"]["d"]; ok {
			t.Errorf("PruneNodes() expected to remove the edge to Node d but it still exists")
		}
		if _, ok := g.Nodes["a"]["e"]; !ok {
			t.Errorf("PruneNodes() expected to not remove the edge to other Node e but removed")
		}
		if _, ok := g.Nodes["c"]["f"]; !ok {
			t.Errorf("PruneNodes() expected to not remove the edge to other Node f but removed")
		}
	})
}

func createTestGraph() Graph {
	return Graph{
		Nodes: map[string]Hashset{
			"a": {"b": true, "c": true},
			"b": {"c": true, "d": true},
			"c": {},
			"d": {},
		},
		mutex: new(sync.RWMutex),
	}
}

func BenchmarkGraph_AddOrReplaceNode(b *testing.B) {
	g := createBenchmarkGraph()
	for i := 0; i < b.N; i++ {
		g.AddOrReplaceNode("a", Hashset{"b": true})
	}
}

func BenchmarkGraph_PruneNode(b *testing.B) {
	g := createBenchmarkGraph()
	for i := 0; i < b.N; i++ {
		g.PruneNode("static1")
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
			"a":       {"b": true, "c": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"b":       {"c": true, "d": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"c":       {"d": true, "e": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"d":       {"e": true, "f": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"e":       {"f": true, "g": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"f":       {"g": true, "h": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"g":       {"h": true, "i": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"h":       {"i": true, "j": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"i":       {"j": true, "k": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"j":       {"k": true, "l": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"k":       {"l": true, "m": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"l":       {"m": true, "n": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"m":       {"n": true, "o": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"n":       {"o": true, "p": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
			"o":       {"p": true, "q": true, "static1": true, "static2": true, "static4": true, "static5": true, "empty1": true, "empty2": true},
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
