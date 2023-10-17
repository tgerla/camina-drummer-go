package drummer

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const (
	sampleRate = 44100
)

type DrumMachine struct {
	State                 string
	Current_pattern_idx   int
	Pattern_length        int
	Current_measure       string
	Tempo                 int
	Measure_changing      bool
	Beat                  int
	Interval              int64
	Patterns              map[string]Pattern
	_sounds               map[string]*audio.Player
	_time_since_last_beat time.Duration
	_current_pattern      Pattern
	_playing_pattern      Pattern
	_prior_measure        string
	last_taps             []int
	_last_time            time.Time

	audioContext *audio.Context
}

func NewDrumMachine() *DrumMachine {
	dm := &DrumMachine{
		State:                 "stopped",
		Current_pattern_idx:   1,
		Pattern_length:        16,
		Current_measure:       "A",
		Tempo:                 120,
		Measure_changing:      false,
		Beat:                  0,
		Interval:              0,
		_time_since_last_beat: 0,
		_prior_measure:        "A",
		last_taps:             []int{},
	}

	dm.audioContext = audio.NewContext(sampleRate)

	dm._sounds = make(map[string]*audio.Player)
	for sound := range SOUNDS {
		f, err := os.Open("kits/1/" + SOUNDS[sound])
		if err != nil {
			log.Fatal(err)
		}

		s, err := wav.DecodeWithoutResampling(f)
		if err != nil {
			log.Fatal(err)
		}

		dm._sounds[sound], err = dm.audioContext.NewPlayer(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	pl := NewPatternLoader()
	dm.Patterns, _ = pl.LoadPatterns()

	dm.SwitchPattern(dm.Current_pattern_idx)
	dm.SetTempo(dm.Tempo)

	return dm
}

func (dm *DrumMachine) Play() {
	dm.State = "playing"
	dm._time_since_last_beat = 0
	dm._playing_pattern = dm._current_pattern
}

func (dm *DrumMachine) Stop() {
	dm.State = "stopped"
}

func (dm *DrumMachine) play_beat() {
	dm.Beat += 1
	if dm.Beat == dm.Pattern_length {
		dm.Beat = 0

		dm._playing_pattern = dm._current_pattern

		if dm.Current_measure == "T" {
			dm.Current_measure = dm._prior_measure
			dm.Measure_changing = true
		} else {
			dm.Measure_changing = false
		}
	}

	drums := dm._playing_pattern.Measures[dm.Current_measure]
	for drum := range drums {
		drumIdx := dm.Beat % len(drums[drum])
		if string(drums[drum][drumIdx]) == string('X') {
			dm._sounds[drum].Rewind()
			dm._sounds[drum].Play()
		}
	}
}

func (dm *DrumMachine) Tick() {
	deltaTime := time.Since(dm._last_time)
	//	log.Println(deltaTime, dm.Interval)
	dm._last_time = time.Now()

	if dm.State == "playing" {
		dm._time_since_last_beat += time.Duration(deltaTime)
		if dm._time_since_last_beat.Milliseconds() >= int64(dm.Interval) {
			dm._time_since_last_beat = 0
			dm.play_beat()
		}
	}
}

// set tempo function
func (dm *DrumMachine) SetTempo(tempo int) {
	dm.Tempo = tempo
	dm.Interval = int64((60 / float64(tempo) / 4) * 1000) // ms

	log.Println("set tempo: ", dm.Tempo)
	log.Println("set interval: ", dm.Interval)
}

func (dm *DrumMachine) SwitchPattern(new_pattern_idx int) {
	dm.Current_pattern_idx = new_pattern_idx
	dm._current_pattern = dm.Patterns[fmt.Sprint(new_pattern_idx)]
	log.Println("current pattern", dm._current_pattern)
	dm.Current_measure = "A"
	dm._playing_pattern = dm._current_pattern
	dm.Pattern_length = dm._current_pattern.Length
}
