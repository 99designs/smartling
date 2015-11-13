package main

type stringSlice []string

func (ss stringSlice) contains(s string) bool {
	for _, t := range ss {
		if t == s {
			return true
		}
	}
	return false
}
