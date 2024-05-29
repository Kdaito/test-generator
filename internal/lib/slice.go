package lib

func Find[S ~[]*E, E any](slice S, f func(*E) bool) *E {
	for _, s := range slice {
		result := f(s)

		if result {
			return s
		}
	}
	return nil
}