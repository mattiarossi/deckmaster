package main

import (
	_ "fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/mdlayher/lmsensors"
	"github.com/muesli/streamdeck"
	"log"
	"time"
)

type StatsWidget struct {
	BaseWidget
}

//
// cpu usage values are read from the /proc/stat pseudo-file with the help of the goprocinfo package...
// For the calculation two measurements are neaded: 'current' and 'previous'...
// More at: func calcSingleCoreUsage(curr, prev)...
//
type MyCPUStats struct {
	Cpu0  float32
	Cpu1  float32
	Cpu2  float32
	Cpu3  float32
	Cpu4  float32
	Cpu5  float32
	Cpu6  float32
	Cpu7  float32
	Cpu8  float32
	Cpu9  float32
	Cpu10 float32
	Cpu11 float32
	Cpu12 float32
	Cpu13 float32
	Cpu14 float32
	Cpu15 float32
}

//
// memory values are read from the /proc/meminfo pseudo-file with the help of the goprocinfo package...
// how are they calculated? Like 'htop' command, see question:
//   - http://stackoverflow.com/questions/41224738/how-to-calculate-memory-usage-from-proc-meminfo-like-htop/
//
type MyMemoInfo struct {
	TotalMachine       uint64
	TotalUsed          uint64
	Buffers            uint64
	Cached             uint64
	NonCacheNonBuffers uint64
}

var currCPUStats *linuxproc.Stat = ReadCPUStats()

var prevCPUStats *linuxproc.Stat = ReadCPUStats()
var coreStats MyCPUStats

var memoryInfo MyMemoInfo
var scanner = lmsensors.New()
var devices []*lmsensors.Device

var updateTime = time.Now().Add(5 * time.Second)

func ReadCPUStats() *linuxproc.Stat {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}
	// fmt.Println(stat)
	return stat
}

func calcMyCPUStats(curr, prev *linuxproc.Stat) *MyCPUStats {

	var stats MyCPUStats

	stats.Cpu0 = calcSingleCoreUsage(curr.CPUStats[0], prev.CPUStats[0])
	stats.Cpu1 = calcSingleCoreUsage(curr.CPUStats[1], prev.CPUStats[1])
	stats.Cpu2 = calcSingleCoreUsage(curr.CPUStats[2], prev.CPUStats[2])
	stats.Cpu3 = calcSingleCoreUsage(curr.CPUStats[3], prev.CPUStats[3])
	stats.Cpu4 = calcSingleCoreUsage(curr.CPUStats[4], prev.CPUStats[4])
	stats.Cpu5 = calcSingleCoreUsage(curr.CPUStats[5], prev.CPUStats[5])
	stats.Cpu6 = calcSingleCoreUsage(curr.CPUStats[6], prev.CPUStats[6])
	stats.Cpu7 = calcSingleCoreUsage(curr.CPUStats[7], prev.CPUStats[7])
	stats.Cpu8 = calcSingleCoreUsage(curr.CPUStats[8], prev.CPUStats[8])
	stats.Cpu9 = calcSingleCoreUsage(curr.CPUStats[9], prev.CPUStats[9])
	stats.Cpu10 = calcSingleCoreUsage(curr.CPUStats[10], prev.CPUStats[10])
	stats.Cpu11 = calcSingleCoreUsage(curr.CPUStats[11], prev.CPUStats[11])
	stats.Cpu12 = calcSingleCoreUsage(curr.CPUStats[12], prev.CPUStats[12])
	stats.Cpu13 = calcSingleCoreUsage(curr.CPUStats[13], prev.CPUStats[13])
	stats.Cpu14 = calcSingleCoreUsage(curr.CPUStats[14], prev.CPUStats[14])
	stats.Cpu15 = calcSingleCoreUsage(curr.CPUStats[15], prev.CPUStats[15])
	return &stats
}

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float32 {

	/*
	 *    source: http://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux
	 *
	 *    PrevIdle = previdle + previowait
	 *    Idle = idle + iowait
	 *
	 *    PrevNonIdle = prevuser + prevnice + prevsystem + previrq + prevsoftirq + prevsteal
	 *    NonIdle = user + nice + system + irq + softirq + steal
	 *
	 *    PrevTotal = PrevIdle + PrevNonIdle
	 *    Total = Idle + NonIdle
	 *
	 *    # differentiate: actual value minus the previous one
	 *    totald = Total - PrevTotal
	 *    idled = Idle - PrevIdle
	 *
	 *    CPU_Percentage = (totald - idled)/totald
	 */

	//
	//  Memory
	//
	//

	PrevIdle := prev.Idle + prev.IOWait
	Idle := curr.Idle + curr.IOWait

	PrevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	NonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal

	PrevTotal := PrevIdle + PrevNonIdle
	Total := Idle + NonIdle
	// fmt.Println(PrevIdle, Idle, PrevNonIdle, NonIdle, PrevTotal, Total)

	//  differentiate: actual value minus the previous one
	totald := Total - PrevTotal
	idled := Idle - PrevIdle

	CPU_Percentage := (float32(totald) - float32(idled)) / float32(totald)

	return CPU_Percentage
}

func ReadMemoInfo() *linuxproc.MemInfo {
	info, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("info read fail")
	}
	// fmt.Printf("Memory info struct:\n%+v", info)
	return info
}

func (w *StatsWidget) Update(dev *streamdeck.Device) {

	if updateTime.Before(time.Now()) {
		currCPUStats = ReadCPUStats()
		coreStats = *calcMyCPUStats(currCPUStats, prevCPUStats)
		prevCPUStats = currCPUStats
		info := ReadMemoInfo()
		memoryInfo.TotalMachine = info.MemTotal
		memoryInfo.TotalUsed = info.MemTotal - info.MemFree
		memoryInfo.Buffers = info.Buffers
		memoryInfo.Cached = info.Cached + info.SReclaimable - info.Shmem
		memoryInfo.NonCacheNonBuffers = memoryInfo.TotalUsed - (memoryInfo.Buffers + memoryInfo.Cached)
		var err error
		devices, err = scanner.Scan()
		/*
		for _, d := range devices {
			for _, s := range d.Sensors {
				switch v := s.(type) {
				case *lmsensors.FanSensor:
					//log.Println(d.Name, v.Name)
				case *lmsensors.TemperatureSensor:
					//log.Println(d.Name, v.Name)
				}
			}
		}
		*/
		if err != nil {
			log.Fatal(err)
		}

		updateTime = time.Now().Add(5 * time.Second)

	} else {
		//fmt.Println("Sensor update skip")
	}

}
