package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"strconv"

	"github.com/golang/freetype"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/streamdeck"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type TopWidget struct {
	BaseWidget
	mode      string
	fillColor string

	lastValue float64
}

func (w *TopWidget) Update(dev *streamdeck.Device) {
	var value float64
	var label string
	switch w.mode {
	case "cpu":
		cpuUsage, err := cpu.Percent(0, false)
		if err != nil {
			log.Fatal(err)
		}

		value = cpuUsage[0]
		label = "CPU"

	case "memory":
		mem, err := mem.VirtualMemory()
		if err != nil {
			log.Fatal(err)
		}
		value = mem.UsedPercent
		label = "MEM"

	default:
		panic("Unknown widget mode: " + w.mode)
	}

	if w.lastValue == value {
		return
	}
	w.lastValue = value

	fill, err := colorful.Hex(w.fillColor)
	if err != nil {
		panic("Invalid color: " + w.fillColor)
	}

	img := image.NewRGBA(image.Rect(0, 0, 72, 72))
	draw.Draw(img,
		image.Rect(12, 6, 60, 54),
		&image.Uniform{color.RGBA{255, 255, 255, 255}},
		image.Point{}, draw.Src)
	draw.Draw(img,
		image.Rect(13, 7, 59, 53),
		&image.Uniform{color.RGBA{0, 0, 0, 255}},
		image.Point{}, draw.Src)
	draw.Draw(img,
		image.Rect(14, 7+int(46*(1-value/100)), 58, 53),
		&image.Uniform{fill},
		image.Point{}, draw.Src)

	drawString(img, ttfFont, strconv.FormatInt(int64(value), 10), 12, freetype.Pt(-1, -1), color.RGBA{255, 255, 255, 255})
	drawString(img, ttfFont, "% "+label, 7, freetype.Pt(-1, img.Bounds().Dx()-4), color.RGBA{255, 255, 255, 255})

	err = dev.SetImage(w.key, img)
	if err != nil {
		log.Fatal(err)
	}
}
