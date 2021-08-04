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

type TempWidget struct {
	BaseWidget
	device string
	temp0   string
	temp1   string
	temp2   string
	temp3   string
	label0   string
	label1   string
	label2   string
	label3   string
}

func stringColorTemp(val float64, alarm bool, high float64) color.RGBA {

	colText := color.RGBA{0, 255, 0, 255} // Green
	if (val > high){
		colText = color.RGBA{255, 255, 0, 0} // Yellow
	}
	if alarm  {
		colText = color.RGBA{255, 0, 0, 255} // Red
	}
	return colText
}


func (w *TempWidget) Update(dev *streamdeck.Device) {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))

	for _, d := range devices {
		for _, s := range d.Sensors {
			if d.Name == w.device {
				switch v := s.(type) {
				case *lmsensors.TemperatureSensor:
					switch v.Name {
					case w.temp0:
						drawString(img, ttfFont, fmt.Sprintf("%s:%.1fC", w.label0, v.Input), 6, freetype.Pt(4, 12), stringColorTemp(v.Input, v.Alarm, v.High))
					case w.temp1:
						drawString(img, ttfFont, fmt.Sprintf("%s:%.1fC", w.label1, v.Input), 6, freetype.Pt(4, 24), stringColorTemp(v.Input, v.Alarm, v.High))
					case w.temp2:
						drawString(img, ttfFont, fmt.Sprintf("%s:%.1fC", w.label2, v.Input), 6, freetype.Pt(4, 36), stringColorTemp(v.Input, v.Alarm, v.High))
					case w.temp3:
						drawString(img, ttfFont, fmt.Sprintf("%s:%.1fC", w.label3, v.Input), 6, freetype.Pt(4, 48), stringColorTemp(v.Input, v.Alarm, v.High))
					}
				}
			}
		}
	}
	drawString(img, ttfFont, fmt.Sprintf("TEMPS"), 6, freetype.Pt(-1, 60), color.RGBA{255, 255, 255, 255})

	err := dev.SetImage(w.key, img)
	if err != nil {
		log.Fatal(err)
	}
}
