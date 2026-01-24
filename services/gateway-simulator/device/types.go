package device

type DeviceType string

const (
	TemperatureSensor DeviceType = "temperature_sensor"
	HumiditySensor    DeviceType = "humidity_sensor"
	PowerMeter        DeviceType = "power_meter"
	AirQualitySensor  DeviceType = "air_quality_sensor"
	MotionSensor      DeviceType = "motion_sensor"
)

type MeasurementRange struct {
	Min   float32
	Max   float32
	Drift float32
}

type SensorProfile struct {
	Type         DeviceType
	Measurements map[string]MeasurementRange
}

var SensorProfiles = map[DeviceType]SensorProfile{
	TemperatureSensor: {
		Type: TemperatureSensor,
		Measurements: map[string]MeasurementRange{
			"temperature": {Min: 15, Max: 35, Drift: 0.5},
			"humidity":    {Min: 30, Max: 80, Drift: 2.0},
		},
	},
	HumiditySensor: {
		Type: HumiditySensor,
		Measurements: map[string]MeasurementRange{
			"humidity":  {Min: 20, Max: 90, Drift: 2.0},
			"dew_point": {Min: 5, Max: 25, Drift: 1.0},
		},
	},
	PowerMeter: {
		Type: PowerMeter,
		Measurements: map[string]MeasurementRange{
			"power":   {Min: 100, Max: 5000, Drift: 50},
			"voltage": {Min: 220, Max: 240, Drift: 2},
			"current": {Min: 0.5, Max: 25, Drift: 0.5},
		},
	},
	AirQualitySensor: {
		Type: AirQualitySensor,
		Measurements: map[string]MeasurementRange{
			"co2":  {Min: 400, Max: 2000, Drift: 20},
			"pm25": {Min: 0, Max: 150, Drift: 5},
			"voc":  {Min: 0, Max: 500, Drift: 10},
		},
	},
	MotionSensor: {
		Type: MotionSensor,
		Measurements: map[string]MeasurementRange{
			"motion_detected": {Min: 0, Max: 1, Drift: 1},
			"occupancy_count": {Min: 0, Max: 10, Drift: 2},
		},
	},
}
