package domain

type Root string

const (
	Milestone Root = "milestone"
	Build     Root = "build"
)

type TransactionType string

const (
	Compensatable TransactionType = "compensatable"
	Pivot         TransactionType = "pivot"
	Retriable     TransactionType = "retriable"
)

type Options string

const (
	Keys         Options = "keys"
	SemanticLock Options = "semanticLock"
)

type TypeOfSematicLock string

func (t TypeOfSematicLock) Is() bool {
	if t == Pending || t == Approval || t == Rejected {
		return true
	}
	return false
}

const (
	Pending  TypeOfSematicLock = "pending"
	Approval TypeOfSematicLock = "approval"
	Rejected TypeOfSematicLock = "rejected"
)
