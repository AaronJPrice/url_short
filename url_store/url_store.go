package url_store

type UrlStore interface {
	// Compress a string to a (reasonably) short string, and store the mapping
	Compress(string) string
	// Expand a short string to a stored long string. Returns false if there
	// is no existing mapping
	Expand(string) (string, bool)
}
