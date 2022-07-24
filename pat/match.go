package pat

import (
	"sort"
	"sync"

	"goji.io/pattern"
)

type match struct {
	pat     *Pattern
	matches []string
}

type matches struct {
	mu      sync.Mutex
	matches []match
}

func (m *matches) AllVariables() map[pattern.Variable]any {
	m.mu.Lock()
	defer m.mu.Unlock()

	var vs map[pattern.Variable]interface{}

	for _, match := range m.matches {
		if len(m.matches) == 1 && len(match.pat.pats) == 0 {
			return nil
		}
		if vs == nil {
			vs = make(map[pattern.Variable]any, len(match.matches))
		}
		for _, p := range match.pat.pats {
			vs[p.name] = m.matches[p.idx]
		}
	}

	return vs
}

func (m *matches) Path() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.matches) == 0 {
		return ""
	}

	match := m.matches[len(m.matches)-1]

	if len(match.matches) == len(match.pat.pats)+1 {
		return match.matches[len(match.matches)-1]
	}

	return ""
}

func (m *matches) Variable(k pattern.Variable) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.matches) == 0 {
		return ""
	}

	for j := len(m.matches) - 1; j >= 0; j-- {
		match := m.matches[j]

		i := sort.Search(len(match.pat.pats), func(i int) bool {
			return match.pat.pats[i].name >= k
		})
		if i < len(match.pat.pats) && match.pat.pats[i].name == k {
			return match.matches[match.pat.pats[i].idx]
		}
	}

	return ""
}

func (m *matches) Add(match match) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.matches = append(m.matches, match)
}
