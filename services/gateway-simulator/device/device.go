package device

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"gateway-simulator/telemetry"
)

type Device struct {
	ID         string
	Type       DeviceType
	GatewayID  string
	TenantID   string
	sensor     *SensorState
}

func NewDevice(deviceType DeviceType, index int, gatewayID, tenantID string) *Device {
	profile, ok := SensorProfiles[deviceType]
	if !ok {
		profile = SensorProfiles[TemperatureSensor]
	}

	return &Device{
		ID:        fmt.Sprintf("dev-%s-%03d", deviceType, index),
		Type:      deviceType,
		GatewayID: gatewayID,
		TenantID:  tenantID,
		sensor:    NewSensorState(profile),
	}
}

func (d *Device) GenerateTelemetry() telemetry.Payload {
	return telemetry.Payload{
		MessageID:    uuid.New().String(),
		GatewayID:    d.GatewayID,
		DeviceID:     d.ID,
		TenantID:     d.TenantID,
		DeviceType:   string(d.Type),
		Timestamp:    time.Now().UTC(),
		Measurements: d.sensor.GenerateReadings(),
	}
}

// GetDeviceTypes restituisce una lista dei tipi di device disponibili
func GetDeviceTypes() []DeviceType {
	return []DeviceType{
		TemperatureSensor,
		HumiditySensor,
		PowerMeter,
		AirQualitySensor,
		MotionSensor,
	}
}
