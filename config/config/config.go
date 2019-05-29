package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"time"

	"github.com/jurgen-kluft/go-home/config"
	logpkg "github.com/jurgen-kluft/go-home/logging"
	"github.com/jurgen-kluft/go-home/pubsub"
)

// Configs holds all the config objects that we can have
type context struct {
	log     *logpkg.Logger
	pubsub  *pubsub.Context
	configs *configs
	watcher *configFileWatcher
}

func newContext(emitter map[string]string) *context {
	ctx := &context{}
	ctx.log = logpkg.New("configs")
	ctx.log.AddEntry("emitter")
	ctx.log.AddEntry("configs")
	ctx.pubsub = pubsub.New(emitter)
	ctx.watcher = newConfigFileWatcher()
	return ctx
}

type configs struct {
	Configurations map[string]*configuration `json:"configurations"`
}

type configuration struct {
	Name           string `json:"name"`
	ConfigFilename string `json:"filename"`
	ChannelName    string `json:"channel"`
	ReflectType    reflect.Type
}

func configFromJSON(jsonstr string) (*configs, error) {
	c := &configs{}
	err := json.Unmarshal([]byte(jsonstr), c)
	return c, err
}

func (c *context) initializeReflectTypes() {
	for name, configuration := range c.configs.Configurations {
		switch name {
		case "aqi":
			configuration.ReflectType = reflect.TypeOf(config.AqiConfig{})
		case "automation":
			configuration.ReflectType = reflect.TypeOf(config.AutomationConfig{})
		case "bravia.tv":
			configuration.ReflectType = reflect.TypeOf(config.BraviaTVConfig{})
		case "calendar":
			configuration.ReflectType = reflect.TypeOf(config.CalendarConfig{})
		case "flux":
			configuration.ReflectType = reflect.TypeOf(config.FluxConfig{})
		case "hue":
			configuration.ReflectType = reflect.TypeOf(config.HueConfig{})
		case "huebridge":
			configuration.ReflectType = reflect.TypeOf(config.HueBridgeConfig{})
		case "presence":
			configuration.ReflectType = reflect.TypeOf(config.PresenceConfig{})
		case "samsung.tv":
			configuration.ReflectType = reflect.TypeOf(config.SamsungTVConfig{})
		case "shout":
			configuration.ReflectType = reflect.TypeOf(config.ShoutConfig{})
		case "suncalc":
			configuration.ReflectType = reflect.TypeOf(config.SuncalcConfig{})
		case "weather":
			configuration.ReflectType = reflect.TypeOf(config.WeatherConfig{})
		case "wemo":
			configuration.ReflectType = reflect.TypeOf(config.WemoConfig{})
		case "xiaomi":
			configuration.ReflectType = reflect.TypeOf(config.XiaomiConfig{})
		case "yee":
			configuration.ReflectType = reflect.TypeOf(config.YeeConfig{})
		}
	}
}

func (c *context) initializeConfigFileWatcher() {
	for name, configuration := range c.configs.Configurations {
		c.watcher.watchConfigFile(configuration.ConfigFilename, name)
	}
}

func (c *context) updateConfigFileWatcher() {
	events := c.watcher.update()
	for _, event := range events {
		if event.Event == MODIFIED {
			_, exists := c.configs.Configurations[event.User]
			if exists {
				c.sendConfigOnChannel(event.User)
			}
		}
	}
}

func (c *context) checkAllConfigurationFiles() (err error) {
	for _, configuration := range c.configs.Configurations {
		var data []byte
		data, err = ioutil.ReadFile(configuration.ConfigFilename)
		if err != nil {
			c.log.LogError("config", err.Error())
		}

		v := reflect.New(configuration.ReflectType).Elem().Interface().(config.Config)
		v, err = v.FromJSON(string(data))
		if err == nil {
			c.log.LogError("config", err.Error())
		}

		_, err := v.ToJSON()
		if err != nil {
			c.log.LogError("config", err.Error())
		}
	}
	return
}

// SendConfigOnChannel will load the JSON based config file and publish it onto pubsub
func (c *context) sendConfigOnChannel(configtype string) (err error) {
	configuration, exists := c.configs.Configurations[configtype]
	if exists {
		var data []byte
		data, err = ioutil.ReadFile(configuration.ConfigFilename)
		if err != nil {
			return err
		}
		v := reflect.New(configuration.ReflectType).Elem().Interface().(config.Config)
		v, err = v.FromJSON(string(data))
		if err == nil {
			jsonstr, err := v.ToJSON()
			if err == nil {
				err = c.pubsub.Publish(configuration.ChannelName, jsonstr)
			}
		}
	}
	return
}

func main() {
	ctx := newContext(config.EmitterIOCfg)

	for {
		connected := true
		for connected {
			register := []string{"config/config/", "config/request/"}
			subscribe := []string{"config/config/", "config/request/"}
			err := ctx.pubsub.Connect("configs", register, subscribe)
			if err == nil {
				ctx.log.LogInfo("emitter", "connected")

				for connected {
					select {
					case msg := <-ctx.pubsub.InMsgs:
						topic := msg.Topic()
						if topic == "client/disconnected/" {
							ctx.log.LogInfo("emitter", "disconnected")
							connected = false
						} else if topic == "config/config/" {
							config, err := configFromJSON(string(msg.Payload()))
							if err == nil {
								ctx.log.LogInfo("configs", "received configuration")
								ctx.configs = config
								ctx.checkAllConfigurationFiles()
								ctx.initializeReflectTypes()
								ctx.initializeConfigFileWatcher()
							} else {
								ctx.log.LogError("configs", err.Error())
							}
						} else if topic == "config/request/" {
							configname := string(msg.Payload())
							ctx.log.LogInfo("configs", "requested configuration for '"+configname+"'.")
							ctx.sendConfigOnChannel(configname)
						}
						break
					case <-time.After(time.Second * 10):
						// Any config files updated ?
						ctx.updateConfigFileWatcher()
						break
					}
				}
			} else {
				connected = false
			}

			if err != nil {
				ctx.log.LogError("configs", err.Error())
			}
		}

		time.Sleep(5 * time.Second)
	}
}
