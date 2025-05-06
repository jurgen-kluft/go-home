package mqtt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/thomasf/yeelight/pkg/yeel"
)

// model has handlermaps of a specific light models.
type model struct {
	printers   map[string]printer
	commanders map[string]commander
}

type printer func(*Device) string
type commander func(*Device, *Command) (yeel.Commander, error)

// all model handlers are stored here.
var models map[string]model

// Command to be sent to Yeebulb's.
type Command struct {
	Command     string
	DeviceName  string
	DeviceModel string
	Value       string
}

// PropUpdate field level property updates from yeebulb tcp connections.
type PropUpdate struct {
	DeviceID string
	Prop     string
	Value    []byte
}

func (p PropUpdate) String() string {
	return string(p.Value)
}

func init() {

	// set up the models map.

	assignedPrinters := make(map[string]bool, 0)
	assignedCommands := make(map[string]bool, 0)

	addCommander := func(name string, m *model) {
		if _, ok := m.commanders[name]; ok {
			panic(fmt.Errorf("commander '%s' already added to '%v'", name, m))
		}
		if cmd, ok := allCommands[name]; ok {
			assignedCommands[name] = true
			m.commanders[name] = cmd
			return
		}
		err := fmt.Errorf("command function '%s' not found", name)
		panic(err)
	}

	addPrinter := func(name string, m *model) {
		if _, ok := m.printers[name]; ok {
			panic(fmt.Errorf("printer '%s' already added to '%v'", name, m))
		}
		if printer, ok := allPrinters[name]; ok {
			assignedPrinters[name] = true
			m.printers[name] = printer
			return
		}
		err := fmt.Errorf("printer function '%s' not found", name)
		panic(err)
	}

	colorModel := model{
		printers:   make(map[string]printer, 0),
		commanders: make(map[string]commander, 0),
	}
	stripModel := model{
		printers:   make(map[string]printer, 0),
		commanders: make(map[string]commander, 0),
	}
	monoModel := model{
		printers:   make(map[string]printer, 0),
		commanders: make(map[string]commander, 0),
	}

	// add common commands and printers
	for _, name := range []string{
		"brightness", "name", "power", "start_animation", "start_flow", "stop_flow",
	} {
		addCommander(name, &colorModel)
		addCommander(name, &stripModel)
		addCommander(name, &monoModel)
	}
	for _, name := range []string{
		"brightness", "flow_params", "id", "is_flowing", "location", "model", "power", "transition",
	} {
		addPrinter(name, &colorModel)
		addPrinter(name, &stripModel)
		addPrinter(name, &monoModel)
	}

	// add color specific commands and printers
	for _, name := range []string{
		"color_temp", "rgb",
	} {
		addCommander(name, &colorModel)
		addCommander(name, &stripModel)
	}
	for _, name := range []string{
		"color_mode", "color_temp", "hue", "rgb", "saturation",
	} {
		addPrinter(name, &colorModel)
		addPrinter(name, &stripModel)
	}

	var unusedCommands []string
	for k, _ := range allCommands {
		if _, ok := assignedCommands[k]; !ok {
			unusedCommands = append(unusedCommands, k)
		}
	}
	var unusedPrinters []string
	for k, _ := range allPrinters {
		if _, ok := assignedPrinters[k]; !ok {
			unusedPrinters = append(unusedPrinters, k)
		}
	}

	if len(unusedCommands) > 0 || len(unusedPrinters) > 0 {
		panic(fmt.Errorf("unused commands: %v, unused printers: %v", unusedCommands, unusedPrinters))
	}

	models = map[string]model{
		"color": colorModel,
		"strip": stripModel,
		"mono":  monoModel,
	}
}

