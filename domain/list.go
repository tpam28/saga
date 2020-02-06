package domain

type StepList []Step

func (s *StepList) Push(step Step) error {
	if len(*s) > 1 {
		tmp := *s
		if tmp[len(*s)-1].T == Repeat && step.T != Repeat {
			return SequenceErr
		}

		if tmp[len(*s)-1].T == Compensatory {
			if step.T == Repeat {
				tmp[len(*s)-1].T = Critical
			}
		}

	}
	*s = append(*s, step)
	return nil
}

type Step struct {
	Name string
	T    TransactionType
	Sl   SemanticLockL
	Keys KeysL
}

type SemanticLockL struct {
	Pending  string
	Approval string
	Rejected string
}

type KeysL struct {
	First  string
	Second string
}
