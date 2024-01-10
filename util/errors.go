package util

func Wrap(err error) func() error {
	return func() error { return err }
}
