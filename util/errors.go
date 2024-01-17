package util

func WrapErr(err error) func() error {
	return func() error { return err }
}
