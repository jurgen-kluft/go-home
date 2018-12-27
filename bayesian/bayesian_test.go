package bayesian

import (
	"math"
	"reflect"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Instance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(0.5, 0.9); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstance_AddInput(t *testing.T) {
	type args struct {
		id             int64
		truestate      bool
		probGivenTrue  float64
		probGivenFalse float64
	}
	tests := []struct {
		name string
		b    *Instance
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.AddInput(tt.args.id, tt.args.truestate, tt.args.probGivenTrue, tt.args.probGivenFalse)
		})
	}
}

func TestInstance_SetInputState(t *testing.T) {
	type args struct {
		id    int64
		state bool
	}
	tests := []struct {
		name string
		b    *Instance
		args args
		want bool
	}{
		{name: "test input 1", b: &Instance{inputs: map[int64]*dataInput{0: &dataInput{id: 0, trueState: true, probGivenTrue: 0.95, probGivenFalse: 0.7, state: true}}, id2Observation: map[int64]int{}, prior: 0.7, threshold: 0.9}, args: args{id: 0, state: false}, want: false},
		{name: "test input 2", b: &Instance{inputs: map[int64]*dataInput{0: &dataInput{id: 0, trueState: true, probGivenTrue: 0.95, probGivenFalse: 0.7, state: true}}, id2Observation: map[int64]int{}, prior: 0.7, threshold: 0.9}, args: args{id: 0, state: true}, want: false},
		{name: "test input 3", b: &Instance{inputs: map[int64]*dataInput{0: &dataInput{id: 0, trueState: true, probGivenTrue: 0.95, probGivenFalse: 0.7, state: true}}, id2Observation: map[int64]int{}, prior: 0.7, threshold: 0.9}, args: args{id: 0, state: false}, want: false},
		{name: "test input 4", b: &Instance{inputs: map[int64]*dataInput{0: &dataInput{id: 0, trueState: true, probGivenTrue: 0.95, probGivenFalse: 0.7, state: true}}, id2Observation: map[int64]int{}, prior: 0.7, threshold: 0.9}, args: args{id: 0, state: true}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.SetInputState(tt.args.id, tt.args.state)
			if got := tt.b.ReadState(); got != tt.want {
				t.Errorf("Instance.SetInputState() -> Instance.ReadState() = %v, want %v (obs: %v)", got, tt.want, len(tt.b.observations))
			}
		})
	}
}

func TestInstance_ReadState(t *testing.T) {
	tests := []struct {
		name string
		b    *Instance
		want bool
	}{
		{name: "test true", b: &Instance{inputs: map[int64]*dataInput{0: &dataInput{id: 0, trueState: true, probGivenTrue: 0.95, probGivenFalse: 0.2, state: true}}, id2Observation: map[int64]int{}, prior: 0.7, threshold: 0.8}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.ReadState(); got != tt.want {
				t.Errorf("Instance.ReadState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstance_processState(t *testing.T) {
	tests := []struct {
		name string
		b    *Instance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.processState()
		})
	}
}

func TestInstance_updateState(t *testing.T) {
	tests := []struct {
		name string
		b    *Instance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.updateState()
		})
	}
}

func TestInstance_computeProbability(t *testing.T) {
	type args struct {
		prior     float64
		probTrue  float64
		probFalse float64
	}
	tests := []struct {
		name string
		b    *Instance
		args args
		want float64
	}{
		{name: "test1", b: &Instance{}, args: args{prior: 0.5, probTrue: 0.90, probFalse: 0.2}, want: float64((0.90 * 0.5) / ((0.90 * 0.5) + 0.2*(1.0-0.5)))},
		{name: "test2", b: &Instance{}, args: args{prior: 0.5, probTrue: 0.95, probFalse: 0.7}, want: float64((0.95 * 0.5) / ((0.95 * 0.5) + 0.7*(1.0-0.5)))},
		{name: "test3", b: &Instance{}, args: args{prior: 0.2, probTrue: 0.90, probFalse: 0.2}, want: float64((0.90 * 0.2) / ((0.90 * 0.2) + 0.2*(1.0-0.2)))},
		{name: "test4", b: &Instance{}, args: args{prior: 0.2, probTrue: 0.95, probFalse: 0.7}, want: float64((0.95 * 0.2) / ((0.95 * 0.2) + 0.7*(1.0-0.2)))},
		{name: "test5", b: &Instance{}, args: args{prior: 0.7, probTrue: 0.90, probFalse: 0.2}, want: float64((0.90 * 0.7) / ((0.90 * 0.7) + 0.2*(1.0-0.7)))},
		{name: "test6", b: &Instance{}, args: args{prior: 0.7, probTrue: 0.95, probFalse: 0.7}, want: float64((0.95 * 0.7) / ((0.95 * 0.7) + 0.7*(1.0-0.7)))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.computeProbability(tt.args.prior, tt.args.probTrue, tt.args.probFalse); !almostEqual(got, tt.want) {
				t.Errorf("Instance.computeProbability() = %v, want %v", got, tt.want)
			}
		})
	}
}
