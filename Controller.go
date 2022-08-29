package main

import "fmt"

var (
	prevWordyHrs    int  = -1
	prevToPast      int  = -1
	prevWordy5Mins  int  = -1
	prevWordyPrefix int  = -1
	prevOclock      bool = false
)

//============================================================

// Convert hours and minutes numbers into the words needed to display current time.
// The display is effectively a stack, with top at the left and bottom at right.
// Most of the time the change required begins with the leftmost word, i.e. the last one
// that was added
// e.g.
// 1) [nearly]        [five] [past] [ten]
// 2) [          ]    [five] [past] [ten]
// 3) [just after]    [five] [past] [ten]
//
// The time display is split into 5 components: prefix, 5mins, toPast, hours and oclock
// Only the hours component is always present, others can be absent.
// Components no longer required for display are removed before adding new ones

func queueCommands(hrs, mins int) {

	var (
		wordyHrs int = hrs
		prefix   int
		toPast   int
		fiveMins int = mins / 5
		fraction int = mins % 5

		exactly   bool = (fraction == 0)
		nearly    bool = (fraction > 2)
		oclock    bool = (mins < 3) || (mins > 57)
		firstHalf bool = (mins < 33)
		cmd       ClockCommand
		////actionNeeded bool
	)

	////actionNeeded = false
	fmt.Printf("\n--- Time %02d:%02d   Queueing...\n", hrs, mins)
	if nearly {
		fiveMins++ // so, if 23 or 24 mins past refer to 25 mins past point, not twenty past
	}
	//fmt.Printf("five Mins before %d\n", fiveMins)
	var fiveMinsWord int
	switch fiveMins {
	case 1:
		fiveMinsWord = 2 //five past
	case 2:
		fiveMinsWord = 3 //ten past
	case 3:
		fiveMinsWord = 4 //quarter
	case 4:
		fiveMinsWord = 5 //twenty
	case 5:
		fiveMinsWord = 6 //twentyfive
	case 12:
		fiveMinsWord = -1 //At the hour
	case 11:
		fiveMinsWord = 2 //five to
	case 10:
		fiveMinsWord = 3 //ten to
	case 9:
		fiveMinsWord = 4 //quarter
	case 8:
		fiveMinsWord = 5 //twenty
	case 7:
		fiveMinsWord = 6 //twentyfive
	default:
		fiveMinsWord = 7 //half
	}
	//fmt.Printf("five Mins after %d\n", fiveMinsWord)

	if exactly {
		prefix = -1
	} else {
		if nearly {
			prefix = PREFIX_NEARLY
		} else {
			prefix = PREFIX_JUSTAFTER
		}

	}

	if !firstHalf {
		wordyHrs++
	}
	if wordyHrs > 12 {
		wordyHrs = wordyHrs - 12
	}
	if wordyHrs == 0 {
		wordyHrs = 12
	}
	wordyHrs += 9 // because the hours start at array index 10

	if oclock {
		fiveMinsWord = -1
		toPast = -1
	} else {
		if firstHalf {
			toPast = 8
		} else {
			toPast = 9
		}
	}

	// Do all the unstacking before considering the stacking or re-stacking.
	// ------------------ do we need to unload anything from the stack? -----------------------
	if prevWordyPrefix != prefix {
		if prevWordyPrefix != -1 { // previous prefix might not be there when its value changes eg just after the hour
			insertUnload(prevWordyPrefix)
			//actionNeeded = true
		}
	}
	if prevWordy5Mins != fiveMinsWord {
		if prevWordy5Mins != -1 {
			insertUnload(prevWordy5Mins)
			//actionNeeded = true
		}
	}
	if prevToPast != toPast {
		if prevToPast != -1 {
			insertUnload(prevToPast)
			//actionNeeded = true
		}
	}

	// Hour has changed.  BTW, this is not when the minute hand reaches 12
	if prevWordyHrs != wordyHrs {
		if prevWordyHrs != -1 {
			insertUnload(prevWordyHrs) //hours are always present (except when starting...)
			//actionNeeded = true
		}
	}

	// Minute hand nearing 12.  Time to include/exclude the "O'Clock"
	if prevOclock != oclock {
		//actionNeeded = true
		if prevWordyHrs != -1 {
			insertUnload(prevWordyHrs) //When o'clock changes the hour does not (surprise!) so hour is re-loaded below
		}
		// only remove oClock word if no longer needed
		if prevOclock {
			insertUnload(OCLOCKWORD)
		}
	}

	// -------- Anything to load onto the stack? ----------

	if prevOclock != oclock {
		if oclock {
			insertLoad(OCLOCKWORD)
			////actionNeeded = true
		}
		insertLoad(wordyHrs) //When o'clock changes the hour word does not change so is re-loaded
		//actionNeeded = true
		prevWordyHrs = wordyHrs //BTW, the hour word changes at 33 mins past the hour...
		prevOclock = oclock
	}

	// hours word is always required, no test for -1
	if prevWordyHrs != wordyHrs {
		insertLoad(wordyHrs)
		//actionNeeded = true
		// moved prevWordyHrs = wordyHrs
	}

	// toPast word can be absent
	if toPast > -1 {
		if prevToPast != toPast {
			insertLoad(toPast)
			//actionNeeded = true
		}
		//moved prevToPast = toPast
	}

	// fiveMinsWord word can be absent
	if fiveMinsWord > -1 {
		if prevWordy5Mins != fiveMinsWord {
			insertLoad(fiveMinsWord)
			//actionNeeded = true
		}
		// moved prevWordy5Mins = fiveMinsWord
	}

	// prefix word can be absent
	if prefix > -1 {
		if prevWordyPrefix != prefix {
			insertLoad(prefix)
			//actionNeeded = true
		}
		// moved prevWordyPrefix = prefix
	}

	//if actionNeeded {
	cmd.action = CMD_REVEAL
	cmd.value = 1
	insertCommand(cmd)
	prevWordyPrefix = prefix
	prevWordy5Mins = fiveMinsWord
	prevToPast = toPast
	prevWordyHrs = wordyHrs
	prevOclock = oclock
	//}

	/*
	 * If tidy is required then we will calculate the gap to close based
	 * on size of words removed & added
	 */
}

// Queue a load command
func insertLoad(wordNum int) {
	//if wordNum > -1 {
	var cmd ClockCommand
	cmd.action = CMD_LOAD
	cmd.value = wordNum
	insertCommand(cmd)
	// } else {
	// 	fmt.Println("skipping load -1 word number")
	// }
}

// Queue an unload command
func insertUnload(wordNum int) {
	// if wordNum > -1 {
	var cmd ClockCommand
	cmd.action = CMD_UNLOAD
	cmd.value = wordNum
	insertCommand(cmd)
	// } else {
	// 	fmt.Println("skipping unload -1 word number")
	// }
}

func insertCommand(cmd ClockCommand) {
	/*
		fmt.Printf("Command :%s ", cmd.action)
		if cmd.action == CMD_LOAD || cmd.action == CMD_UNLOAD {
			fmt.Print(words[cmd.value].text)
		}
		fmt.Println()
	*/
	cmdChannel <- cmd
}
