// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/mfrc522"
)

func main() {

	var (
		rstPin machine.Pin = machine.GP10
		irqPin machine.Pin = machine.GP11
	)
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 10,
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN, //Miso
		SDI:       machine.SPI0_SDI_PIN, //Mosi
		Mode:      0})

	rfid, err := mfrc522.NewSPI(machine.SPI0, rstPin, irqPin)
	if err != nil {
		println(err.Error())
	}

	// Idling device on exit.
	defer rfid.Halt()

	// Setting the antenna signal strength.
	rfid.SetAntennaGain(5)

	timedOut := false
	cb := make(chan []byte)
	timer := time.NewTimer(10 * time.Second)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		timer.Stop()
		timedOut = true
		close(cb)
	}()

	go func() {
		fmt.Println("Started %s", rfid.String())

		for {
			// Trying to read card UID.
			uid, err := rfid.ReadUID(10 * time.Second)

			// If main thread timed out just exiting.
			if timedOut {
				return
			}

			// Some devices tend to send wrong data while RFID chip is already detected
			// but still "too far" from a receiver.
			// Especially some cheap CN clones which you can find on GearBest, AliExpress, etc.
			// This will suppress such errors.
			if err != nil {
				continue
			}

			cb <- uid
			return
		}
	}()

	for {
		select {
		case <-timer.C:
			fmt.Println("Didn't receive any data")
			return
		case data := <-cb:
			fmt.Println("UID:", hex.EncodeToString(data))
			return
		}
	}
}
