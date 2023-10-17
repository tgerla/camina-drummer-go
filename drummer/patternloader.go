package drummer

import (
	"encoding/json"
	"fmt"
	"os"
)

const PATTERNS_FILE = "assets/patterns.json"

var SOUNDS = map[string]string{
	"AC": "clave_low_stereo_16.wav",
	"BD": "bass_stereo_16.wav",
	"SD": "snare_stereo_16.wav",
	"LT": "tom1_stereo_16.wav",
	"MT": "tom2_stereo_16.wav",
	"HT": "tom3_stereo_16.wav",
	"CH": "hh_closed_stereo_16.wav",
	"OH": "hh_open_stereo_16.wav",
	"CY": "ride_stereo_16.wav",
	"RS": "rim_stereo_16.wav",
	"CP": "clave_high_stereo_16.wav",
	"CB": "cowbell_high_stereo_16.wav",
}

type Pattern struct {
	Name       string                       `json:"name"`
	TimeSig    string                       `json:"time_signature"`
	TempoRange string                       `json:"tempo_range"`
	Length     int                          `json:"length"`
	Measures   map[string]map[string]string `json:"measures"`
}

type Patterns struct {
	Patterns map[string]Pattern `json:"patterns"`
}

type PatternLoader struct {
	filename string
}

func NewPatternLoader() *PatternLoader {
	return &PatternLoader{
		filename: PATTERNS_FILE,
	}
}

func (p *PatternLoader) LoadPatterns() (map[string]Pattern, error) {
	data, err := os.ReadFile(p.filename)
	if err != nil {
		return make(map[string]Pattern), err
	}

	patterns := Patterns{}
	err = json.Unmarshal([]byte(data), &patterns)
	if err != nil {
		return make(map[string]Pattern), err
	}
	fmt.Printf("Loaded %d patterns from %s\n", len(patterns.Patterns), p.filename)

	return patterns.Patterns, nil
}
