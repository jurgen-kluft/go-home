package main

import "sort"

type DataInput struct {
	ID             int64
	TrueState      bool
	ProbGivenTrue  float64
	ProbGivenFalse float64

	State bool
}

type Observation struct {
	ProbTrue  float64
	ProbFalse float64
}

type Bayesian struct {
	Name      string
	Prior     float64
	Threshold float64

	Inputs []DataInput

	Observations       []*Observation
	ObservationsSorted []int
	ID2Observation     map[int64]int

	Probability float64
}

func (b *Bayesian) ReadState() bool {
	return b.Probability >= b.Threshold
}

func (b *Bayesian) OnChange() {
	b.ProcessState()
	b.Update()
}

func (b *Bayesian) ProcessState() {
	for _, input := range b.Inputs {
		// Add entity to current observations if state conditions are met
		should_trigger := input.State
		if should_trigger {
			probtrue := input.ProbGivenTrue
			probfalse := input.ProbGivenFalse

			var obs *Observation
			obsi, exists := b.ID2Observation[input.ID]
			if !exists {
				obs = &Observation{ProbTrue: probtrue, ProbFalse: probfalse}
				obsi = len(b.Observations)

				b.Observations = append(b.Observations, obs)

				b.ObservationsSorted = append(b.ObservationsSorted, obsi)
				sort.Ints(b.ObservationsSorted)

				b.ID2Observation[input.ID] = obsi
			}
			obs.ProbTrue = probtrue
			obs.ProbFalse = probfalse
		} else {
			obsi, exists := b.ID2Observation[input.ID]
			if exists {
				sobsi := sort.SearchInts(b.ObservationsSorted, obsi)
				b.ObservationsSorted = append(b.ObservationsSorted[:sobsi], b.ObservationsSorted[sobsi+1:]...)
				delete(b.ID2Observation, input.ID)
			}
		}
	}
}

func (b *Bayesian) Update() {
	prior := b.Prior
	for _, obs := range b.Observations {
		prior = b.UpdateProbability(prior, obs.ProbTrue, obs.ProbFalse)
	}
	b.Probability = prior
}

func (b *Bayesian) UpdateProbability(prior float64, prob_true float64, prob_false float64) float64 {
	// Update probability using Bayes' rule.
	numerator := prob_true * prior
	denominator := numerator + prob_false*(1-prior)
	probability := numerator / denominator
	return probability
}

func main() {

}
