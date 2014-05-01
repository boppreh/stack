package stack

func RunOps(inputs []Value, ops []Op) ([]Value, error) {
	s := New(inputs)

	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			return nil, err
		}
	}

	return s.Exhaust(), nil
}
