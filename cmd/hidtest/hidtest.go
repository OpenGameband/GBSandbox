package main

import (
	"GBSandbox/pkg/comms"
	"GBSandbox/pkg/imgutil"
	"fmt"
	"image"
	_ "image/gif"
	"os"
	"time"
)

var data = comms.GBData{
	Header: comms.GamebandHeader{
		Timezone:                5,
		AltTimezone:             0,
		TzChange:                0,
		Orientation:             1,
		TransitionFrameDuration: 47,
		ScreenCount:             3,
		AnimationDataLength:     12,
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
		{
			Header: comms.AnimationHeader{
				ScreenType:    32,
				PauseMode:     0,
				PauseDuration: 1428,
				FrameDuration: 1428,
				AnimationType: 0,
				DataLength:    0,
			},
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

	z, offset := time.Now().Zone()
	data.Header.Timezone = uint8((offset / 60 / 60) * 4)
	fmt.Printf("Setting Gameband Timezone: %s (%d)\n", z, offset/60/60)
	if _, end := time.Now().ZoneBounds(); end.After(time.Now()) {
		data.Header.TzChange = uint32(end.Unix())
		z, offset := end.Zone()
		fmt.Printf("Setting Gameband Next Timezone: %s (%d)\n", z, offset)
		data.Header.AltTimezone = uint8((offset / 60 / 60) * 4)
	}

	if len(os.Args) > 1 && os.Args[1] != "" {
		imageFile, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer imageFile.Close()
		img, _, err := image.Decode(imageFile)
		if err != nil {
			panic(err)
		}

		frame := comms.Frame{}
		frame.Data = make([]byte, 20)
		for i := 0; i < 20; i += 2 {
			frame.Data[i] = byte(imgutil.GetTwoColumns(img, i))
			frame.Data[i+1] = byte(imgutil.GetTwoColumns(img, i) >> 8)
		}
		data.Animations[2].Frames = []comms.Frame{frame}
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
	err = os.WriteFile("gb-"+time.Now().Format(time.DateTime)+".dump", data, 0644)
	if err != nil {
		panic(err)
	}

}
