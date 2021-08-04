package main

import (
	"image"
	"log"
	"strconv"
	"time"
	"image/color"
	"github.com/golang/freetype"
	"github.com/muesli/streamdeck"
)

type DateWidget struct {
	BaseWidget
}

func (w *DateWidget) Update(dev *streamdeck.Device) {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))

	t := time.Now()
	day := t.Day()
	month := t.Month().String()
	year := t.Year()
	drawString(img, ttfBoldFont, strconv.Itoa(day), 13, freetype.Pt(-1, 20), color.RGBA{255, 255, 255, 255})
	drawString(img, ttfFont, month, 13, freetype.Pt(-1, 43), color.RGBA{255, 255, 255, 255})
	drawString(img, ttfThinFont, strconv.Itoa(year), 13, freetype.Pt(-1, 66), color.RGBA{255, 255, 255, 255})

	err := dev.SetImage(w.key, img)
	if err != nil {
		log.Fatal(err)
	}
}
