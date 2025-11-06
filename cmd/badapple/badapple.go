package main

import (
	"GBSandbox/pkg/comms"
	"GBSandbox/pkg/imgutil"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var data = comms.GBData{
	Header: comms.GamebandHeader{
		Timezone:                5,
		AltTimezone:             0,
		TzChange:                0,
		Orientation:             1,
		TransitionFrameDuration: 47,
		ScreenCount:             1,
	},
	//Animations: nil,
	Animations: []comms.Animation{
		{
			Header: comms.AnimationHeader{
				ScreenType:    32,
				PauseMode:     0,
				PauseDuration: 1000,
				FrameDuration: 250,
				AnimationType: 0,
				DataLength:    0,
			},
		},
	},
}

func main() {
	frames := make([]*os.File, 160)

	c := 0
	err := filepath.WalkDir("cmd/badapple/frames", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		frames[c] = f
		c++
		return nil
	})

	if err != nil {
		panic(err)
	}

	gb, err := comms.OpenHid()
	if err != nil {
		panic(err)
	}

	gbFrames := make([]comms.Frame, 100)

	for i := range frames {
		if i == len(gbFrames) {
			break
		}
		img, _, err := image.Decode(frames[i])
		if err != nil {
			continue
		}
		tempImg := image.NewGray(image.Rect(0, 0, 20, 7))
		draw.Draw(tempImg, img.Bounds(), img, image.Point{}, draw.Over)

		gbFrame := comms.Frame{}
		gbFrame.Data = make([]byte, 20)
		for i := 0; i < 20; i += 2 {
			gbFrame.Data[i] = byte(imgutil.GetTwoColumns(tempImg, i))
			gbFrame.Data[i+1] = byte(imgutil.GetTwoColumns(tempImg, i) >> 8)
		}
		gbFrames[i] = gbFrame
	}
	data.Animations[0].Frames = gbFrames
	fmt.Printf("Frame Count: %d", len(data.Animations[0].Frames))

	err = gb.WriteGBData(data)
	if err != nil {
		panic(err)
	}

	for _, frame := range frames {
		frame.Close()
	}
}
