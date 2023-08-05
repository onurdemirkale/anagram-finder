package inputsource

type InputSource interface {
	GetWords() ([]string, error)
}
