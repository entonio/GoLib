package wordml

type WordML interface {
	PlainText() (string, error)
}
