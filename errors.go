package rss

type InvalidStatusCodeError struct {
	StatusCode int
	Status     string
}

func (e *InvalidStatusCodeError) Temporary() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

func (e *InvalidStatusCodeError) Error() string {
	return e.Status
}
