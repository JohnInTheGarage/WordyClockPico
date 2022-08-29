package main

import (
	"fmt"
	"machine"
	"strings"
	"time"
)

const (
	HARDWARE_PRESENT = false // determines if the calibrate code should be run
	//LED_ONBOARD int = 13
	SWITCH_PIN int = 7

	//Stepper motor clues
	FORWARD    int = 0
	BACKWARD   int = 1
	INTERLEAVE int = 1

	//parameters for display activity
	PUSHFROM_OFFSET      int = 8
	PUSHFROM_FINAL_EXTRA int = 5
	PICK_LOAD_POINT      int = 7
	PICK_SWING_SPACE     int = 10 //(in mm)
	TIDY_INITIAL_OFFSET  int = 3
	TIDY_SWING_OFFSET    int = 8
	STACK_START          int = 55
	DICT_SLOT_END        int = 48

	SERVO_PAUSE = 500 // (millisec)

	right bool = true
	left  bool = false

	PREFIX_NEARLY    int = 0
	PREFIX_JUSTAFTER int = 1
	OCLOCKWORD       int = 22

	CMD_REVEAL   = "show"   // Reveal display, if its lying flat for loading/unloading
	CMD_LOAD     = "load"   // Load from dictionary.   i.e. move word onto physical stack
	CMD_UNLOAD   = "unload" // Unload to dictionary    i.e. move word off of physical stack
	CMD_DICTMOVE = "dict"   // for manual testing
	CMD_PICKMOVE = "pick"   // for manual testing - Pick horizontal position
	CMD_ARM      = "arm"    // for manual testing - Angle of Pick arm. Point Left = 180, point up = 90, point right = 0
	CMD_HALT     = "halt"   // for manual testing

)

type (
	ClockCommand struct {
		action string
		value  int
	}

	ClockWord struct {
		text            string
		size            int
		dictionarySteps int
	}

	DisplayWord struct {
		text            string
		size            int
		dictionarySteps int
		stackPos        int
	}
)

// using Map as a Stack as its very small
var (
	wordStack  = make(map[int]DisplayWord)
	words      = [23]ClockWord{}
	ignoreGap  = false
	oldDictPos = 0
	oldPickPos = 0
	removedGap = 0
	totalWords = -1
	wordsShown = false
	led        = machine.LED // Echos the display up & down status for debugging
)

/*
* Build the table of words and their sizes and positions in the dictionary
* If running attached to a physical clock, calibrate also
 */
func initClock() {
	// (In K1 implementation, Size includes 2mm bump)
	//                        Word       Steps Size
	words[0] = makeClockWord("nearly     ", 0, 31)
	words[1] = makeClockWord("just after ", 8, 43)
	words[2] = makeClockWord("five       ", 16, 20)
	words[3] = makeClockWord("ten        ", 24, 18)
	words[4] = makeClockWord("quarter    ", 32, 35)
	words[5] = makeClockWord("twenty     ", 40, 33)
	words[6] = makeClockWord("twentyfive ", 48, 47)
	words[7] = makeClockWord("half       ", 56, 21)
	words[8] = makeClockWord("past       ", 64, 23)
	words[9] = makeClockWord("to         ", 72, 18)
	words[10] = makeClockWord("one       ", 80, 21)
	words[11] = makeClockWord("two       ", 88, 21)
	words[12] = makeClockWord("three     ", 97, 27)
	words[13] = makeClockWord("four      ", 106, 21)
	words[14] = makeClockWord("five      ", 115, 21)
	words[15] = makeClockWord("six       ", 124, 18)
	words[16] = makeClockWord("seven     ", 133, 30)
	words[17] = makeClockWord("eight     ", 142, 26)
	words[18] = makeClockWord("nine      ", 151, 23)
	words[19] = makeClockWord("ten       ", 160, 18)
	words[20] = makeClockWord("eleven    ", 169, 33)
	words[21] = makeClockWord("twelve    ", 178, 33)
	words[22] = makeClockWord("o'clock   ", 187, 35)

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	if HARDWARE_PRESENT {
		calibrate()
	}

}

// Convert three parameters into a ClockWord struct
func makeClockWord(text string, steps int, size int) ClockWord {
	totalWords++
	var w ClockWord
	w.text = strings.TrimSpace(text)
	w.dictionarySteps = steps
	w.size = size
	return w
}

// ==============================================================
func pickFlatLHS() { // pick-arm to travelling position  i.e. turned away from stack by servo
	//servoPick.write(180)
	//delay(servoPause)
}

//==============================================================
/*
 * Move pick arm from flat RHS to vertical but slowly
 */

func readyLoadArm() { // Move from flat LHS to vertical but slowly
	var ix int
	for ix = 135; ix > 95; ix -= 5 {
		slowArm(ix)
	}
}

//==============================================================
/*
 * Move pick arm from flat RHS to vertical but slowly
 */
func readyUnloadArm() {
	var ix int
	for ix = 45; ix < 95; ix += 5 {
		slowArm(ix)
	}
}

// ==============================================================
func slowArm(angle int) {
	delay(100)
	//servoPick.write(angle)
}

// ==============================================================
func pickflatRHS() {
	//servoPick.write(0)
	delay(SERVO_PAUSE)
}

