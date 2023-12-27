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

func TestGraph_RemoveNode(t *testing.T) {
	t.Run("Key already exists and gets removed", func(t *testing.T) {
		g := &Graph{
			Nodes: map[string]Hashset{
				"a": {"b": true},
			},
			mutex: new(sync.RWMutex),
		}
		g.RemoveNode("a")
		if _, ok := g.Nodes["a"]; ok {
			t.Errorf("RemoveNode() expected to remove the key but it still exists")
		}
	})
	t.Run("Key does not exist and continues to not exist", func(t *testing.T) {
		g := &Graph{
			Nodes: make(map[string]Hashset),
			mutex: new(sync.RWMutex),
		}
		g.RemoveNode("a")
		if _, ok := g.Nodes["a"]; ok {
			t.Errorf("RemoveNode() expected to remove the key but it exists")
		}
	})
}
