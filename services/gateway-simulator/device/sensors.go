package device

import "math/rand"

type SensorState struct {
	currentValues map[string]float32
	profile       SensorProfile
}

func NewSensorState(profile SensorProfile) *SensorState {
	state := &SensorState{
		currentValues: make(map[string]float32),
		profile:       profile,
	}

	for name, r := range profile.Measurements {
		state.currentValues[name] = r.Min + rand.Float32()*(r.Max-r.Min)
	}

	return state
}

func (s *SensorState) GenerateReadings() map[string]interface{} {
	readings := make(map[string]interface{})

	for name, r := range s.profile.Measurements {
		current := s.currentValues[name]

		drift := (rand.Float32()*2 - 1) * r.Drift
		newValue := current + drift

		if newValue < r.Min {
			newValue = r.Min
		}
		if newValue > r.Max {
			newValue = r.Max
		}

		s.currentValues[name] = newValue

		if name == "motion_detected" {
			readings[name] = newValue > 0.5
		} else {
			readings[name] = float32(int(newValue*100)) / 100
		}
	}

	return readings
}