// ==============================================================
func sendDictionaryTo(steps int) {
	var temp int = steps - oldDictPos
	var clockwise bool = true

	if temp < 0 {
		clockwise = false
		temp = -temp
	}

	// fmt.Println("INFO:New Dictionary pos :" + String(steps) )
	// fmt.Println("INFO:Old Dictionary pos :" + String(oldDictPos) )
	// fmt.Println("INFO:steps :" + String(temp))
	// fmt.Println("INFO:fwd :" + String(clockwise) )

	moveDictionary(temp, clockwise)
	oldDictPos = steps
}

//==============================================================
/*
 * "Forward", i.e. down the list of words, is BACKWARD for the Stepper
 */
func moveDictionary(numSteps int, clockwise bool) {
	numSteps += numSteps // Because of interleave
	// if (clockwise){
	//     stepperDict->step(numSteps, BACKWARD, INTERLEAVE)
	// } else {
	//     stepperDict->step(numSteps, FORWARD, INTERLEAVE)
	// }
	// stepperDict->release()
	//fmt.Println("Released Dictionary Stepper")
}

//==============================================================
/*
 * Target for Pick arm is millimeters from LHS
 */
func sendPickTo(target int) {
	var clockwise bool = true

	var temp int = target - oldPickPos

	if temp < 0 {
		clockwise = false
		temp = -temp
	}
	//fmt.Println("INFO:New Pick arm pos :" + String(target) )
	//fmt.Println("INFO:Old Pick arm pos :" + String(oldPickPos) )

	movePickArm(temp, clockwise)
	oldPickPos = target
}

//==============================================================
/*
 * "Forward", i.e. Right to Left, is anti-clockwise and FORWARD for the Stepper
 */
func movePickArm(millimeters int, clockwise bool) {

	//var numSteps int = millimeters * 10 // 1 rotation approx 42 mm.  200 steps per rotation but double for interleave
	//fmt.Println("INFO:Pick arm steps :", numSteps)
	// if (clockwise){
	//     stepperPick->step(numSteps, FORWARD, INTERLEAVE)
	// } else {
	//     stepperPick->step(numSteps, BACKWARD, INTERLEAVE)
	// }

	//stepperPick->release()
	//fmt.Println("Released Pick Stepper")
}

//=========================
/*
* Push a new word onto the stack, set its position on the stack (in mm)
* and reduce removedGap as necessary.
 */
func pushAndSize(w ClockWord) {
	var (
		dw       DisplayWord
		topIndex int
		thatWord DisplayWord
		delta    int
	)

	delta = w.size
	fmt.Printf(" ====== LOAD size : %v %s \n", w.size, w.text)

	if removedGap > 0 {
		fmt.Println("filling gap of :", removedGap)
		if removedGap < w.size {
			delta = w.size - removedGap
			removedGap = 0
		} else {
			delta = 0
			fmt.Println("INFO:removedGap was :", removedGap)
			removedGap -= w.size
			fmt.Println("INFO:removedGap now :", removedGap)
		}
	}

	topIndex = len(wordStack)

	// Update the postions of other words on the display
	for i := 0; i < topIndex; i++ {
		thatWord = wordStack[i]
		thatWord.stackPos += delta
		wordStack[i] = thatWord
	}

	dw.size = w.size
	dw.text = w.text
	dw.stackPos = STACK_START
	wordStack[topIndex] = dw

	/*
	* When multiple words are inserted into a gap big enough, their
	* ".stackPos" can be wrong, but this is corrected by the tidy-up routine.
	* For accurate reporting of stack positions during debug it would need to
	* have the gap held in the internal stack similarly to the way words
	* are held.  Not done yet as it only affects debugging, not the actual
	* operation of the clock.
	 */

}

//=========================
/*
 * Take the next word off the stack
 */
func pop() DisplayWord {

	var topIndex = len(wordStack) - 1

	if topIndex < 0 {
		fmt.Println("Popped but stack already empty - aborting")
		for {
			time.Sleep(1 * time.Second)
		}
	}

	var w = wordStack[topIndex]
	delete(wordStack, topIndex)

	fmt.Printf(" ====== UNLOAD size : %v %s \n", w.size, w.text)
	removedGap += w.size

	// But if stack is now empty, reflect that in removedGap
	if len(wordStack) == 0 {
		removedGap = 0
	}

	return w
}

// =========================
func setDisplayLoading() {
	led.Low()
	wordsShown = false
	//fmt.Println("display down - needs hardware")
}

// =========================
func setDisplayVisible() {
	led.High()
	wordsShown = true
	//fmt.Println("display up - needs hardware")
}

//=========================
/*
 * debugging
 */
func printStack() {
	fmt.Print("----- Stack -----")
	/*var ix int = len(wordStack) - 1
	for ix > -1 {
		var num int = wordStack[ix].wordNum
		fmt.Println("{} ", words[num].text)
		ix -= 1
	}
	fmt.Println()*/
	// Word positions
	fmt.Print("      ")
	var ix = len(wordStack) - 1
	for ix > -1 {
		var dw DisplayWord = wordStack[ix]
		fmt.Printf("[@%d  %s  %d]   ", dw.stackPos, dw.text, dw.stackPos+dw.size)
		ix -= 1
	}
	fmt.Println()
}

//=========================
/*
 * Get the pick to the start position
 */
func calibrate() {
	var switchOpen bool = true
	//fmt.Println("calibrate routine")
	pickFlatLHS()

	//Go forward (right) 5 mm in case we are on the switch already
	movePickArm(5, right)

	//go left until microswitch activates
	for switchOpen {
		movePickArm(2, left)
	}
	//stepperPick->release()
	delay(500)

	// ease off the switch
	movePickArm(3, right)
	//stepperPick->release()
	oldPickPos = 0

}

func delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
