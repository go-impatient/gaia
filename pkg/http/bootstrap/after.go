package bootstrap

type AfterServerStopFunc func() error

// CloseDatabase ...
func CloseDatabase() AfterServerStopFunc {
	return func() error {
		return nil
	}
}
