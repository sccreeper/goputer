package sound

import (
	"errors"
	"math"
	"sccreeper/goputer/pkg/constants"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

var sr beep.SampleRate = beep.SampleRate(44100)

//https://dev.to/thilanka/sine-wave-generator-using-golang-nom

type SineWave struct {
	sample_factor float64 // Just for ease of use so that we don't have to calculate every sample
	phase         float64
	volume        float64
}

func (g *SineWave) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples { // increment = ((2 * PI) / SampleRate) * freq
		v := math.Sin(g.phase * 2.0 * math.Pi) // period of the wave is thus defined as: 2 * PI.
		samples[i][0] = v * g.volume
		samples[i][1] = v * g.volume
		_, g.phase = math.Modf(g.phase + g.sample_factor)
	}

	return len(samples), true
}

func (*SineWave) Err() error {
	return nil
}

type SquareWave struct {
	sample_factor float64
	phase         float64
	volume        float64
}

func (g *SquareWave) Stream(samples [][2]float64) (n int, ok bool) {

	for i := range samples {

		var v float64

		if g.phase <= math.Pi {
			v = 1.0
		} else {
			v = -1.0
		}

		samples[i][0] = v * float64(g.phase) * math.Pi * g.volume
		samples[i][1] = v * float64(g.phase) * math.Pi * g.volume

		_, g.phase = math.Modf(g.phase + g.sample_factor)

	}

	return len(samples), true

}

func (g *SquareWave) Err() error {
	return nil
}

func SineTone(sr beep.SampleRate, freq float64, volume float64) (beep.Streamer, error) {
	dt := freq / float64(sr)

	if dt >= 1.0/2.0 {
		return nil, errors.New("samplerate must be at least 2 times grater then frequency")
	}

	return &SineWave{dt, 0.1, volume}, nil
}

func SquareTone(sr beep.SampleRate, freq float64, volume float64) (beep.Streamer, error) {

	dt := freq / float64(sr)

	if dt >= 1.0/2.0 {
		return nil, errors.New("samplerate must be at least 2 times grater then frequency")
	}

	return &SquareWave{dt, 0.1, volume}, nil

}

func SoundInit() {
	speaker.Init(sr, sr.N(time.Second/10))
}

func PlaySound(sound_type uint32, sound_tone uint32, volume uint32) {

	var str beep.Streamer

	speaker.Clear()

	volume_float := float64(volume) / 255.0

	switch constants.SoundWaveType(sound_type) {
	case constants.SWSine:
		str, _ = SineTone(sr, float64(sound_tone), volume_float)
	case constants.SWSquare:
		str, _ = SquareTone(sr, float64(sound_tone), volume_float)

	}

	speaker.Play(str)

}
