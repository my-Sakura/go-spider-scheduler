package Interface

type Task interface {
	Crawl(string) (string, error)
	Parse(string) error
}
