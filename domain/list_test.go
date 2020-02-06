package domain

import (
	"testing"
)

func TestStepList_Push(t *testing.T) {
	var sl StepList
	err := sl.Push(Step{T: Compensatory})
	if err != nil {
		t.Fatal(err)
	}
	err = sl.Push(Step{T: Compensatory})
	if err != nil {
		t.Fatal(err)
	}
	err = sl.Push(Step{T: Compensatory})
	if err != nil {
		t.Fatal(err)
	}
	err = sl.Push(Step{T: Repeat})
	if err != nil {
		t.Fatal(err)
	}
	if sl[len(sl)-2].T != Critical {
		t.Fatal("expected change type", sl[len(sl)-1].T)
	}
	err = sl.Push(Step{T: Repeat})
	if err != nil {
		t.Fatal(err)
	}
	err = sl.Push(Step{T: Compensatory})
	if err == nil {
		t.Fatal("expected err")
	}
}
