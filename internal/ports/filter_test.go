package ports

import (
	"testing"
)

func TestFilter_Allow_NilFilter(t *testing.T) {
	var f *Filter
	ps := PortState{Port: 80, Proto: "tcp", Open: true}
	if !f.Allow(ps) {
		t.Error("nil filter should allow all ports")
	}
}

func TestFilter_Allow_ExcludedPort(t *testing.T) {
	f := NewFilter([]int{8080}, nil)
	if f.Allow(PortState{Port: 8080, Proto: "tcp", Open: true}) {
		t.Error("expected port 8080 to be excluded")
	}
}

func TestFilter_Allow_NonExcludedPort(t *testing.T) {
	f := NewFilter([]int{8080}, nil)
	if !f.Allow(PortState{Port: 443, Proto: "tcp", Open: true}) {
		t.Error("expected port 443 to be allowed")
	}
}

func TestFilter_Allow_ExcludedProto(t *testing.T) {
	f := NewFilter(nil, []string{"udp"})
	if f.Allow(PortState{Port: 53, Proto: "udp", Open: true}) {
		t.Error("expected udp proto to be excluded")
	}
}

func TestFilter_Allow_NonExcludedProto(t *testing.T) {
	f := NewFilter(nil, []string{"udp"})
	if !f.Allow(PortState{Port: 53, Proto: "tcp", Open: true}) {
		t.Error("expected tcp proto to be allowed")
	}
}

func TestFilter_Apply_FiltersCorrectly(t *testing.T) {
	f := NewFilter([]int{22, 8080}, []string{"udp"})
	input := []PortState{
		{Port: 22, Proto: "tcp", Open: true},
		{Port: 80, Proto: "tcp", Open: true},
		{Port: 53, Proto: "udp", Open: true},
		{Port: 443, Proto: "tcp", Open: true},
		{Port: 8080, Proto: "tcp", Open: true},
	}
	got := f.Apply(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].Port != 80 || got[1].Port != 443 {
		t.Errorf("unexpected filtered results: %v", got)
	}
}

func TestFilter_Apply_EmptyInput(t *testing.T) {
	f := NewFilter([]int{80}, nil)
	got := f.Apply([]PortState{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestNewFilter_EmptySlices(t *testing.T) {
	f := NewFilter(nil, nil)
	ps := PortState{Port: 9999, Proto: "tcp", Open: true}
	if !f.Allow(ps) {
		t.Error("filter with no exclusions should allow all")
	}
}
