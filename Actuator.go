package main

import (
	"fmt"
	"strings"
)

// Receives commands from clockController and drives the clock motors
func ClockActuator(cmds chan ClockCommand) {
	var numValue int

	for {
		item := <-cmds
		numValue = item.value
		debugDelay()
		switch item.action {
		case CMD_LOAD:
			loading(numValue)
		case CMD_UNLOAD:
			unloading(numValue)
		case CMD_REVEAL:
			reveal()

			// ---------- Initial calibration for prototype build --------
			// Expected to be via MANUAL input so OK not sent back to sender
		case CMD_PICKMOVE:
			sendPickTo(numValue)
		case CMD_DICTMOVE:
			sendDictionaryTo(words[numValue].dictionarySteps)
		case CMD_ARM:
			//servoPick.write(numValue);
		case CMD_HALT:
			for {
				delay(1000)
			}
		}
	}

}

//==============================================================
/*
* Fetch the supplied word and put it on the stacks
* (both virtual and real  i.e. the display)
 */
func loading(wordNum int) {
	if wordsShown {
		setDisplayLoading()
	}

	var activeWord = words[wordNum]
	pushAndSize(activeWord)
	pickFlatLHS()
	sendDictionaryTo(activeWord.dictionarySteps)
	sendPickTo(PICK_LOAD_POINT)
	readyLoadArm()
	// Push word onto stack until pick leaves dictionary
	sendPickTo(STACK_START)
	// Back up the pick for 10mm, leaving room to swing arm either way
	sendPickTo(STACK_START - PICK_SWING_SPACE) // leave room to swing arm left or right

	//fmt.Println("OK:loaded " + activeWord.text)

}

// ==============================================================
func unloading(wordNum int) {
	if wordsShown {
		setDisplayLoading()
	}

	var pushFrom int
	var activeWord = pop()
	ignoreGap = (strings.ToLower(activeWord.text) == "nearly")
	pushFrom = activeWord.stackPos + activeWord.size - PUSHFROM_OFFSET
	/*
	 * When unloading the last word from the stack go further right
	 * as vibration can move the final word
	 */
	if len(wordStack) == 0 {
		pushFrom += PUSHFROM_FINAL_EXTRA
	}
	/*
		fmt.Println("INFO: ====== UNLOAD ===========")
		fmt.Println("INFO:word at   :", activeWord.stackPos)
		fmt.Println("INFO:word size :", activeWord.size)
		fmt.Println("INFO:push from :", pushFrom)
		fmt.Println("INFO:gap size  :", removedGap)
	*/

	pickflatRHS() //load routine has to leave space to swing arm to the right side
	sendPickTo(pushFrom)
	sendDictionaryTo(activeWord.dictionarySteps)
	readyUnloadArm()
	/* push the word as far in as it will go
	 * to stop vibration moving it rightwards and into contact with
	 * stationary parts when revolving the dictionary
	 */
	sendPickTo(DICT_SLOT_END - (DICT_SLOT_END - activeWord.size))
	sendPickTo(STACK_START) // i.e.until pick leaves dictionary so it can rotate

	//fmt.Println("OK:un-loaded " + activeWord.text)

}

/*
 * The "Reveal" code is sent at the change of each minute but often there is
 * no change in the display, so the words may already be showing.
 */
func reveal() {
	if !wordsShown {
		tidyUp()
		setDisplayVisible()
		printStack()
	}
}

//==============================================================
/*
 * The original tidy-up took a long time, so the new idea is
 * to delay tidy-up until after the loads, on the basis that the gap
 * caused by removing words is often filled, at least partly, by
 * the next words to be loaded.
 * So "removedGap" is increased during unloading and decreased
 * during loading (although if we unload ALL the words it becomes zero).
 * The Pi tells us when loading is finished by sending a "!"
 * if removedGap is > zero at that time, then we can tidy-up.
 * The scheme of avoiding tidyup after removing NEARLY still applies
 * because the next state (<exact>-time), has no load, and then
 * loading JUSTAFTER over-fills the gap left by NEARLY.
 */

func tidyUp() {
	if ignoreGap {
		return
	}

	// Calc extent of words on stack
	var stackMax = STACK_START + removedGap
	for _, item := range wordStack {
		stackMax += item.size
	}

	if removedGap > 0 {
		fmt.Println("INFO: ====== TIDY UP ===========")
		fmt.Println("INFO:stackMax   :", stackMax)

		pickflatRHS()
		sendPickTo(stackMax - TIDY_INITIAL_OFFSET)
		readyUnloadArm()
		stackMax -= removedGap

		//close the gap
		sendPickTo(stackMax)
		pickflatRHS()
		// allow arm-swing room
		sendPickTo(stackMax + TIDY_SWING_OFFSET)
		fmt.Println("INFO:Removed Gap :", removedGap)
	}
	removedGap = 0

	// Stack positions updated when loading finished
	stackMax = STACK_START
	for _, item := range wordStack {
		item.stackPos = stackMax
		stackMax += item.size
	}
	/*
		ix = len(wordStack) - 1
		for ix > -1 {
			var thisWord DisplayWord = wordStack[ix]
			thisWord.stackPos = stackMax
			stackMax += thisWord.size
			wordStack[ix] = thisWord
			//msg = fmt.Sprintf("INFO: stack(%d) pos:%d size:%d",
			//	ix, wordstack[ix].stackPos, wordstack[ix].size)
			//fmt.Println(msg)
			ix -= 1
		}*/
	fmt.Printf("------- Display tidied, stackMax :%d,   removedGap :%d", stackMax, removedGap)
	fmt.Println()

}

// ==============================================================

func debugDelay() {
	delay(1000)
}
