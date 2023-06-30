package components

// componentFactory
type componentFactory struct {
}

// NewComponentFactory
func NewComponentFactory() (*componentFactory, error) {
	return &componentFactory{}, nil
}
