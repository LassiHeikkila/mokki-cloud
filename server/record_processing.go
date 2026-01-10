package server

import (
	"errors"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

func MeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil {
		return nil, errors.New("nil record")
	}

	switch r.Field() {
	case "pressure":
		return pressureMeasurementFromRecord(r)
	case "humidity":
		return humidityMeasurementFromRecord(r)
	case "temperature":
		return temperatureMeasurementFromRecord(r)
	case "batteryvoltage":
		return batteryVoltageMeasurementFromRecord(r)
	case "co2":
		return co2MeasurementFromRecord(r)
	case "pm2p5":
		return pm2p5MeasurementFromRecord(r)
	case "":
		return nil, errors.New("empty field")
	default:
		return nil, errors.New("unknown field: " + r.Field())
	}
}

func pressureMeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "pressure" {
		return nil, errors.New("not a pressure record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("pressure: missing sensormac field")
	}
	recordTime := r.Time()

	switch pressure := r.Value().(type) {
	case int64:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case int32:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case int16:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case int8:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case int:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case uint8:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case uint16:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case uint32:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case uint64:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case float32:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	case float64:
		return &PressureMeasurement{
			SensorID_: mac,
			Pressure_: int(pressure),
			Time_:     recordTime,
		}, nil
	}

	return nil, errors.New("pressure: cannot cast value int")

}

func humidityMeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "humidity" {
		return nil, errors.New("not a humidity record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("humidity: missing sensormac field")
	}
	humidity, ok := r.Value().(float64)
	if !ok {
		return nil, errors.New("humidity: cannot cast value to float64")
	}
	recordTime := r.Time()

	return &HumidityMeasurement{
		SensorID_: mac,
		Humidity_: humidity,
		Time_:     recordTime,
	}, nil
}

func temperatureMeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "temperature" {
		return nil, errors.New("not a temperature record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("temperature: missing sensormac field")
	}
	temperature, ok := r.Value().(float64)
	if !ok {
		return nil, errors.New("temperature: cannot cast value to float64")
	}
	recordTime := r.Time()

	return &TemperatureMeasurement{
		SensorID_:    mac,
		Temperature_: temperature,
		Time_:        recordTime,
	}, nil
}

func co2MeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "co2" {
		return nil, errors.New("not a co2 record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("co2: missing sensormac field")
	}
	recordTime := r.Time()

	switch co2 := r.Value().(type) {
	case int64:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case int32:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case int16:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case int8:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case int:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case uint8:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case uint16:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case uint32:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case uint64:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case float32:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	case float64:
		return &CO2Measurement{
			SensorID_: mac,
			CO2_: int(co2),
			Time_:     recordTime,
		}, nil
	}

	return nil, errors.New("pressure: cannot cast value int")
}

func pm2p5MeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "pm2p5" {
		return nil, errors.New("not a pm2p5 record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("pm2p5: missing sensormac field")
	}
	pm2p5, ok := r.Value().(float64)
	if !ok {
		return nil, errors.New("pm2p5: cannot cast value to float64")
	}
	recordTime := r.Time()

	return &PM2p5Measurement{
		SensorID_: mac,
		PM2p5_:    pm2p5,
		Time_:     recordTime,
	}, nil
}

func batteryVoltageMeasurementFromRecord(r *query.FluxRecord) (Measurement, error) {
	if r == nil || r.Field() != "batteryvoltage" {
		return nil, errors.New("not a batteryvoltage record")
	}

	mac, ok := r.ValueByKey("sensormac").(string)
	if !ok || mac == "" {
		return nil, errors.New("batteryvoltage: missing sensormac field")
	}
	voltage, ok := r.Value().(float64)
	if !ok {
		return nil, errors.New("batteryvoltage: cannot cast value to float64")
	}
	recordTime := r.Time()

	return &BatteryVoltageMeasurement{
		SensorID_: mac,
		Voltage_:  voltage,
		Time_:     recordTime,
	}, nil
}
