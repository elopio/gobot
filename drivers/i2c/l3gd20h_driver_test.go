package i2c

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*HMC6352Driver)(nil)

// --------- HELPERS
func initTestL3GD20HDriver() (driver *L3GD20HDriver) {
	driver, _ = initTestL3GD20HDriverWithStubbedAdaptor()
	return
}

func initTestL3GD20HDriverWithStubbedAdaptor() (*L3GD20HDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewL3GD20HDriver(adaptor), adaptor
}

// --------- TESTS

func TestNewL3GD20HDriver(t *testing.T) {
	// Does it return a pointer to an instance of HMC6352Driver?
	var d interface{} = NewL3GD20HDriver(newI2cTestAdaptor())
	_, ok := d.(*L3GD20HDriver)
	if !ok {
		t.Errorf("NewL3GD20HDriver() should have returned a *L3GD20HDriver")
	}
}

func TestL3GD20HDriver(t *testing.T) {
	d := initTestL3GD20HDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

// Methods
func TestL3GD20HDriverStart(t *testing.T) {
	d, _ := initTestL3GD20HDriverWithStubbedAdaptor()

	gobottest.Assert(t, d.Start(), nil)
}

func TestL3GD20HDriverHalt(t *testing.T) {
	d := initTestL3GD20HDriver()

	gobottest.Assert(t, d.Halt(), nil)
}

func TestL3GD20HDriverScale(t *testing.T) {
	d := initTestL3GD20HDriver()
	gobottest.Assert(t, d.Scale(), L3GD20HScale250dps)

	d.SetScale(L3GD20HScale500dps)
	gobottest.Assert(t, d.Scale(), L3GD20HScale500dps)
}

func TestL3GD20HDriverMeasurement(t *testing.T) {
	d, adaptor := initTestL3GD20HDriverWithStubbedAdaptor()
	rawX := 5
	rawY := 8
	rawZ := -3
	adaptor.i2cReadImpl = func() ([]byte, error) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, int16(rawX))
		binary.Write(buf, binary.LittleEndian, int16(rawY))
		binary.Write(buf, binary.LittleEndian, int16(rawZ))
		return buf.Bytes(), nil
	}

	d.Start()
	x, y, z, err := d.XYZ()
	gobottest.Assert(t, err, nil)
	var sensitivity float32 = 0.00875
	gobottest.Assert(t, x, float32(rawX)*sensitivity)
	gobottest.Assert(t, y, float32(rawY)*sensitivity)
	gobottest.Assert(t, z, float32(rawZ)*sensitivity)
}
