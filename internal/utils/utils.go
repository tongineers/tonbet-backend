package utils

func DupSliceOfPointers[T any](src []*T) []*T {
	c := make([]*T, len(src))
	for i, p := range src {
		if p == nil {
			continue
		}

		v := *p
		c[i] = &v
	}

	return c
}
