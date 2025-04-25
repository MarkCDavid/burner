package internal

import "github.com/sirupsen/logrus"

type Difficulty interface {
	Update(difficulty Difficulty)
	Set(difficulty float64)
	Adjust(event *Event)
	GetLambda(power float64) float64
}

func NewProofOfWorkDifficulty(epochLength int64, blockFrequency float64) Difficulty {
	return &ProofOfWorkDifficulty{
		EpochIndex:       0,
		EpochLength:      epochLength,
		BlockFreqency:    blockFrequency,
		EpochTimeElapsed: 0,
		Difficulty:       1,
	}
}

type ProofOfWorkDifficulty struct {
	EpochIndex       int64
	EpochLength      int64
	BlockFreqency    float64
	EpochTimeElapsed float64
	Difficulty       float64
}

func (target *ProofOfWorkDifficulty) Update(difficulty Difficulty) {
	source, ok := difficulty.(*ProofOfWorkDifficulty)
	if !ok {
		panic("not a proof of work difficulty")
	}
	target.EpochIndex = source.EpochIndex
	target.EpochLength = source.EpochLength
	target.BlockFreqency = source.BlockFreqency
	target.EpochTimeElapsed = source.EpochTimeElapsed
	target.Difficulty = source.Difficulty
}

func (d *ProofOfWorkDifficulty) Set(difficulty float64) {
	d.Difficulty = difficulty
}

func (d *ProofOfWorkDifficulty) GetLambda(power float64) float64 {
	return power / (d.BlockFreqency * d.Difficulty)
}

func (d *ProofOfWorkDifficulty) Adjust(event *Event) {
	d.EpochIndex += 1
	d.EpochTimeElapsed += event.Duration()

	if d.EpochIndex >= d.EpochLength {
		adjustment := (d.BlockFreqency * float64(d.EpochLength)) / d.EpochTimeElapsed
		if adjustment > 4 {
			adjustment = 4
		}
		if adjustment < 0.25 {
			adjustment = 0.25
		}

		logrus.Infof("Epoch Time: %f, Average Time: %f, Epoch Index: %d, Adjustment: %f", d.EpochTimeElapsed, d.EpochTimeElapsed/float64(d.EpochIndex), d.EpochIndex, adjustment)

		d.Difficulty *= adjustment
		d.EpochIndex = 0
		d.EpochTimeElapsed = 0
	}
}
