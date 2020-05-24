package utils

var _empty struct{}

// StringSet provides set operations for sets of strings.
type StringSet map[string]struct{}

// NewStringSet instantiates a new StringSet from a list of keys.
func NewStringSet(keys []string) StringSet {
	A := StringSet{}
	for _, k := range keys {
		A.Add(k)
	}
	return A
}

// Add adds a key to a StringSet.
func (A StringSet) Add(key string) {
	A[key] = _empty
}

// Del removes a key from a StringSet.
func (A StringSet) Del(key string) {
	delete(A, key)
}

// Contains returns true iff A contains s.
func (A StringSet) Contains(s string) bool {
	_, ok := A[s]
	return ok
}

// ContainsSet returns true iff all keys of B are in A.
func (A StringSet) ContainsSet(B StringSet) bool {
	cntIn := 0
	for k := range B {
		if A.Contains(k) {
			cntIn++
		}
	}
	if cntIn == len(B) {
		return true
	}
	return false
}

// IsEqual returns true iff sets of keys of A and B are the same.
func (A StringSet) IsEqual(B StringSet) bool {
	if len(A) != len(B) {
		return false
	}
	return A.ContainsSet(B)
}

// IsSupersetOf returns true iff key set of A is a superset of key set of B.
func (A StringSet) IsSupersetOf(B StringSet) bool {
	return A.ContainsSet(B) && len(A) > len(B)
}

// IsSubsetOf returns true iff key set of A is a subset of key set of B.
func (A StringSet) IsSubsetOf(B StringSet) bool {
	return B.IsSupersetOf(A)
}
