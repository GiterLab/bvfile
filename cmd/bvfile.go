package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/GiterLab/bvfile"
)

func ReadBVFile(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return buf
}

func main() {
	var bvBytes []byte

	appVersion := flag.Bool("v", false, "Version of app")
	debug := flag.Bool("d", false, "debug mode")
	bvFileName := flag.String("f", "", "bv filename")
	flag.Parse()

	if *appVersion {
		fmt.Println("Version:", bvfile.Version)
	}

	if *bvFileName != "" {
		bvBytes = ReadBVFile(*bvFileName)
	} else {
		fmt.Println("please specify the bv filename first")
		os.Exit(0)
	}

	if *debug {
		// debug
		fmt.Printf("%v\n", hex.Dump(bvBytes))
		bvfile.Debug(true)
		bvfile.SetUserDebug(func(format string, level int, v ...interface{}) {
			switch level {
			case bvfile.LevelInformational:
				fmt.Println(fmt.Sprintf(format, v...))
			case bvfile.LevelError:
				fmt.Println(fmt.Sprintf(format, v...))
			}
		})
	}

	bvInfo := &bvfile.BvFileInfo{}
	err := bvInfo.Parse(bvBytes)
	if err != nil {
		fmt.Println("parse bv file failed:", err)
	}
	if *debug {
		fmt.Println(bvInfo)
	}

	fmt.Println("CustomerSeriNum:", bvInfo.CustomerSeriNum)
	fmt.Println("ManufacSerial:", bvInfo.ManufacSerial)
	fmt.Println("LoggerStartDate:", bvInfo.LoggerStartDate, time.Unix(bvInfo.LoggerStartDate, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("LoggerStopDate:", bvInfo.LoggerStopDate, time.Unix(bvInfo.LoggerStopDate, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("StartupDelay:", bvInfo.StartupDelay)
	fmt.Println("TemperatureUnits:", bvInfo.TemperatureUnits)
	fmt.Println("AVR:", bvInfo.AVR)
	fmt.Println("MKT:", bvInfo.MKT)
	fmt.Println("HighBinTempTripPoint_B_Temp & Humi:", bvInfo.HighBinTempTripPoint_B_Temp, bvInfo.HighBinTempTripPoint_B_Humi)
	fmt.Println("HighBinTempTripPoint_A_Temp & Humi:", bvInfo.HighBinTempTripPoint_A_Temp, bvInfo.HighBinTempTripPoint_A_Humi)
	fmt.Println("LowBinTempTripPoint_A_Temp & Humi:", bvInfo.LowBinTempTripPoint_A_Temp, bvInfo.LowBinTempTripPoint_A_Humi)
	fmt.Println("LowBinTempTripPoint_B_Temp & Humi:", bvInfo.LowBinTempTripPoint_B_Temp, bvInfo.LowBinTempTripPoint_B_Humi)
	fmt.Println("LogerName:", bvInfo.LogerName)
	fmt.Println("TimeBase:", bvInfo.TimeBase)
	fmt.Println("ConfigBy:", bvInfo.ConfigBy)
	fmt.Println("NumberOfPoints:", bvInfo.NumberOfPoints)

	// 按照温度从小到大排序
	sort.Sort(bvfile.SensorDataList(bvInfo.Data))
	if len(bvInfo.Data) != 0 {
		fmt.Println("MinT:", bvInfo.Data[0])
		fmt.Println("MaxT:", bvInfo.Data[len(bvInfo.Data)-1])
	}
}
