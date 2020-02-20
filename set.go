package main

import "sort"

// Set is a set of strings
type Set map[string]struct{}

// NewSet creates a new set
func NewSet() *Set {
	s := Set(make(map[string]struct{}))
	return &s
}

// Add adds a single value to the set
func (s *Set) Add(value string) *Set {
	(*s)[value] = struct{}{}
	return s
}

// AddAll adds all values to the set
func (s *Set) AddAll(values []string) *Set {
	for _, v := range values {
		s.Add(v)
	}
	return s
}

// Contains checks if the set contains the given value
func (s Set) Contains(value string) bool {
	if _, ok := s[value]; ok {
		return true
	}
	return false
}

func (s Set) isValid(validSet Set) bool {
	for i := range s {
		if !validSet.Contains(i) {
			return false
		}
	}
	return true
}

func (s Set) containsAtLeastOne(validSet Set) bool {
	for i := range validSet {
		if s.Contains(i) {
			return true
		}
	}
	return false
}

// ToSlice transforms a set of strings into a slice of strings
func (s *Set) ToSlice() []string {
	keys := make([]string, len(*s))

	i := 0
	for k := range *s {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}