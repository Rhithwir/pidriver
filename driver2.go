package main

import (
	"fmt"
	"math/bits"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var pins = []rpio.Pin{
	rpio.Pin(4),
	rpio.Pin(25),
	rpio.Pin(24),
	rpio.Pin(23),
	rpio.Pin(22),
	rpio.Pin(27),
	rpio.Pin(18),
	rpio.Pin(2),
	rpio.Pin(3),
	rpio.Pin(8),
	rpio.Pin(7),
	rpio.Pin(10),
	rpio.Pin(9),
	rpio.Pin(11),
	rpio.Pin(6),
	rpio.Pin(13),
	rpio.Pin(19),
	rpio.Pin(26),
	rpio.Pin(12),
	rpio.Pin(16),
}
var strobe = rpio.Pin(5)
var rw = rpio.Pin(20)

func main() {

	rpio.Open()
	defer rpio.Close()

	strobe.Output()
	rw.Output()

	strobe.High()
	rw.High()

	board := uint(0)
	pwm := 15

	for sel := 0; sel < 8; sel++ {
		for phase := 0; phase < 4; phase++ {
			for pwm = 0; pwm < 16; pwm++ {
				t0 := time.Now()
				//	fmt.Println(pwm, phase, board, sel)
				write(uint(pwm), uint(phase), board, uint(sel))
				//time.Sleep(25 * time.Millisecond)
				fmt.Println(time.Since(t0))
			}
		}
	}

	for sel := 0; sel < 8; sel++ {
		write(0, 0, board, uint(sel))
	}
	// i := 0
	// f := ""
	// for i < 100 {
	for quad := 0; quad < 4; quad++ {
		for anadr := 0; anadr < 16; anadr++ {
			//fmt.Println(quad, anadr)
			t0 := time.Now()
			v := read(board, uint(quad), uint(anadr))
			fmt.Println(time.Since(t0))
			fmt.Println(v)
			///	f += strconv.FormatUint(uint64(v), 10) + " "
		}
		//	f += "\n"
	}
	//f += "\n"
	// 	i++
	// }
	// files.Write(f, "8.txt")
	// for {
	// 	t0 := time.Now()
	// 	for sel := 0; sel < 8; sel++ {
	// 		write(15, 0, board, uint(sel))
	// 	}
	// 	for sel := 0; sel < 8; sel++ {
	// 		write(0, 0, board, uint(sel))
	// 	}
	// 	fmt.Println(time.Since(t0))
	// }
}

func write(pwm, phase, board, sel uint) {
	word := sel + board<<6 + uint(bits.Reverse8(uint8(pwm)<<4))<<16 + phase<<14
	for i, p := range pins {
		p.Output()
		switch word >> uint(len(pins)-1-i) & 1 {
		case 0:
			p.Low()
		case 1:
			p.High()
		}
	}
	strobe.Low()
	strobe.High()
}

func read(board, quad, anadr uint) (data int) {
	// write command
	word := anadr + quad<<4 + board<<6
	for i, p := range pins {
		p.Output()
		switch word >> uint(len(pins)-1-i) & 1 {
		case 0:
			p.Low()
		case 1:
			p.High()
		}
	}
	strobe.Low()
	strobe.High()
	rw.Low()
	// get data
	for i := 0; i < 8; i++ {
		p := pins[i]
		p.Input()
		data += int(p.Read()) << uint(8-i-1)
	}
	// return
	return
}
