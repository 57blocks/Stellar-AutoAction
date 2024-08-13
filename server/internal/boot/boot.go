package boot

type (
	Initializer interface {
		Boot() error
	}

	WrapFunc func() error
)

func (f WrapFunc) Boot() error {
	return f()
}

func Wrap(f func() error) WrapFunc {
	return f
}

func Boots(inits ...Initializer) error {
	for _, initializer := range inits {
		err := initializer.Boot()
		if err != nil {
			return err
		}
	}

	return nil
}
