package main

import (
	"fmt"
	"machine"
	"time"
)

const (
	PROD_NUM_ADDR0 = 0x0000
	PROD_NUM_ADDR1 = 0x0001

	PROD_IQS550_VAL = 40
	PROD_IQS572_VAL = 58
	PROD_IQS525_VAL = 52

	PROJECT_NUM_ADDR0 = 0x0002
	PROJECT_NUM_ADDR1 = 0x0003
	MAJ_VER_ADDR      = 0x0004
	MIN_VER_ADDR      = 0x0005
	BOOTLOADER_ADDR   = 0x0006

	BOOTLOADER_IS_AVAIL_VAL = 0xA5
	NO_BOOTLOADER_VAL       = 0xEE

	// 0-3 bits for row and 4-7 bits for column
	MAX_TOUCH_ADDR       = 0x000B
	PREV_CYCLE_TIME_ADDR = 0x000C

	// 0 - single tap
	// 1 - press and hold
	// 2 - swipe negative X
	// 3 - swipe positive X
	// 4 - swipe negative Y
	// 5 - swipe positive Y
	// 6 - not used
	// 7 - not used
	GESTURE_ADDR0 = 0x0D
	// 0 - double tap
	// 1 - scroll
	// 2 - zoom
	// 3 - not used
	// 4 - not used
	// 5 - not used
	// 6 - not used
	// 7 - not used
	GESTURE_ADDR1 = 0x0E

	SYSTEM_INFO_ADDR0            = 0x0F
	SYSTEM_INFO_ADDR1            = 0x10
	NUM_OF_FINGERS_ADDR          = 0x11
	FINGER0_REL_X_ADDR0          = 0x12
	FINGER0_REL_X_ADDR1          = 0x13
	FINGER0_REL_Y_ADDR0          = 0x14
	FINGER0_REL_Y_ADDR1          = 0x15
	FINGER0_ABS_X_ADDR0          = 0x16
	FINGER0_ABS_X_ADDR1          = 0x17
	FINGER0_ABS_Y_ADDR0          = 0x18
	FINGER0_ABS_Y_ADDR1          = 0x19
	FINGER0_TOUCH_STRENGTH_ADDR0 = 0x1A
	FINGER0_TOUCH_STRENGTH_ADDR1 = 0x1B
	FINGER0_TOUCH_AREA_SIZE_ADDR = 0x1C
	FINGER1_ABS_X_ADDR0          = 0x1D
	FINGER1_ABS_X_ADDR1          = 0x1E
	FINGER1_ABS_Y_ADDR0          = 0x1F
	FINGER1_ABS_Y_ADDR1          = 0x20
	FINGER1_TOUCH_STRENGTH_ADDR0 = 0x21
	FINGER1_TOUCH_STRENGTH_ADDR1 = 0x22
	FINGER1_TOUCH_AREA_SIZE_ADDR = 0x23
	FINGER2_ABS_X_ADDR0          = 0x24
	FINGER2_ABS_X_ADDR1          = 0x25
	FINGER2_ABS_Y_ADDR0          = 0x26
	FINGER2_ABS_Y_ADDR1          = 0x27
	FINGER2_TOUCH_STRENGTH_ADDR0 = 0x28
	FINGER2_TOUCH_STRENGTH_ADDR1 = 0x29
	FINGER2_TOUCH_AREA_SIZE_ADDR = 0x2A
	FINGER3_ABS_X_ADDR0          = 0x2B
	FINGER3_ABS_X_ADDR1          = 0x2C
	FINGER3_ABS_Y_ADDR0          = 0x2D
	FINGER3_ABS_Y_ADDR1          = 0x2E
	FINGER3_TOUCH_STRENGTH_ADDR0 = 0x2F
	FINGER3_TOUCH_STRENGTH_ADDR1 = 0x30
	FINGER3_TOUCH_AREA_SIZE_ADDR = 0x31
	FINGER4_ABS_X_ADDR0          = 0x32
	FINGER4_ABS_X_ADDR1          = 0x33
	FINGER4_ABS_Y_ADDR0          = 0x34
	FINGER4_ABS_Y_ADDR1          = 0x35
	FINGER4_TOUCH_STRENGTH_ADDR0 = 0x36
	FINGER4_TOUCH_STRENGTH_ADDR1 = 0x37
	FINGER4_TOUCH_AREA_SIZE_ADDR = 0x38

	COMM_ADDR = 0x74

	READ  = 0xE9
	WRITE = 0xE8

	END_COMM_ADDR = 0xEE
)

