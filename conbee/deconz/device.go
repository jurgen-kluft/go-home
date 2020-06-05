package deconz

// Devices is a map of devices indexed by their id
type Devices map[string]Device // UniqueID -> Device

// Device is a device like a Light, Sensor(Motion, Contact, Switch)
type Device struct {
	Type     string
	Name     string
	DeviceID string
}
