package ports

// Filter defines criteria for excluding ports from scan results.
type Filter struct {
	// ExcludePorts is a set of port numbers to ignore.
	ExcludePorts map[int]struct{}
	// ExcludeProtos is a set of protocol strings ("tcp", "udp") to ignore.
	ExcludeProtos map[string]struct{}
}

// NewFilter constructs a Filter from slices of ports and protocols.
func NewFilter(excludePorts []int, excludeProtos []string) *Filter {
	f := &Filter{
		ExcludePorts:  make(map[int]struct{}, len(excludePorts)),
		ExcludeProtos: make(map[string]struct{}, len(excludeProtos)),
	}
	for _, p := range excludePorts {
		f.ExcludePorts[p] = struct{}{}
	}
	for _, proto := range excludeProtos {
		f.ExcludeProtos[proto] = struct{}{}
	}
	return f
}

// Allow returns true when the given PortState should be kept (not filtered out).
func (f *Filter) Allow(ps PortState) bool {
	if f == nil {
		return true
	}
	if _, excluded := f.ExcludePorts[ps.Port]; excluded {
		return false
	}
	if _, excluded := f.ExcludeProtos[ps.Proto]; excluded {
		return false
	}
	return true
}

// Apply returns a new slice containing only the PortStates that pass the filter.
func (f *Filter) Apply(states []PortState) []PortState {
	result := make([]PortState, 0, len(states))
	for _, s := range states {
		if f.Allow(s) {
			result = append(result, s)
		}
	}
	return result
}
