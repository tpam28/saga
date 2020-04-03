package domain

type StepList []Step

func (s *StepList) Push(step Step) error {
	prev := len(*s)-1
	tmp := *s

	if len(*s) > 0 {
		tmp[prev].Next = &step
		step.Prev = &tmp[prev]
	}

	if len(*s) > 1 {
		//Returns err if sequence is invalid.
		if tmp[prev].T == Retriable && step.T != Retriable {
			return SequenceErr
		}
		//Change type to Critical if list of repeated tx is started.
		if tmp[prev].T == Compensatable {
			if step.T == Retriable {
				tmp[len(*s)-1].T = Pivot
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
	Prev *Step
	Next *Step
}

type SemanticLockL struct {
	Pending  string
	Approval string
	Rejected string
}
