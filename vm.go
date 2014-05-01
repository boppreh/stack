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

type OpCode int
type Program []OpCode

func Run(program Program) ([]Value, error) {
	s := new(Stack)

	r := s.Push
	p := s.Pop

	for _, opCode := range program {
		switch opCode {
			case 0:
				r(0)
			case 1:
				r(p() + 1)
			case 2:
				r(p() + p())
		}
	}

	return s.Exhaust(), nil
}
