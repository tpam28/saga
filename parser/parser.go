package parser

import (
	"fmt"

	"github.com/tpam28/saga/domain"
)

func NamesFromMap(m []map[string]map[string]interface{}) []string {
	res := make([]string, len(m))
	for k, v := range m {
		for name := range v {
			res[k] = name
			break
		}
	}
	return res
}


func GetSemantic(m map[string]interface{}) (*domain.SemanticLockL, error) {
	el, flag1 := m[string(domain.States)]
	if !flag1 {
		return nil, domain.ErrNoRows
	}
	tmp1, ok := el.(map[interface{}]interface{})
	if !ok {
		return nil, domain.KeysInvalidFormatErr
	}
	tmp := make(map[domain.TypeOfSematicLock]string)

	for k, v := range tmp1 {
		if domain.TypeOfSematicLock(fmt.Sprint(k)).Is() {
			tmp[domain.TypeOfSematicLock(fmt.Sprint(k))] = fmt.Sprint(v)
		}
	}
	res := &domain.SemanticLockL{}

	hash, ok := el.(map[interface{}]interface{})
	if !ok {
		return nil, domain.InvalidType
	}

	pending, ok := hash[string(domain.Pending)]
	if !ok {

	}
	approval, ok := hash[string(domain.Approval)]
	if !ok {

	}

	rejected, ok := hash[string(domain.Rejected)]
	if !ok {
		rejected = ""
	}
	res.Approval = fmt.Sprint(approval)
	res.Pending = fmt.Sprint(pending)
	res.Rejected = fmt.Sprint(rejected)
	return res, nil
}

func ParseStep(name string, m map[string]map[string]interface{}) (*domain.Step, error) {
	semantic, err := GetSemantic(m[name])
	if err != nil {
		return nil, err
	}
	semanticType := domain.Compensatable
	if semantic.Rejected == "" {
		semanticType = domain.Retriable
	}

	return &domain.Step{
		Name: name,
		T:    semanticType,
		Sl:   *semantic,
	}, nil
}

func ParseConfigSlice(ss []map[string]map[string]interface{}) (domain.StepList, error) {
	res := domain.StepList{}
	names := NamesFromMap(ss)
	for i, s := range ss {
		step, err := ParseStep(names[i], s)
		if err != nil {
			return nil, err
		}
		err = res.Push(*step)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
