# camina-drummer-go
A drum loop player in Go.

This will eventually run on one of these: https://www.adafruit.com/product/802. Or maybe a RaspberryPi.

![Screenshot](screenshot.png)

Features:

- YAML patterns
- Three different measure styles per pattern, A, B, and transition/break
- Tap tempo

Prototype keyboard control:

```
Q or Esc - quit
Left arrow - previous pattern
Right arrow - next pattern
Spacebar - play/pause
Plus - increase tempo
Minus - decrease tempo
A/B/T - switch to a different measure style, A and B styles, plus "T" for a transition or break.
```

The samples are from https://github.com/jstrait/beats and patterns are from https://github.com/montoyamoraga/drum-machine-patterns, using my script to convert them to YAML at https://github.com/tgerla/drum-machine-patterns. 

TODO:

- swing
- switch measure styles
- tap tempo