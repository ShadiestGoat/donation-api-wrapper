package donations

func withError[T any](v *T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}

	return v, nil
}

func withErrorArr[T any](v []*T, err error) ([]*T, error) {
	if err != nil {
		return nil, err
	}

	return v, nil
}