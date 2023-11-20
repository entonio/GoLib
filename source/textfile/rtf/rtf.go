package rtf

type Rtf interface {
	PlainText() (string, error)
}
