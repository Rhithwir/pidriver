package pidriver

import (
	"fmt"
	"math"

	"github.com/stianeikeland/go-rpio"
)

// func main() {
// 	// usage: rpicmd board=36 phase=2 ...
// 	word := 0
// 	pi := connect()
// 	for i, a := range os.Args {
// 		if i != 0 {
// 			cmds := strings.Split(a, "=")
// 			cmd := cmds[0]
// 			val, _ := strconv.Atoi(cmds[1])
// 			switch cmd {
// 			case "rw": // both
// 				val = val & 1
// 				word += val & 1
// 			case "sel": // w only
// 				val = val & 7
// 				word += (val << 1) & 14
// 			case "quad": // r only
// 				val = val & 3
// 				word += (val << 1) & 6
// 			case "anadr": // r only
// 				val = val & 15
// 				word += (val << 3) & 120
// 			case "board": // both
// 				val = val & 127
// 				word += (val << 8) & 32512
// 			case "phase": // w only
// 				val = val & 3
// 				word += (val << 17) & 393216
// 			case "pwm": // w only
// 				val = val & 15
// 				word += (val << 19) & 7864320
// 			}
// 		}
// 	}
// 	fmt.Println(word)
// 	fmt.Println(strconv.FormatInt(int64(word), 2))

// }

//Rpi raspberry pi
type Rpi struct {
	Pins []rpio.Pin
}

// Connect gives you a new Rpi
func Connect() Rpi {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
	}
	pi := Rpi{}
	pinids := []int{2, 3, 4, 14, 15, 17, 18, 27, 22, 23, 24, 10, 9, 25, 11, 8, 7, 5, 6, 12, 13, 19, 18, 26, 20, 21}
	for i := range pinids {
		pi.Pins = append(pi.Pins, rpio.Pin(pinids[i]))
	}
	return pi
}

func (pi Rpi) strobe() {
	pi.Pins[15].Output()
	pi.Pins[15].High()
	pi.Pins[15].Low()
}

// Command commands the Pi GPIO
// set any unused args to 0, e.g. pi.command(0,0,24,12,2,0,0) or pi.command(12,1,24,0,0,2,1)
func (pi Rpi) Command(pwm, phase, board, anadr, quad, sel, rw int) (data int) {
	word := 0
	if rw == 0 { // read
		quad = quad & 3
		word += (quad << 1) & 6
		anadr = anadr & 15
		word += (anadr << 3) & 120
		board = board & 127
		word += (board << 8) & 32512
		data = pi.read(word)
	} else { // write
		sel = sel & 7
		word += (sel << 1) & 14
		board = board & 127
		word += (board << 8) & 32512
		phase = phase & 3
		word += (phase << 17) & 393216
		pwm = pwm & 15
		word += (pwm << 19) & 7864320
		pi.write(word)
		return -1
	}
	return
}

func (pi Rpi) write(word int) {
	if word > 8388607 {
		fmt.Println("word too large,", word, "> 8388607")
		return
	}
	for i, p := range pi.Pins {
		p.Output()
		if word&int(math.Pow(2, float64(i)))>>uint(i) == 1 {
			p.High()
		} else {
			p.Low()
		}
	}
	pi.strobe()
}

func (pi Rpi) read(cmd int) (data int) {
	pi.write(cmd)
	for i, p := range pi.Pins {
		if i < 8 {
			p.Input()
		} else {
			p.Output()
		}
	}
	pi.strobe()
	for i, p := range pi.Pins {
		if i < 8 {
			data = (data << 1) + int(p.Read())
		}
	}
	return data
}
