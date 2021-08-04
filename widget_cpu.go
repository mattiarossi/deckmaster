package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/muesli/streamdeck"
	"image"
	"image/color"
	_ "image/draw"
	"log"
	"reflect"
	_ "strconv"
)

type CpuWidget struct {
	BaseWidget
	cpu0 string
	cpu1 string
	cpu2 string
	cpu3 string
}

func getAttr(obj interface{}, fieldName string) reflect.Value {
	pointToStruct := reflect.ValueOf(obj) // addressable
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(fieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField
}

func stringColor(val float32) color.RGBA {

	colText := color.RGBA{0, 255, 0, 255} // Green
	if val >= 0.5 {
		colText = color.RGBA{255, 128, 0, 255} // Yellow
	}
	if val >= 1.0 {
		colText = color.RGBA{255, 0, 0, 255} // Red
	}
	return colText
}

func (w *CpuWidget) Update(dev *streamdeck.Device) {
	img := image.NewRGBA(image.Rect(0, 0, 72, 72))

  cpu0 := float32(getAttr(&coreStats, w.cpu0).Float())
	drawString(img, ttfFont, fmt.Sprintf("%s:%.3f", w.cpu0, cpu0), 6, freetype.Pt(4, 12), stringColor(cpu0))
  cpu1 := float32(getAttr(&coreStats, w.cpu1).Float())
	drawString(img, ttfFont, fmt.Sprintf("%s:%.3f", w.cpu1, cpu1), 6, freetype.Pt(4, 24), stringColor(cpu1))
  cpu2 := float32(getAttr(&coreStats, w.cpu2).Float())
	drawString(img, ttfFont, fmt.Sprintf("%s:%.3f", w.cpu2, cpu2), 6, freetype.Pt(4, 36), stringColor(cpu2))
  cpu3 := float32(getAttr(&coreStats, w.cpu3).Float())
	drawString(img, ttfFont, fmt.Sprintf("%s:%.3f", w.cpu3, cpu3), 6, freetype.Pt(4, 48), stringColor(cpu3))
	err := dev.SetImage(w.key, img)
	if err != nil {
		log.Fatal(err)
	}
}
