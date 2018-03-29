package mist

import (
	"fmt"
	"strings"
)

type (
	// Filter ...
	Filter struct {
		filters    map[string][]int16
		varToIndex map[string]int16
		varValues  []bool
		rtstack    []bool
	}
)

const (
	opxmask = 0xF000
	opxand  = 0x1000
	opxor   = 0x2000
	opxxor  = 0x3000
	opxnot  = 0x4000
)

func newFilter() (filter *Filter) {
	filter = &Filter{
		filters:    map[string][]int16{},
		varToIndex: map[string]int16{},
		varValues:  []bool{},
		rtstack:    []bool{},
	}
	return
}

func (filter *Filter) addVar(varname string) int16 {
	i, exists := filter.varToIndex[varname]
	if !exists {
		i = int16(len(filter.varValues))
		filter.varToIndex[varname] = i
		filter.varValues = append(filter.varValues, false)
	}
	return i
}

func (filter *Filter) addFilter(expression string) error {
	elems := strings.Split(expression, " ")
	compiled := []int16{}
	for _, e := range elems {
		e = strings.Trim(e, " ")
		if e == "|" {
			compiled = append(compiled, opxor)
		} else if e == "&" {
			compiled = append(compiled, opxand)
		} else if e == "^" {
			compiled = append(compiled, opxxor)
		} else if e == "!" {
			compiled = append(compiled, opxnot)
		} else {
			vi := filter.addVar(e)
			compiled = append(compiled, vi)
		}
	}

	err := filter.validate(compiled)
	if err == nil {
		filter.filters[expression] = compiled
	}
	return err
}

// Add sorts the keys and then attempts to add them
func (filter *Filter) Add(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return filter.add(keys)
}

// add ...
func (filter *Filter) add(keys []string) error {
	errstr := ""
	for _, k := range keys {
		err := filter.addFilter(k)
		if err != nil {
			errstr = errstr + err.Error()
		}
	}
	if errstr == "" {
		return nil
	}
	return fmt.Errorf(errstr)
}

// Remove sorts the keys and then attempts to remove them
func (filter *Filter) Remove(keys []string) {

	if len(keys) == 0 {
		return
	}

	filter.remove(keys)
}

// remove ...
func (filter *Filter) remove(keys []string) {
	for _, k := range keys {
		delete(filter.filters, k)
	}
}

// Match sorts the keys and then attempts to find a match
func (filter *Filter) Match(keys []string) bool {
	return filter.match(keys)
}

func (filter *Filter) validate(expr []int16) error {
	r := []bool{}
	for _, s := range expr {
		i := len(r)
		switch s {
		case opxand:
			if i < 2 {
				return fmt.Errorf("RPN evaluation, expression is incorrect")
			}
			r = r[:i-1]
			break
		case opxor:
			if i < 2 {
				return fmt.Errorf("RPN evaluation, expression is incorrect")
			}
			r = r[:i-1]
			break
		case opxxor:
			if i < 2 {
				return fmt.Errorf("RPN evaluation, expression is incorrect")
			}
			r = r[:i-1]
			break
		case opxnot:
			if i < 1 {
				return fmt.Errorf("RPN evaluation, expression is incorrect")
			}
			break
		default:
			r = append(r, filter.varValues[s])
			break
		}
	}
	if len(r) == 1 {
		return nil
	}
	return fmt.Errorf("RPN evaluation results not in 1 answer")
}

func (filter *Filter) evaluate(expr []int16) bool {
	r := []bool{}
	for _, s := range expr {
		i := len(r)
		switch s {
		case opxand:
			r[i-2] = r[i-2] && r[i-1]
			r = r[:i-1]
			break
		case opxor:
			r[i-2] = r[i-2] || r[i-1]
			r = r[:i-1]
			break
		case opxxor:
			r[i-2] = (r[i-2] && !r[i-1]) || (!r[i-2] && r[i-1])
			r = r[:i-1]
			break
		case opxnot:
			r[i-1] = !r[i-1]
			break
		default:
			r = append(r, filter.varValues[s])
			break
		}
	}
	return r[0]
}

// ​match ...
func (filter *Filter) match(keys []string) bool {
	if len(keys) == 0 {
		return false
	}
	for i := range filter.varValues {
		filter.varValues[i] = false
	}
	for _, k := range keys {
		vi, exists := filter.varToIndex[k]
		if exists {
			filter.varValues[vi] = true
		}
	}

	for _, expr := range filter.filters {
		res := filter.evaluate(expr)
		if res {
			return true
		}
	}
	return false
}

// ToSlice ...
func (filter *Filter) ToSlice() (list [][]string) {
	for k := range filter.filters {
		list = append(list, []string{k})
	}
	return
}
