package bvfile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/proto"
)

const (
	BitsUint8  = 8
	BitsUint16 = 16
	BitsUint32 = 32
	BitsUint64 = 64
	BitsUintX  = 100 // string
)

// bv bin data format
type BvFileIndex struct {
	Name  string
	Size  int
	Index int
	Bits  int
}

func init() {
	// 检查 BvFileIndexTable 表 index是否正确
	var index int = 0

	for _, v := range BvFileIndexTable {
		if v.Index != index {
			panic("painc: BvFileIndexTable")
		}
		index = index + v.Size
	}
}

// BvFileIndexTable bv文件变量位置表
var BvFileIndexTable = []BvFileIndex{
	{Name: "CustomerSeriNum", Size: 8, Index: 0, Bits: BitsUint32},
	{Name: "ManufacSerial", Size: 8, Index: 8, Bits: BitsUint32},
	{Name: "CustomerResets", Size: 4, Index: 16, Bits: BitsUint16},
	{Name: "ConfigTime", Size: 8, Index: 20, Bits: BitsUint32},
	{Name: "LoggerStartDate", Size: 8, Index: 28, Bits: BitsUint32},
	{Name: "LoggerStopDate", Size: 8, Index: 36, Bits: BitsUint32},
	{Name: "StartupDelay", Size: 8, Index: 44, Bits: BitsUint32},
	{Name: "MeasurInterval", Size: 8, Index: 52, Bits: BitsUint32},
	{Name: "TemperatureUnits", Size: 1, Index: 60, Bits: BitsUint8},
	{Name: "HighBinTempTripPoint_A", Size: 4, Index: 61, Bits: BitsUint16},
	{Name: "HighBinFuncCounter_A", Size: 8, Index: 65, Bits: BitsUint32},
	{Name: "HighFuncLimit_A", Size: 8, Index: 73, Bits: BitsUint32},
	{Name: "HighBinTempTripPoint_B", Size: 4, Index: 81, Bits: BitsUint16},
	{Name: "HighBinFuncCounter_B", Size: 8, Index: 85, Bits: BitsUint32},
	{Name: "HighFuncLimit_B", Size: 8, Index: 93, Bits: BitsUint32},
	{Name: "LowBinTempTripPoint_A", Size: 4, Index: 101, Bits: BitsUint16},
	{Name: "LowBinFuncCounter_A", Size: 8, Index: 105, Bits: BitsUint32},
	{Name: "LowFuncLimit_A", Size: 8, Index: 113, Bits: BitsUint32},
	{Name: "LowBinTempTripPoint_B", Size: 4, Index: 121, Bits: BitsUint16},
	{Name: "LowBinFuncCounter_B", Size: 8, Index: 125, Bits: BitsUint32},
	{Name: "LowFuncLimit_B", Size: 8, Index: 133, Bits: BitsUint32},
	{Name: "AVR", Size: 4, Index: 141, Bits: BitsUint16},
	{Name: "MKT", Size: 4, Index: 145, Bits: BitsUint16},
	{Name: "TotalTime_HighB", Size: 8, Index: 149, Bits: BitsUint32},
	{Name: "TotalTime_HighA", Size: 8, Index: 157, Bits: BitsUint32},
	{Name: "TotalTime_2_8", Size: 8, Index: 165, Bits: BitsUint32},
	{Name: "TotalTime_LowA", Size: 8, Index: 173, Bits: BitsUint32},
	{Name: "TotalTime_LowB", Size: 8, Index: 181, Bits: BitsUint32},
	{Name: "NoOfViola_HighB", Size: 8, Index: 189, Bits: BitsUint32},
	{Name: "NoOfViola_HighA", Size: 8, Index: 197, Bits: BitsUint32},
	{Name: "NoOfViola_LowA", Size: 8, Index: 205, Bits: BitsUint32},
	{Name: "NoOfViola_LowB", Size: 8, Index: 213, Bits: BitsUint32},
	{Name: "LongestTime_HighB", Size: 8, Index: 221, Bits: BitsUint32},
	{Name: "LongestTime_HighA", Size: 8, Index: 229, Bits: BitsUint32},
	{Name: "LongestTime_LowA", Size: 8, Index: 237, Bits: BitsUint32},
	{Name: "LongestTime_LowB", Size: 8, Index: 245, Bits: BitsUint32},
	{Name: "LogerName", Size: 22, Index: 253, Bits: BitsUintX},
	{Name: "TimeBase", Size: 10, Index: 275, Bits: BitsUintX},
	{Name: "ConfigBy", Size: 20, Index: 285, Bits: BitsUintX},
	{Name: "Notes", Size: 96, Index: 305, Bits: BitsUintX},
	{Name: "NumberOfPoints", Size: 8, Index: 401, Bits: BitsUint32},
	//{Name: "PDFPassword", Size: 16, Index: 409, Bits: BitsUint64},
	//{Name: "Temperature_0", Size: 4, Index: 425, Bits: BitsUint16},
	{Name: "Temperature_0", Size: 4, Index: 409, Bits: BitsUint16},
}

