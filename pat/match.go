package pat

import (
	"context"
	"sort"
	"sync"

	"goji.io/pattern"
)

type matchesKey struct{}

type match struct {
	pat     *Pattern
	matches []string
}

type Matches struct {
	mu      sync.Mutex
	matches []match
}

func getMatches(ctx context.Context) *Matches {
	if mi := ctx.Value(matchesKey{}); mi != nil {
		if m, ok := mi.(*Matches); ok {
			return m
		}
	}

	return nil
}

func ensureMatches(ctx context.Context) (context.Context, *Matches) {
	if m := getMatches(ctx); m != nil {
		return ctx, m
	}

	ms := new(Matches)
	return context.WithValue(ctx, matchesKey{}, ms), ms
}

func (ms *Matches) add(m match) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.matches = append(ms.matches, m)
}

func (ms *Matches) AllVariables() map[pattern.Variable]interface{} {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	size := 0
	for _, m := range ms.matches {
		if m.pat == nil {
			continue
		}
		size += len(m.pat.pats)
	}

	if size == 0 {
		return nil
	}

	vs := make(map[pattern.Variable]interface{}, len(ms.matches))
	// iterate oldest to newest match here so if there are duplicate
	// variables the newest wins
	for _, m := range ms.matches {
		if m.pat == nil {
			continue
		}

		for _, p := range m.pat.pats {
			vs[p.name] = m.matches[p.idx]
		}
	}

	return vs
}

/*
GetAllVariables is a standard value which, when passed to context.Context.Value,
returns all variable bindings present in the context, with bindings in newer
contexts overriding values deeper in the stack. The concrete type

	map[Variable]interface{}

is used for this purpose. If no variables are bound, nil should be returned
instead of an empty map.
*/
func GetAllVariables(ctx context.Context) map[pattern.Variable]interface{} {
	ms := getMatches(ctx)
	if ms == nil {
		return nil
	}

	return ms.AllVariables()
}

func (ms *Matches) Path() string {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if len(ms.matches) == 0 {
		return ""
	}

	m := ms.matches[len(ms.matches)-1]
	if m.pat != nil && len(m.matches) == len(m.pat.pats)+1 {
		return m.matches[len(m.matches)-1]
	}
	return ""
}

func (ms *Matches) Variable(k pattern.Variable) (string, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// search most recent match to least recent
	for i := len(ms.matches) - 1; i >= 0; i-- {
		m := ms.matches[i]

		if m.pat == nil {
			continue
		}

		i := sort.Search(len(m.pat.pats), func(i int) bool {
			return m.pat.pats[i].name >= k
		})
		if i < len(m.pat.pats) && m.pat.pats[i].name == k {
			return m.matches[m.pat.pats[i].idx], true
		}
	}

	return "", false
}
