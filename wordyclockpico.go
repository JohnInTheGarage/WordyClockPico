package main

//to compile for pico : "tinygo flash -target=pico ./ "

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var (
	cmdChannel = make(chan ClockCommand, 10)
	hr, min    int
	prevMin    int = -1

	testRange     = [...]string{} //{"0928", "0929", "0930", "0931"}
	testIndex int = 0
)

func main() {

	time.Sleep(20 * time.Second)

	if !initRTC() {
		println("------------- Real Time Clock not available ---------------")
	}

	initClock()

	//led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//led.High()
	//time.Sleep(time.Second)
	//led.Low()
	//time.Sleep(time.Second)

	go ClockActuator(cmdChannel)
	fmt.Println("Clock actuator started.")

	if len(testRange) > 0 {
		useFakeTimes()
	} else {
		useRealTimes()
	}

}

func useRealTimes() {

	for {
		hr, min = getTime() // from attached RTC (for microcontrollers, or use getPCtime if on a PC)
		if min != prevMin {
			queueCommands(hr, min)
			// In a quiet period following major updates
			if min == 4 || min == 34 {
				//showMemStats()
				runtime.GC()
				//showMemStats()
			}

		}
		prevMin = min
		time.Sleep(5 * time.Second)
	}

}

func showMemStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Println("Heap used:", mem.HeapInuse, ", heap sys:", mem.HeapSys, ", heap released:", mem.HeapReleased)
}

func useFakeTimes() {
	for testIndex < len(testRange) {
		var item = testRange[testIndex]
		hr, err := strconv.Atoi(item[0:2])
		if err != nil {
			fmt.Printf("Test Hours incorrect? Failed with: '%s'\n", err)
		}
		min, err := strconv.Atoi(item[2:])
		if err != nil {
			fmt.Printf("Test Hours incorrect? Failed with: '%s'\n", err)
		}
		testIndex++
		queueCommands(hr, min)
		time.Sleep(15 * time.Second)
		//fmt.Printf(" testindex %d len(data) %d\n", testIndex, len(testRange))
	}

}

// For use on "PC" type hardware,
// i.e. with an operating system; doubt it would work on a microcontroller.
func getPCTime() (hrs, mins int) {
	now := time.Now()
	hrs = now.Hour()
	mins = now.Minute()
	return
}