func GetBvFileIndexByName(name string) *BvFileIndex {
	for _, v := range BvFileIndexTable {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

// BvFileInfo bv文件头信息
type BvFileInfo struct {
	CustomerSeriNum             string  // u32 顾客系列号
	ManufacSerial               string  // u32 厂家系列号
	CustomerResets              int64   // u16 顾客复位次数
	ConfigTime                  int64   // u32 配置时间
	LoggerStartDate             int64   // u32 启动时间
	LoggerStopDate              int64   // u32 停止时间
	StartupDelay                int64   // u32 启动延时时间
	MeasurInterval              int64   // u32 测量间隔
	TemperatureUnits            string  // u8 温度单位
	HighBinTempTripPoint_A      int64   // u16 上限配置值 ?
	HighBinTempTripPoint_A_Temp float32 // 上限配置值 - 温度
	HighBinTempTripPoint_A_Humi float32 // 上限配置值 - 湿度
	HighBinFuncCounter_A        int64   // u32 超上限秒数
	HighFuncLimit_A             string  // u32 上限值 ?
	HighBinTempTripPoint_B      int64   // u16 上上限配置值 ?
	HighBinTempTripPoint_B_Temp float32 // 上上限配置值 - 温度
	HighBinTempTripPoint_B_Humi float32 // 上上限配置值 - 湿度
	HighBinFuncCounter_B        int64   // u32 超上上限秒数
	HighFuncLimit_B             string  // u32 上上限值 ?
	LowBinTempTripPoint_A       int64   // u16 下限配置值 ?
	LowBinTempTripPoint_A_Temp  float32 // 下限配置值 - 温度
	LowBinTempTripPoint_A_Humi  float32 // 下限配置值 - 湿度
	LowBinFuncCounter_A         int64   // u32 超下限秒数
	LowFuncLimit_A              string  // u32 下限值 ?
	LowBinTempTripPoint_B       int64   // u16 下下限配置值 ?
	LowBinTempTripPoint_B_Temp  float32 // 下限配置值 - 温度
	LowBinTempTripPoint_B_Humi  float32 // 下限配置值 - 湿度
	LowBinFuncCounter_B         int64   // u32 超下下限秒数
	LowFuncLimit_B              string  // u32 下下限值 ?
	AVR                         float32 // u16 平均温度
	MKT                         float32 // u16  MKT 值
	TotalTime_HighB             int64   // u32 超上上限总的秒数
	TotalTime_HighA             int64   // u32 超上限总的秒数
	TotalTime_2_8               int64   // u32 未超上限和下限的秒数
	TotalTime_LowA              int64   // u32 超下限总的秒数
	TotalTime_LowB              int64   // u32 超下下限总的秒数
	NoOfViola_HighB             int64   // u32 超上上限总的点数
	NoOfViola_HighA             int64   // u32 超上限总的点数
	NoOfViola_LowA              int64   // u32 超下限总的点数
	NoOfViola_LowB              int64   // u32 超下下限总的点数
	LongestTime_HighB           int64   // u32 超上上限最长一次秒数
	LongestTime_HighA           int64   // u32 超上限最长一次秒数
	LongestTime_LowA            int64   // u32 超下限最长一次秒数
	LongestTime_LowB            int64   // u32 超下下限紧长一次秒数
	LogerName                   string  // u8 记录仪名称字符串
	TimeBase                    string  // u8 基于格林威治时间的各国时间描述符
	ConfigBy                    string  // u8 配置者描述符
	Notes                       string  // u8 留言/笔记字符串
	NumberOfPoints              int64   // u32 数据点数
	PDFPassword                 string  // u8 PDF文件的打开密码 ?

	Temperature_0 int64 // u16/u32  4字节时为温度或者 8字节时为温湿度
	// .
	// .
	// Temperature_N int64 // u16/u32

	Data []*SensorData // 传感器数据点
}

// Parse bv文件解析
func (bv *BvFileInfo) Parse(b []byte) error {
	if bv != nil {
		if len(b) < 53 {
			return errors.New("bv file format error")
		}

		bvV := reflect.ValueOf(bv)
		for _, v := range BvFileIndexTable {
			indexInfo := GetBvFileIndexByName(v.Name)
			if indexInfo.Bits == BitsUintX {
				stringV := DecodeToString(b[indexInfo.Index:indexInfo.Index+indexInfo.Size], indexInfo.Bits)
				if debugEnable {
					TraceInfo("%v of indexInfo: %v", v.Name, indexInfo)
					TraceInfo("% 02X\n", b[indexInfo.Index:indexInfo.Index+indexInfo.Size])
					TraceInfo("%v", stringV)
				}
				if bvV.Kind() != reflect.Ptr && !bvV.Elem().CanSet() {
					continue
				} else {
					f := bvV.Elem().FieldByName(v.Name)
					if !f.IsValid() {
						return fmt.Errorf("%s reflect invalid", v.Name)
					}
					switch f.Kind() {
					case reflect.String:
						f.SetString(stringV)
					}
				}
			} else {
				uint64V, err := DecodeToUInt(b[indexInfo.Index:indexInfo.Index+indexInfo.Size], indexInfo.Bits)
				if debugEnable {
					TraceInfo("%v of indexInfo: %v", v.Name, indexInfo)
					TraceInfo("% 02X", b[indexInfo.Index:indexInfo.Index+indexInfo.Size])
					TraceInfo("%v, %v", uint64V, err)
				}
				if err != nil {
					return err
				}
				if bvV.Kind() != reflect.Ptr && !bvV.Elem().CanSet() {
					continue
				} else {
					f := bvV.Elem().FieldByName(v.Name)
					if !f.IsValid() {
						return fmt.Errorf("%s reflect invalid", v.Name)
					}
					switch f.Kind() {
					case reflect.String:
						f.SetString(strconv.Itoa(int(uint64V)))
					case reflect.Int64:
						f.SetInt(int64(uint64V))
					case reflect.Float32:
						f.SetFloat(float64(uint64V))
					}
				}
			}

			// 特殊变量转换特殊处理
			switch v.Name {
			case "TemperatureUnits":
				if bv.TemperatureUnits == "12" {
					bv.TemperatureUnits = "C"
				}
				if bv.TemperatureUnits == "15" {
					bv.TemperatureUnits = "F"
				}
			case "PDFPassword":
				intV, err := strconv.Atoi(bv.PDFPassword)
				if err != nil {
					bv.PDFPassword = ""
				}
				bv.PDFPassword = fmt.Sprintf("%016X", uint64(intV))
			case "HighBinTempTripPoint_A":
				bv.HighBinTempTripPoint_A_Temp = float32(uint16(bv.HighBinTempTripPoint_A)&0xFF) / 10.0
				bv.HighBinTempTripPoint_A_Humi = float32((uint16(bv.HighBinTempTripPoint_A) >> 8) & 0xFF)
			case "HighBinTempTripPoint_B":
				bv.HighBinTempTripPoint_B_Temp = float32(uint16(bv.HighBinTempTripPoint_B)&0xFF) / 10.0
				bv.HighBinTempTripPoint_B_Humi = float32((uint16(bv.HighBinTempTripPoint_B) >> 8) & 0xFF)
			case "LowBinTempTripPoint_A":
				bv.LowBinTempTripPoint_A_Temp = float32(uint16(bv.LowBinTempTripPoint_A)&0xFF) / 10.0
				bv.LowBinTempTripPoint_A_Humi = float32((uint16(bv.LowBinTempTripPoint_A) >> 8) & 0xFF)
			case "LowBinTempTripPoint_B":
				bv.LowBinTempTripPoint_B_Temp = float32(uint16(bv.LowBinTempTripPoint_B)&0xFF) / 10.0
				bv.LowBinTempTripPoint_B_Humi = float32((uint16(bv.LowBinTempTripPoint_B) >> 8) & 0xFF)
			case "AVR":
				bv.AVR = bv.AVR / 10.0
			case "MKT":
				bv.MKT = bv.MKT / 10.0
			}
		}

		// parse t/h
		dataLen := (bv.LoggerStopDate-bv.LoggerStartDate)/bv.MeasurInterval + 1
		if dataLen != bv.NumberOfPoints {
			return fmt.Errorf("dataLen error: %v, %v", dataLen, bv.NumberOfPoints)
		}
		dataendIndex := bytes.Index(b, []byte("dataend"))
		if dataendIndex == -1 {
			return errors.New("dataend format error")
		}
		bits := (len(b[409:dataendIndex])) / int(dataLen)
		if (len(b[409:dataendIndex]))%int(dataLen) != 0 {
			return errors.New("data format error")
		}
		if debugEnable {
			TraceInfo("data area length: %v, %v, %v", len(b[409:dataendIndex]), int(dataLen*4), bits)
		}
		switch bits {
		case 4:
			// t only
			for i := 0; i < int(dataLen); i++ {
				sensorData := &SensorData{}
				sensorData.At = proto.Int64(bv.LoggerStartDate + int64(i*int(bv.MeasurInterval)))
				uint64V, err := DecodeToUInt(b[409+i*4:409+i*4+4], BitsUint16)
				if err != nil {
					return fmt.Errorf("[t]t: %v", err)
				}
				sensorData.SetT(float32(uint64V) / 10.0)
				bv.Data = append(bv.Data, sensorData)
			}
		case 8:
			// t and h
			for i := 0; i < int(dataLen); i++ {
				sensorData := &SensorData{}
				sensorData.At = proto.Int64(bv.LoggerStartDate + int64(i*int(bv.MeasurInterval)))
				uint64V, err := DecodeToUInt(b[409+i*4:409+i*4+4], BitsUint16)
				if err != nil {
					return fmt.Errorf("[th]t: %v", err)
				}
				sensorData.SetT(float32(uint64V) / 10.0)
				uint64V, err = DecodeToUInt(b[409+i*4+4:409+i*4+8], BitsUint16)
				if err != nil {
					return fmt.Errorf("[th]h: %v", err)
				}
				sensorData.SetH(float32(uint64V) / 10.0)
				bv.Data = append(bv.Data, sensorData)
			}
		default:
			return errors.New("unsupport bits length")
		}

		// parse end
		return nil
	}
	return errors.New("bv is nil")
}

// String 序列化成字符串
func (bv BvFileInfo) String() string {
	body, err := json.Marshal(bv)
	if err != nil {
		return ""
	}
	return string(body)
}
