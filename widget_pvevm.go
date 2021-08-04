package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/muesli/streamdeck"
	"image"
	"image/color"
	_ "image/draw"
	"log"
	"os/exec"
	"strconv"
	_ "strconv"
	"sync"
	"time"
	"strings"
)

type PveVMWidget struct {
	BaseWidget

	vmid           string
	vmlabel        string
	doublearm      string
	updateinterval string

	armed       bool
	doublearmed bool
	status      string
	init        sync.Once
	updatetime  time.Time
	_updateinterval time.Duration
	_armedAction string
	_doublearm	bool
}

func proxmoxVmCtrl(vm string, command string) {
	_, err := exec.Command("/usr/local/bin/vmctrl.sh", vm, command).Output()
	if err != nil {
		log.Fatal(err)
	}

}

func proxmoxStatus(vm string) string {
	out, err := exec.Command("/usr/local/bin/vmcheck.sh", vm).Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func stringColorVM(vmstatus string) color.RGBA {
	colText := color.RGBA{255, 255, 255, 255} // White

	switch vmstatus {
	case "running":
		colText = color.RGBA{0, 255, 0, 255} // Green
	case "stopped":
		colText = color.RGBA{255, 0, 0, 255} // Red
	case "paused":
		colText = color.RGBA{255, 128, 0, 255} // Yellow

	}
	return colText
}

func (w *PveVMWidget) TriggerAction(hold bool) {
	//log.Println("PVEVM Action: ", hold)
	w.updatetime = time.Now()
	if hold {
		if (w.armed){
			if ((w._doublearm && w.doublearmed) || (!w._doublearm)){
				switch (w.status){
				case "running":
					w._armedAction = "stopping.."
					proxmoxVmCtrl(w.vmid, "stop")
				case "stopped":
					w._armedAction = "starting.."
					proxmoxVmCtrl(w.vmid, "start")
				case "paused":
					w._armedAction = "resuming.."
					proxmoxVmCtrl(w.vmid, "resume")
				}
				w.armed = false
				w.doublearmed = false
				w._armedAction = ""
			}else{
				w.doublearmed = true
			}
		}else{
			w.armed = true
			if (w._doublearm){
				w._armedAction = "hold again"
			}else{
				switch (w.status){
				case "running":
					w._armedAction = "hold to stop"
				case "stopped":
					w._armedAction = "hold to start"
				case "paused":
					w._armedAction = "hold to resume"
				}

			}
		}
	}else{
		//TODO - DIsplay VM Graphs ?
	}
}

func (w *PveVMWidget) Update(dev *streamdeck.Device) {
	if w.updatetime.Before(time.Now()) {
		//log.Println("PVEVM Update ", w.vmlabel)
		img := image.NewRGBA(image.Rect(0, 0, 72, 72))
		oldstatus := w.status
		w.status = proxmoxStatus(w.vmid)
		if (oldstatus != w.status){
			w._armedAction=""
		}
		drawString(img, ttfFont, fmt.Sprintf("%s",w.vmlabel), 10, freetype.Pt(-1, 16), color.RGBA{255, 255, 255, 255})
		drawString(img, ttfFont, fmt.Sprintf("%s",w.status), 8, freetype.Pt(-1, 32), stringColorVM(w.status))
		drawString(img, ttfFont, fmt.Sprintf("%s",w._armedAction), 6, freetype.Pt(-1, 48), stringColorVM(w.status))
		drawString(img, ttfFont, fmt.Sprintf("VM"), 6, freetype.Pt(-1, 62), color.RGBA{255, 255, 255, 255})
		err := dev.SetImage(w.key, img)
		if err != nil {
			log.Fatal(err)
		}
		w.updatetime = time.Now().Add(w._updateinterval * time.Second)

	}
	w.init.Do(func() {
		//log.Println("PVEVM Init")
		w.status = proxmoxStatus(w.vmid)
		ui, _ := strconv.Atoi(w.updateinterval)
		w._updateinterval = time.Duration(ui)
		w.updatetime = time.Now().Add(w._updateinterval * time.Second)
		w._armedAction = ""
		w._doublearm = false
		if (w.doublearm == "true"){
			w._doublearm = true
		}
	})
}
