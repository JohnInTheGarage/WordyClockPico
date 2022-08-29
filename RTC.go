package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/ds3231"
)

var (
	sensor   ds3231.Device = ds3231.New(machine.I2C0)
	location time.Location
)

// ===============================
func initRTC() bool {
	/* --------------  Allow time to start monitor
	time.Sleep(time.Second * 20)
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()
	//-----------------------------*/

	machine.I2C0.Configure(machine.I2CConfig{
		SDA: 0,
		SCL: 1,
	})

	sensor.Configure()

	rtcUp := sensor.IsRunning()
	if !rtcUp {
		println("ds3231 not detected")
	}
	return rtcUp
}

/*
* Get current hrs, mins - adjusted for DST.
* No timezone database on microcontrollers! Can't use location, err := time.LoadLocation("Europe/London")
* Some short-cuts are used here :
* Timezone is fixed for UK, so we can use UTC unless Daylight Savings time (i.e. British Summertime) is in force.
* The clock is intended to be external, in the garden, so will only display in daytime in order to
* avoid noise problems for the neighbours.
* This means DST/BST adjustment can be based on date alone, rather than date + time of day
* Which for the UK, is 2am on the final Sunday of March & October
 */

func getTime() (hrs, mins int) {

	rtc, _ := sensor.ReadTime()
	localTime := time.Date(rtc.Year(), rtc.Month(), rtc.Day(), rtc.Hour(), rtc.Minute(), rtc.Second(), rtc.Nanosecond(), time.UTC)
	//debug fmt.Println(" UTC Time : ", localTime)

	dstBegin := finalSunday(localTime.Year(), time.March)
	dstEnd := finalSunday(localTime.Year(), time.October)

	//
	if localTime.After(dstBegin) {
		//debug println("today is after DST start")
		if localTime.Before(dstEnd) {
			//debug println("today is before DST end -add an hour")
			localTime = localTime.Add(time.Hour * 1)
		}
	}

	//debug fmt.Println(" Local time", localTime)

	hrs = localTime.Hour()
	mins = localTime.Minute()
	/*
		rawTemp, _ := sensor.ReadTemperature()
		temp := float32(rawTemp) / 1000

		var all = rtcTime.Weekday()

		fmt.Printf("Temp: %.2f °C, Time %0d:%0d  %d", temp, hrs, mins, all)
		println()
	*/
	return
}

func finalSunday(year int, month time.Month) (testDate time.Time) {

	dd := 31
	var testDay time.Weekday = time.Friday

	for testDay != time.Sunday {
		testDate = time.Date(year, month, dd, 2, 0, 0, 0, time.UTC)
		testDay = testDate.Weekday()
		dd--
	}
	//debug fmt.Println("for ", month, "in ", year, ", Last Sunday is ", testDate)
	return
}
