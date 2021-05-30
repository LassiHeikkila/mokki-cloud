package server

import (
	"time"
)

type Measurement interface {
	SensorID() string
	Measurement() string
	Value() interface{}
	Time() time.Time
}

type BatteryVoltageMeasurement struct {
	SensorID_ string `json:"sensorID"`

	// Voltage given in volts
	Voltage_ float64 `json:"voltage"`

	// Time when measurement was recorded
	Time_ time.Time `json:"time"`
}

func (m *BatteryVoltageMeasurement) Measurement() string {
	return "batteryvoltage"
}

func (m *BatteryVoltageMeasurement) SensorID() string {
	return m.SensorID_
}

func (m *BatteryVoltageMeasurement) Value() interface{} {
	return m.Voltage_
}

func (m *BatteryVoltageMeasurement) Time() time.Time {
	return m.Time_
}

type TemperatureMeasurement struct {
	SensorID_ string `json:"sensorID"`

	// Temperature given in Celsius
	Temperature_ float64 `json:"temperature"`

	// Time when measurement was recorded
	Time_ time.Time `json:"time"`
}

func (m *TemperatureMeasurement) Measurement() string {
	return "temperature"
}

func (m *TemperatureMeasurement) SensorID() string {
	return m.SensorID_
}

func (m *TemperatureMeasurement) Value() interface{} {
	return m.Temperature_
}

func (m *TemperatureMeasurement) Time() time.Time {
	return m.Time_
}

type HumidityMeasurement struct {
	SensorID_ string `json:"sensorID"`

	// Humidity given in percent
	Humidity_ float64 `json:"humidity"`

	// Time when measurement was recorded
	Time_ time.Time `json:"time"`
}

func (m *HumidityMeasurement) Measurement() string {
	return "humidity"
}

func (m *HumidityMeasurement) SensorID() string {
	return m.SensorID_
}

func (m *HumidityMeasurement) Value() interface{} {
	return m.Humidity_
}

func (m *HumidityMeasurement) Time() time.Time {
	return m.Time_
}

type PressureMeasurement struct {
	SensorID_ string `json:"sensorID"`

	// Pressure given in Pascal
	Pressure_ int `json:"pressure"`

	// Time when measurement was recorded
	Time_ time.Time `json:"time"`
}

func (m *PressureMeasurement) Measurement() string {
	return "pressure"
}

func (m *PressureMeasurement) SensorID() string {
	return m.SensorID_
}

func (m *PressureMeasurement) Value() interface{} {
	return m.Pressure_
}

func (m *PressureMeasurement) Time() time.Time {
	return m.Time_
}
