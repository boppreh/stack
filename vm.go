package stack

func Run(inputs []Value, ops []Op) ([]Value, error) {
	s := new(Stack)

	for _, value := range inputs {
		s.Push(value)
	}

	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			return nil, err
		}
	}

	return s.Exhaust(), nil
}