// allPrinters contains all different property printing commands supported by
// the yeelight protocol.
var allPrinters = map[string](func(*Device) string){

	"brightness": func(d *Device) string {
		return fmt.Sprintf("%d", d.Brightness)
	},

	"rgb": func(d *Device) string {
		if d.ColorMode != 1 || d.Flowing == 1 {
			return ""
		}
		return yeel.IntToRGB(d.RGB).String()
	},

	"power": func(d *Device) string {
		return d.Power
	},

	"location": func(d *Device) string {
		return d.Location
	},

	"id": func(d *Device) string {
		return d.ID
	},

	"model": func(d *Device) string {
		return d.Model
	},

	"color_mode": func(d *Device) string {
		if d.Flowing == 1 {
			return "FLOW"
		}
		switch d.ColorMode {
		case 1:
			return "RGB"
		case 2:
			return "TEMP"
		case 3:
			return "HSV"
		default:
			return "???"
		}
	},

	"color_temp": func(d *Device) string {
		if d.ColorMode != 2 {
			return ""
		}
		return fmt.Sprintf("%d", d.ColorTemprature)

	},

	"hue": func(d *Device) string {
		if d.ColorMode != 3 || d.Flowing == 1 {
			return ""
		}
		return fmt.Sprintf("%d", d.Hue)
	},

	"saturation": func(d *Device) string {
		if d.ColorMode != 3 || d.Flowing == 1 {
			return ""
		}
		return fmt.Sprintf("%d", d.Saturation)
	},

	"transition": func(d *Device) string {
		return fmt.Sprintf("%d", d.Transition/time.Second)
	},

	"flow_params": func(d *Device) string {
		if d.Flowing == 1 {
			return d.FlowParams
		}
		return ""
	},

	"is_flowing": func(d *Device) string {
		if d.Flowing == 1 {
			return "on"
		}
		return "off"
	},
}

// allCommands contains all the printers needed to construct the various
// models.
var allCommands = map[string](func(*Device, *Command) (yeel.Commander, error)){

	"color_temp": func(d *Device, msg *Command) (yeel.Commander, error) {
		v, err := strconv.ParseInt(msg.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		cmd := yeel.SetColorTempratureCommand{
			ColorTemprature: int(v),
			Duration:        d.Transition,
		}
		return cmd, nil
	},

	"brightness": func(d *Device, msg *Command) (yeel.Commander, error) {
		bri, err := strconv.ParseUint(msg.Value, 10, 8)
		if err != nil {
			return nil, err
		}

		cmd := yeel.SetBrightnessCommand{
			Brightness: int(bri),
			Duration:   d.Transition,
		}
		return cmd, nil
	},
	"power": func(d *Device, msg *Command) (yeel.Commander, error) {
		var power bool
		switch msg.Value {
		case "on":
			power = true
		case "off":
			power = false
		case "toggle":
			power = d.Power != "on"
		default:
			return yeel.GetPropCommand{
				Properties: []string{"power"},
			}, nil
		}

		if (power && d.Power == "on") || (!power && d.Power == "off") {
			return nil, nil
		}
		return yeel.SetPowerCommand{
			On:       power,
			Duration: d.Transition,
		}, nil
	},

	"rgb": func(d *Device, msg *Command) (yeel.Commander, error) {
		if d.Model != "color" {
			return nil, fmt.Errorf("not an rgb bulb: %s", msg.Command)
		}
		strvals := strings.Split(msg.Value, ",")
		if len(strvals) != 3 {
			return nil, fmt.Errorf("expected three commma separated values: %s", msg.Value)
		}
		var intvals []uint8
		for _, v := range strvals {
			i, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return nil, err
			}
			intvals = append(intvals, uint8(i))
		}

		command := yeel.SetRGBCommand{
			Duration: d.Transition,
			R:        int(intvals[0]),
			G:        int(intvals[1]),
			B:        int(intvals[2]),
		}
		return command, nil
	},

	"start_animation": func(d *Device, msg *Command) (yeel.Commander, error) {
		if v, ok := animations[msg.Value]; ok {
			return yeel.StartColorFlowCommand{
				Count:     0,
				Action:    1,
				Animation: v,
			}, nil
		}
		return nil, errors.New("some-error")
	},

	"stop_flow": func(d *Device, msg *Command) (yeel.Commander, error) {
		return yeel.StopColorFlowCommand{}, nil

	},

	"start_flow": func(d *Device, msg *Command) (yeel.Commander, error) {
		return yeel.StartColorFlowCommand{
			Count:  10,
			Action: 1,
			Animation: yeel.Animation{
				yeel.RawKeyframe{RawExpression: msg.Value},
			},
		}, nil

	},

	"name": func(d *Device, msg *Command) (yeel.Commander, error) {
		if msg.Value != "" {
			return yeel.SetNameCommand{Name: msg.Value}, nil
		}
		return nil, errors.New("no name")
	}}
