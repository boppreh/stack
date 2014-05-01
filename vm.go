package stack

type OpCode int
type Program []OpCode

func RunOps(inputs []Value, ops []Op) ([]Value, error) {
	s := New(inputs)

	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			return nil, err
		}
	}

	return s.Exhaust(), nil
}

func Run(program Program) ([]Value, error) {
	s := new(Stack)

	for _, opCode := range program {
		switch opCode {
			case 0:
				s.Push(10)
		}
	}

	return s.Exhaust(), nil
}
