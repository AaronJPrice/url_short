package shortener

type Shortener interface {
	// result should be reasonably short
	// collisions should be unlikely
	Shorten(string) string
}