func main() {
	machine.InitSerial()
	time.Sleep(5 * time.Second)
	readyPin := machine.P0_06

	readyPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	for readyPin.Get() {
		fmt.Println("Waiting for ready")
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Ready")
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.SCL_PIN,
		SDA: machine.SDA_PIN,
	})
	if err != nil {
		fmt.Println("failed to configure i2c")
		fmt.Println(err)
	}

	i2c.Tx(COMM_ADDR, []byte{0x00}, nil)
	i2c.Tx(COMM_ADDR, []byte{GESTURE_ADDR0}, nil)
	fmt.Println("reading initial data")
	var data = make([]byte, 16)
	i2c.Tx(COMM_ADDR, nil, data)
	fmt.Printf("initial data: %08b\n", data)
	for {
		var data = make([]byte, 1)
		i2c.Tx(COMM_ADDR, []byte{0x00, GESTURE_ADDR0}, data)

		// As we can have press and hold combined with other gestures, we want to check it out of switch
		var pressHold = uint8(data[0])>>1&0x01 == 0x01
		if pressHold {
			fmt.Println("press and hold")
		}

		switch {
		case uint8(data[0])&0x01 == 0x01:
			fmt.Println("tapped")
		case uint8(data[0])>>2&0x01 == 0x01:
			fmt.Println("swipe down")
		case uint8(data[0])>>3&0x01 == 0x01:
			fmt.Println("swiped up")
		case uint8(data[0])>>4&0x01 == 0x01:
			fmt.Println("swiped right")
		case uint8(data[0])>>5&0x01 == 0x01:
			fmt.Println("swiped left")
		default:
		}

		i2c.Tx(COMM_ADDR, []byte{0x00, GESTURE_ADDR1}, data)
		switch {
		case uint8(data[0])&0x01 == 0x01:
			fmt.Println("2 finger tapped")
		case uint8(data[0])>>1&0x01 == 0x01:
			fmt.Println("scroll")
		case uint8(data[0])>>2&0x01 == 0x01:
			fmt.Println("zoom")
		default:
		}
		i2c.Tx(COMM_ADDR, []byte{0x00, SYSTEM_INFO_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, SYSTEM_INFO_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, NUM_OF_FINGERS_ADDR}, data)
		// if uint8(data[0]) > 1 {
		// 	fmt.Printf("more than one finger: %d\n", int(data[0]))
		// }

		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_REL_X_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_REL_X_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_REL_Y_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_REL_Y_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_ABS_X_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_REL_X_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_ABS_Y_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_ABS_Y_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_TOUCH_STRENGTH_ADDR0}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_TOUCH_STRENGTH_ADDR1}, data)
		i2c.Tx(COMM_ADDR, []byte{0x00, FINGER0_TOUCH_AREA_SIZE_ADDR}, data)

		i2c.Tx(COMM_ADDR, []byte{END_COMM_ADDR, END_COMM_ADDR}, nil)
		i2c.Tx(COMM_ADDR, []byte{0x00}, nil)

		// i2c.Tx(COMM_ADDR, []byte{0x00}, nil)
		// i2c.WriteRegister(COMM_ADDR, 0x00, nil)
		// fmt.Printf("data: %08b\n", data) // (data)
		time.Sleep(30 * time.Millisecond)
	}

}
