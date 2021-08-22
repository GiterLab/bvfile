package bvfile

import (
	"encoding/json"
	"errors"

	"google.golang.org/protobuf/proto"
)

// SensorData 传感器数据点
type SensorData struct {
	At *int64   `json:"at,omitempty"`
	T  *float32 `json:"t,omitempty"`
	H  *float32 `json:"h,omitempty"`
}

// GetAt 获取传感器采集时间点
func (d *SensorData) GetAt() int64 {
	if d != nil && d.At != nil {
		return *d.At
	}
	return 0
}

// GetT 获取传感器温度值
func (d *SensorData) GetT() float32 {
	if d != nil && d.T != nil {
		return *d.T
	}
	return -10000
}

// SetT 设置传感器温度值
func (d *SensorData) SetT(t float32) error {
	if d == nil {
		return errors.New("d is nil")
	}
	d.T = proto.Float32(t)
	return nil
}

// GetH 获取传感器湿度值
func (d *SensorData) GetH() float32 {
	if d != nil && d.H != nil {
		return *d.H
	}
	return -10000
}

// SetH 设置传感器湿度值
func (d *SensorData) SetH(h float32) error {
	if d == nil {
		return errors.New("d is nil")
	}
	d.H = proto.Float32(h)
	return nil
}

// String 序列化成字符串
func (d SensorData) String() string {
	body, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(body)
}

// SensorDataList 传感器数据点
type SensorDataList []*SensorData

func (a SensorDataList) Len() int           { return len(a) }
func (a SensorDataList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SensorDataList) Less(i, j int) bool { return a[i].GetT() < a[j].GetT() }
