package bayesian

import "sort"

type dataInput struct {
	id             int64
	trueState      bool
	probGivenTrue  float64
	probGivenFalse float64
	state          bool
}

type observation struct {
	probTrue  float64
	probFalse float64
}

// Instance is the Bayesian object
type Instance struct {
	name               string
	inputs             map[int64]*dataInput
	observations       []*observation
	observationsSorted []int
	id2Observation     map[int64]int

	prior       float64
	threshold   float64
	probability float64
}

// New creates a new instance of a Bayesian object
func New() *Instance {
	inst := &Instance{}

	return inst
}

// AddInput is there to add an input to the bayesian observation
func (b *Instance) AddInput(id int64, truestate bool, probGivenTrue float64, probGivenFalse float64) {
	input := &dataInput{id: id, trueState: truestate, probGivenTrue: probGivenTrue, probGivenFalse: probGivenFalse, state: false}
	b.inputs[id] = input
}

// SetInputState will set the input state of one of the inputs that is identified by ID 'id'
func (b *Instance) SetInputState(id int64, state bool) {
	input, contains := b.inputs[id]
	if contains {
		input.state = state
	}
}

// ReadState will return true/false according to the bayesian computation
func (b *Instance) ReadState() bool {
	b.processState()
	b.updateState()
	return b.probability >= b.threshold
}

// Example will test all possible combinations
// presence := Bayesian.New()
// presence.Add("livingroom sensor", 0.95, 0.7)
func Example() {
	ishome := New()
	ishome.AddInput(0, true, 0.9, 0.2) // iPhoneX is home/away
	ishome.AddInput(1, true, 0.9, 0.2) // Macbook is home/away
	ishome.AddInput(2, true, 0.9, 0.2) // Kindle is home/away

	// Possibility A
	ishome.SetInputState(0, false) // iPhoneX is home
	ishome.SetInputState(0, false) // Macbook is away
	ishome.SetInputState(0, false) // Kindle is away
	ishome.ReadState()

	// Possibility B
	ishome.SetInputState(0, true)  // iPhoneX is home
	ishome.SetInputState(0, false) // Macbook is away
	ishome.SetInputState(0, false) // Kindle is away
	ishome.ReadState()

	// Possibility C
	ishome.SetInputState(0, false) // iPhoneX is home
	ishome.SetInputState(0, true)  // Macbook is away
	ishome.SetInputState(0, false) // Kindle is away
	ishome.ReadState()

	// Possibility D
	ishome.SetInputState(0, false) // iPhoneX is home
	ishome.SetInputState(0, false) // Macbook is away
	ishome.SetInputState(0, true)  // Kindle is away
	ishome.ReadState()

	// Possibility E
	ishome.SetInputState(0, true)  // iPhoneX is home
	ishome.SetInputState(0, true)  // Macbook is away
	ishome.SetInputState(0, false) // Kindle is away
	ishome.ReadState()
}

func (b *Instance) processState() {
	for _, input := range b.inputs {
		// Add entity to current observations if state conditions are met
		if input.state {
			probtrue := input.probGivenTrue
			probfalse := input.probGivenFalse

			var obs *observation
			obsi, exists := b.id2Observation[input.id]
			if !exists {
				obs = &observation{probTrue: probtrue, probFalse: probfalse}
				obsi = len(b.observations)

				b.observations = append(b.observations, obs)

				b.observationsSorted = append(b.observationsSorted, obsi)
				sort.Ints(b.observationsSorted)

				b.id2Observation[input.id] = obsi
			}
			obs.probTrue = probtrue
			obs.probFalse = probfalse
		} else {
			obsi, exists := b.id2Observation[input.id]
			if exists {
				sobsi := sort.SearchInts(b.observationsSorted, obsi)
				b.observationsSorted = append(b.observationsSorted[:sobsi], b.observationsSorted[sobsi+1:]...)
				delete(b.id2Observation, input.id)
			}
		}
	}
}

func (b *Instance) updateState() {
	prior := b.prior
	for _, obs := range b.observations {
		prior = b.computeProbability(prior, obs.probTrue, obs.probFalse)
	}
	b.probability = prior
}

func (b *Instance) computeProbability(prior float64, probTrue float64, probFalse float64) float64 {
	// Update probability using Bayes' rule.
	numerator := probTrue * prior
	denominator := numerator + probFalse*(1-prior)
	probability := numerator / denominator
	return probability
}
