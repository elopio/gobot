package dragonboard410c

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gobot.io/x/gobot/sysfs"
)

// Adaptor is the gobot.Adaptor representation for the Beaglebone
type Adaptor struct {
	name   string
	kernel string
	usrLed string
}

// NewAdaptor returns a new Beaglebone Adaptor
func NewAdaptor() *Adaptor {
	b := &Adaptor{
		name: "DragonBoard410c",
	}

	b.setSlots()
	return b
}

func (b *Adaptor) setSlots() {
	b.kernel = getKernel()
	b.usrLed = "/sys/class/leds/apq8016-sbc:green:"
}

// Name returns the Adaptor name
func (b *Adaptor) Name() string { return b.name }

// SetName sets the Adaptor name
func (b *Adaptor) SetName(n string) { b.name = n }

// Kernel returns the Linux kernel version for the BeagleBone
func (b *Adaptor) Kernel() string { return b.kernel }

// Connect initializes the pwm and analog dts.
func (b *Adaptor) Connect() error {
	return nil
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (b *Adaptor) Finalize() error {
	return nil
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (b *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	if strings.Contains(pin, "user") {
		fi, err := sysfs.OpenFile(b.usrLed+pin+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil {
			return err
		}
		_, err = fi.WriteString(strconv.Itoa(int(val)))
		return err
	}
	return nil
}

func getKernel() string {
	result, _ := exec.Command("uname", "-r").Output()

	return strings.TrimSpace(string(result))
}
