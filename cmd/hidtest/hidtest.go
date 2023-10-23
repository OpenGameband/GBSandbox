package main

import (
	"GBSandbox/pkg/comms"
	"fmt"
	"os"
	"time"
)

var data = comms.GBData{
	Header: comms.GamebandHeader{
		Timezone:                4,
		AltTimezone:             0,
		TzChange:                1698541200,
		Orientation:             1,
		TransitionFrameDuration: 47,
		ScreenCount:             2,
		AnimationDataLength:     12,
		Checksum0:               143,
		Checksum1:               152,
	},
	//Animations: nil,
	Animations: []comms.Animation{
		{
			Header: comms.AnimationHeader{
				ScreenType:    0,
				PauseMode:     0,
				PauseDuration: 0,
				FrameDuration: 3000,
				AnimationType: 3,
				DataLength:    0,
			},
			Frames: nil,
		},
		{
			Header: comms.AnimationHeader{
				ScreenType:    2,
				PauseMode:     0,
				PauseDuration: 0,
				FrameDuration: 3000,
				AnimationType: 3,
				DataLength:    0,
			},
			Frames: nil,
		},
	},
}

func main() {
	gb := new(comms.Gameband)
	gb, err := comms.OpenHid()
	if err != nil {
		panic(err)
	}
	fmt.Println("Setting Gameband Time:")
	err = gb.SetTime()
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing Gameband Data")
	err = gb.WriteGBData(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Finishing")
	err = gb.Commit()
	if err != nil {
		panic(err)
	}

	fmt.Println("Dumping Gameband flash...")
	data, err := gb.ReadGameband()
	err = os.WriteFile("gb-"+time.Now().String()+".dump", data, 0644)
	if err != nil {
		panic(err)
	}
}
