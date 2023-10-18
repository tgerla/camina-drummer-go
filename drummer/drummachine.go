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
	State               string
	CurrentPatternIndex int
	PatternLength       int
	CurrentMeasure      string
	Tempo               int
	MeasureChanging     bool
	Beat                int
	Interval            int64
	Patterns            map[string]Pattern
	sounds              map[string]*audio.Player
	timeSinceLastBeat   time.Duration
	currentPattern      Pattern
	playingPattern      Pattern
	priorMeasure        string
	lastTaps            []int
	lastTime            time.Time

	audioContext *audio.Context
}

func NewDrumMachine() *DrumMachine {
	dm := &DrumMachine{
		State:               "stopped",
		CurrentPatternIndex: 1,
		PatternLength:       16,
		CurrentMeasure:      "A",
		Tempo:               120,
		MeasureChanging:     false,
		Beat:                0,
		Interval:            0,
		timeSinceLastBeat:   0,
		priorMeasure:        "A",
		lastTaps:            []int{},
	}

	dm.audioContext = audio.NewContext(sampleRate)

	dm.sounds = make(map[string]*audio.Player)
	for sound := range SOUNDS {
		f, err := os.Open("kits/1/" + SOUNDS[sound])
		if err != nil {
			log.Fatal(err)
		}

		s, err := wav.DecodeWithoutResampling(f)
		if err != nil {
			log.Fatal(err)
		}

		dm.sounds[sound], err = dm.audioContext.NewPlayer(s)
		if err != nil {
			log.Fatal(err)
		}
	}

	pl := NewPatternLoader()
	dm.Patterns, _ = pl.LoadPatterns()

	dm.SwitchPattern(dm.CurrentPatternIndex)
	dm.SetTempo(dm.Tempo)

	return dm
}

func (dm *DrumMachine) Play() {
	dm.State = "playing"
	dm.timeSinceLastBeat = 0
	dm.playingPattern = dm.currentPattern
}

func (dm *DrumMachine) Stop() {
	dm.State = "stopped"
}

func (dm *DrumMachine) play_beat() {
	dm.Beat += 1
	if dm.Beat == dm.PatternLength {
		dm.Beat = 0

		dm.playingPattern = dm.currentPattern

		if dm.CurrentMeasure == "T" {
			dm.CurrentMeasure = dm.priorMeasure
			dm.MeasureChanging = true
		} else {
			dm.MeasureChanging = false
		}
	}

	drums := dm.playingPattern.Measures[dm.CurrentMeasure]
	for drum := range drums {
		drumIdx := dm.Beat % len(drums[drum])
		if string(drums[drum][drumIdx]) == string('X') {
			dm.sounds[drum].Rewind()
			dm.sounds[drum].Play()
		}
	}
}

func (dm *DrumMachine) Tick() {
	deltaTime := time.Since(dm.lastTime)
	//	log.Println(deltaTime, dm.Interval)
	dm.lastTime = time.Now()

	if dm.State == "playing" {
		dm.timeSinceLastBeat += time.Duration(deltaTime)
		if dm.timeSinceLastBeat.Milliseconds() >= int64(dm.Interval) {
			dm.timeSinceLastBeat = 0
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
	dm.CurrentPatternIndex = new_pattern_idx
	dm.currentPattern = dm.Patterns[fmt.Sprint(new_pattern_idx)]
	log.Println("current pattern", dm.currentPattern)
	dm.CurrentMeasure = "A"
	dm.playingPattern = dm.currentPattern
	dm.PatternLength = dm.currentPattern.Length
}
