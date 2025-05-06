package yeel

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/thomasf/lg"
)

// Conn .
type Conn struct {
	Device Device

	// NotificationC chan string
	commandC      chan Command
	NotificationC chan Notification
}

func (c *Conn) Close() error {
	lg.Warningln("close not implemented")
	return nil
}

// SendCommand sends a single named command to the Yeelight hub via TCP
func (c *Conn) Open() error {

	c.commandC = make(chan Command, 0)
	responseC := make(chan string, 10)
	resultC := make(chan Result, 10)
	notificationC := make(chan Notification, 10)
	exitC := make(chan bool, 0)
	var err error

	conn, err := net.DialTimeout("tcp", c.Device.LocationAddr(), time.Second*5)

	if err != nil {
		lg.Errorln("Failed to connect: %s\n", err)
		return err
	}
	lg.Infoln("opened conn to ", c.Device.Name, c.Device.Location)

	go func() {
		// TrackedCommand .
		type TrackedCommand struct {
			Command Command
			Created time.Time
		}
		commands := make(map[int]TrackedCommand, 0)
		cleanupTicker := time.NewTimer(30 * time.Second)
		for {
			select {
			case <-exitC:
				lg.Infoln("exiting...")
				for _, v := range commands {
					v.Command.resultC <- Result{
						Err: ErrClosing,
					}
					close(v.Command.resultC)
				}
				return

			case command := <-c.commandC:
				commands[command.ID] = TrackedCommand{command, time.Now()}
				data, err := json.Marshal(&command)
				if err != nil {
					lg.Fatal(err)
				}

				lg.V(10).Infoln(string(data))

				if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
					lg.Errorln(err)
				}

				_, err = fmt.Fprintf(conn, "%s\r\n", data)
				if err != nil {
					lg.Fatal(err)
				}

			case r := <-resultC:
				lg.V(10).Infoln(r)
				if v, ok := commands[r.ID]; ok {
					go func(c Command, r Result) {
						c.resultC <- r
						close(c.resultC)
					}(v.Command, r)
					delete(commands, r.ID)
				} else {
					lg.Warningln("response from untracked command: %v", r)
				}

			case n := <-notificationC:
				lg.V(10).Infoln(n)
				if c.NotificationC != nil {
					c.NotificationC <- n
				}

			case <-cleanupTicker.C:
				minuteAgo := time.Now().Add(-time.Minute)
				for k, v := range commands {
					if v.Created.Before(minuteAgo) {
						v.Command.resultC <- Result{
							Err: ErrTimeout,
						}
						close(v.Command.resultC)
						delete(commands, k)
					}
				}
			}
		}
	}()

	go func() {
	loop:
		for str := range responseC {
			var r resultOrNotification
			err := json.Unmarshal([]byte(str), &r)
			if err != nil {
				lg.Errorln(err, str)
			}
			if r.Notification == nil && r.Result == nil {
				log.Printf("could not parse message from: %s", str)
				continue loop
			}
			if r.Notification != nil {
				n := *r.Notification
				n.DeviceID = c.Device.ID
				notificationC <- n
			}
			if r.Result != nil {
				resultC <- *r.Result
			}
		}
	}()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		response := scanner.Text()
		lg.V(10).Infoln("response", response)
		responseC <- response
	}
	close(responseC)

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var ErrTimeout = errors.New("command timeout")
var ErrClosing = errors.New("connection closing")

func (c *Conn) ExecCommand(commander Commander) (Result, error) {
	command, err := commander.Command()
	if err != nil {
		return Result{}, err
	}
	if _, ok := c.Device.Support[command.Method]; !ok {
		return Result{}, fmt.Errorf("command method %s not supported", command.Method)
	}
	c.commandC <- command
	result := <-command.resultC
	return result, nil
}
