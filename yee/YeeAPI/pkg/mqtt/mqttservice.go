package mqtt

import (
	"log"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/thomasf/lg"
	"github.com/thomasf/yeelight/pkg/ssdp"
	"github.com/thomasf/yeelight/pkg/yeel"
)

// MQTTService .
type MQTTService struct {
	MQTTConnectionStr string
	ClientID          string
	MulticastIF       string // which network interface to use for multicast udp

}

func (m *MQTTService) Start() error {

	{
		errorLog := log.New(os.Stdout, "", 0)
		MQTT.ERROR = errorLog
		lg.CopyLoggerTo("error", errorLog)

		warnLog := log.New(os.Stdout, "", 0)
		MQTT.WARN = warnLog
		lg.CopyLoggerTo("warning", warnLog)

		fatalLog := log.New(os.Stdout, "", 0)
		MQTT.CRITICAL = fatalLog
		lg.CopyLoggerTo("fatal", fatalLog)

	}

	commandCh := make(chan Command, 100)
	var mqconf mqttconfig
	{
		var err error
		mqconf, err = parseConnStr(m.MQTTConnectionStr)
		if err != nil {
			return err
		}
	}
	lg.Infoln(mqconf)

	clientID := m.ClientID
	if clientID == "" {
		clientID = "yeelight-proxy"
	}
	mqttOpts := MQTT.NewClientOptions().
		AddBroker(mqconf.Broker).
		SetClientID(clientID)

	if mqconf.User != "" {
		mqttOpts.SetUsername(mqconf.User)
	}

	if mqconf.Password != "" {
		mqttOpts.SetPassword(mqconf.Password)
	}

	mqttClient := MQTT.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	handleMessage := func(client MQTT.Client, msg MQTT.Message) {
		lg.V(10).Infof("message: %s - %s\n", msg.Topic, msg.Payload)
		cmdMessage, err := parseMessage(msg)
		if err != nil {
			lg.Errorln(err)
			return
		}
		commandCh <- cmdMessage

	}
	for _, v := range []string{"ylcolor/command/#", "ylmono/command/#", "ylstrip/command/#"} {
		if token := mqttClient.Subscribe(v, 0, handleMessage); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	deviceCh := make(chan yeel.Device, 100)
	conn := ssdp.Conn{}
	err := conn.Start(m.MulticastIF)
	if err != nil {
		lg.Fatal(err)
	}

	go func() {
		err := conn.Listen(deviceCh)
		if err != nil {
			lg.Errorln(err)
		}
	}()
	go func() {
		for {
			err := conn.Search()
			if err != nil {
				lg.Errorln(err)
			}
			time.Sleep(time.Hour)
		}
	}()

	conns := make(map[string]*yeel.Conn, 0)
	devices := make(map[string]*Device, 0)
	devicesRetries := make(map[string]int, 0)
	connEndedC := make(chan string, 0)
	notificationC := make(chan yeel.Notification, 0)

	checkPower := time.NewTicker(time.Minute)
	checkDelayCommands := time.NewTicker(100 * time.Millisecond)

	type delayedCommand struct {
		Created time.Time
		Command Command
	}
	var delayedCommands []delayedCommand

	handleConn := func(conn *yeel.Conn) {
		connectedBulbs.Inc()
		defer connectedBulbs.Dec()
		err := conn.Open()
		if err != nil {
			lg.Warningln(err)
		}
		connEndedC <- conn.Device.ID
	}

	findDevice := func(cmdMsg Command) *Device {
		var td *Device
		td, ok := devices[cmdMsg.DeviceName]

		if !ok {
		find:
			for _, v := range devices {
				if cmdMsg.DeviceName == v.Device.Name && cmdMsg.DeviceModel == v.Device.Model {
					td = v
					ok = true
					break find
				}
			}
		}
		if !ok {
			lg.Errorf("device not found: %v", cmdMsg)
		}
		return td
	}

	sendCommand := func(td *Device, cmdMsg Command) {
		command, err := td.Command(cmdMsg)
		if err != nil {
			lg.Warningln(err)
		}

		if command != nil {
			if conn, ok := conns[td.Device.ID]; ok {
				go func(conn *yeel.Conn, commander yeel.Commander) {
					res, err := conn.ExecCommand(command)
					if err != nil {
						lg.Warningln(err)
					} else {
						lg.V(10).Infoln(res)
					}
				}(conn, command)
			}
		}
	}

	requiresPowerOn := map[string]bool{
		"rgb":        true,
		"color_temp": true,
		"brightness": true,
	}

loop:
	for {
		select {
		case <-checkPower.C:
			for _, conn := range conns {
				go func(c *yeel.Conn) {
					res, err := c.ExecCommand(
						yeel.GetPropCommand{
							Properties: []string{"power"},
						},
					)
					if err != nil {
						lg.Errorln(err)
					}
					d := devices[c.Device.ID]
					d.Device.Power = res.Result[0]
					publishUpdates(d, mqttClient)

					lg.V(10).Infoln(c.Device.ID, res)
				}(conn)
			}

		case <-checkDelayCommands.C:
			if len(delayedCommands) < 1 {
				continue loop
			}
			var newDelayCommands []delayedCommand
			now := time.Now()
		cmds:
			for _, v := range delayedCommands {
				if now.After(v.Created.Add(2 * time.Second)) {
					continue cmds
				}
				td := findDevice(v.Command)
				if td != nil && td.Power == "on" {
					sendCommand(td, v.Command)
					continue cmds
				}
				newDelayCommands = append(newDelayCommands, v)
			}
			delayedCommands = newDelayCommands

		case cmdMsg := <-commandCh:
			td := findDevice(cmdMsg)
			if td == nil {
				continue loop
			}
			if td.Power == "off" && requiresPowerOn[cmdMsg.Command] {
				delayedCommands = append(delayedCommands, delayedCommand{
					Created: time.Now(),
					Command: cmdMsg,
				})
				continue loop
			}
			sendCommand(td, cmdMsg)

		case id := <-connEndedC:
			delete(conns, id)
			if d, ok := conns[id]; ok {
				retries := devicesRetries[id]
				if retries < 50 {
					conn := &yeel.Conn{
						Device:        d.Device,
						NotificationC: notificationC,
					}
					conns[d.Device.ID] = conn
					devicesRetries[id] = retries + 1
					go handleConn(conn)
				} else {
					delete(devices, id)
					delete(devicesRetries, id)
				}
			}

		case n := <-notificationC:
			if td, ok := devices[n.DeviceID]; ok {
				td.Update(n)
				publishUpdates(td, mqttClient)
			} else {
				lg.Errorln(n)
			}

		case yd := <-deviceCh:
			if td, ok := devices[yd.ID]; ok {
				td.LastUpdated = time.Now()
				td.Device = yd
				if publishUpdates(td, mqttClient) {
					lg.V(8).Infoln("updated device:", td.Device.ID)
				}
			} else {
				td = &Device{
					Device:      yd,
					Transition:  time.Second,
					LastUpdated: time.Now(),
				}
				devices[yd.ID] = td
				publishUpdates(td, mqttClient)
			}
			{
				prevConn, ok := conns[yd.ID]
				if ok && prevConn.Device.Location != yd.Location {
					err := prevConn.Close()
					if err != nil {
						lg.Warningln(err)
					}
				}
				if !ok || prevConn.Device.Location != yd.Location {
					conn := &yeel.Conn{
						Device:        yd,
						NotificationC: notificationC,
					}
					conns[yd.ID] = conn
					go handleConn(conn)
				} else {
					lg.Infoln("already connected to", yd.ID, yd.Location)
				}
			}
		}
	}
}
