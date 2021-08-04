package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/muesli/streamdeck"
	"image"
	_ "image/draw"
	"log"
	_ "strconv"
	"image/color"
	"github.com/mdlayher/lmsensors"
)

type FanWidget struct {
	BaseWidget
	device string
	fan0   string
	fan1   string
	fan2   string
	fan3   string
	label0   string
	label1   string
	label2   string
	label3   string
}

func stringColorFan(val bool) color.RGBA {

	colText := color.RGBA{0, 255, 0, 255} // Green
	if val  {
		colText = color.RGBA{255, 0, 0, 255} // Red
	}
	return colText
}


func (w *FanWidget) Update(dev *streamdeck.Device) {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))

	for _, d := range devices {
		for _, s := range d.Sensors {
			if d.Name == w.device {
				switch v := s.(type) {
				//case *lmsensors.CurrentSensor:
				//log.Println(v)
				//case *lmsensors.VoltageSensor:
				//log.Println(v)
				case *lmsensors.FanSensor:
					switch v.Name {
					case w.fan0:
						drawString(img, ttfFont, fmt.Sprintf("%s:%+v", w.label0, v.Input), 6, freetype.Pt(4, 12), stringColorFan(v.Alarm))
					case w.fan1:
						drawString(img, ttfFont, fmt.Sprintf("%s:%+v", w.label1, v.Input), 6, freetype.Pt(4, 24), stringColorFan(v.Alarm))
					case w.fan2:
						drawString(img, ttfFont, fmt.Sprintf("%s:%+v", w.label2, v.Input), 6, freetype.Pt(4, 36), stringColorFan(v.Alarm))
					case w.fan3:
						drawString(img, ttfFont, fmt.Sprintf("%s:%+v", w.label3, v.Input), 6, freetype.Pt(4, 48), stringColorFan(v.Alarm))
					}
				}
			}
		}
	}
	drawString(img, ttfFont, fmt.Sprintf("FANS"), 6, freetype.Pt(-1, 60), color.RGBA{255, 255, 255, 255})

	err := dev.SetImage(w.key, img)
	if err != nil {
		log.Fatal(err)
	}
}
