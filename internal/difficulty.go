package internal

type Difficulty interface {
	Update(difficulty Difficulty)
	Set(difficulty float64)
	CanMine(simulation *Simulation, previousBlock *Block, power float64) bool
	Adjust(event *Event)
	GetLambda(power float64) float64
}

func NewSlimcoinProofOfBurnDifficulty(enabled bool) Difficulty {
	return &SlimcoinProofOfBurnDifficulty{
		Enabled:    enabled,
		Difficulty: 1,
	}
}

type SlimcoinProofOfBurnDifficulty struct {
	Enabled    bool
	Difficulty float64
}

func (target *SlimcoinProofOfBurnDifficulty) Update(difficulty Difficulty) {
	_, ok := difficulty.(*SlimcoinProofOfBurnDifficulty)
	if !ok {
		panic("not a slimcoin proof of burn difficulty")
	}
}
func (d *SlimcoinProofOfBurnDifficulty) CanMine(simulation *Simulation, previousBlock *Block, power float64) bool {
	if !d.Enabled {
		return false
	}

	if previousBlock.Type != ProofOfWork {
		return false
	}

	chance := power / d.Difficulty // TODO: Model differently, should not be based on power
	return simulation.Random.Chance(chance)
}

func (d *SlimcoinProofOfBurnDifficulty) Set(difficulty float64) {
	d.Difficulty = difficulty
}

func (d *SlimcoinProofOfBurnDifficulty) GetLambda(power float64) float64 {
	return 1
}

func (d *SlimcoinProofOfBurnDifficulty) Adjust(event *Event) {
}

func NewProofOfWorkDifficulty(enabled bool, epochLength int64, blockFrequency float64) Difficulty {
	return &ProofOfWorkDifficulty{
		Enabled:          enabled,
		EpochIndex:       0,
		EpochLength:      epochLength,
		BlockFreqency:    blockFrequency,
		EpochTimeElapsed: 0,
		Difficulty:       1,
	}
}

type ProofOfWorkDifficulty struct {
	Enabled          bool
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

func (d *ProofOfWorkDifficulty) CanMine(simulation *Simulation, previousBlock *Block, difficulty float64) bool {
	return d.Enabled
}

func (d *ProofOfWorkDifficulty) Set(difficulty float64) {
	d.Difficulty = difficulty
}

func (d *ProofOfWorkDifficulty) GetLambda(power float64) float64 {
	return power / (d.BlockFreqency * d.Difficulty)
}

func (d *ProofOfWorkDifficulty) Adjust(event *Event) {
	if event.Block.Type != ProofOfWork {
		return
	}
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

		// logrus.Infof("Epoch Time: %f, Average Time: %f, Epoch Index: %d, Adjustment: %f", d.EpochTimeElapsed, d.EpochTimeElapsed/float64(d.EpochIndex), d.EpochIndex, adjustment)

		d.Difficulty *= adjustment
		d.EpochIndex = 0
		d.EpochTimeElapsed = 0
	}
}
