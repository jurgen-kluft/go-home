package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/jurgen-kluft/go-home/config"
	"github.com/jurgen-kluft/go-home/micro-service"
)

// Configs holds all the config objects that we can have
type context struct {
	configs *configs
	watcher *configFileWatcher
	micro   *microservice.Service
}

func newContext() *context {
	ctx := &context{}
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
}

func configFromJSON(data []byte) (*configs, error) {
	c := &configs{}
	err := json.Unmarshal(data, c)
	return c, err
}

func (c *context) configFromJSON(configname string, jsondata []byte) (config.Config, error) {
	var ci config.Config
	var err error
	c.micro.Logger.LogInfo(c.micro.Name, fmt.Sprintf("configuration %s, FromJSON", configname))
	switch configname {
	case "aqi":
		ci, err = config.AqiConfigFromJSON(jsondata)
	case "automation":
		ci, err = config.AutomationConfigFromJSON(jsondata)
	case "bravia.tv":
		ci, err = config.BraviaTVConfigFromJSON(jsondata)
	case "calendar":
		ci, err = config.CalendarConfigFromJSON(jsondata)
	case "flux":
		ci, err = config.FluxConfigFromJSON(jsondata)
	case "hue":
		ci, err = config.HueConfigFromJSON(jsondata)
	case "huebridge":
		ci, err = config.HueBridgeConfigFromJSON(jsondata)
	case "presence":
		ci, err = config.PresenceConfigFromJSON(jsondata)
	case "samsung.tv":
		ci, err = config.SamsungTVConfigFromJSON(jsondata)
	case "shout":
		ci, err = config.ShoutConfigFromJSON(jsondata)
	case "suncalc":
		ci, err = config.SuncalcConfigFromJSON(jsondata)
	case "weather":
		ci, err = config.WeatherConfigFromJSON(jsondata)
	case "wemo":
		ci, err = config.WemoConfigFromJSON(jsondata)
	case "xiaomi":
		ci, err = config.XiaomiConfigFromJSON(jsondata)
	case "yee":
		ci, err = config.YeeConfigFromJSON(jsondata)
	}
	return ci, err
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
	for name, configuration := range c.configs.Configurations {
		var data []byte
		data, err = ioutil.ReadFile(configuration.ConfigFilename)
		if err != nil {
			c.micro.Logger.LogError(c.micro.Name, err.Error())
		} else {
			if data != nil {
				v, err := c.configFromJSON(name, data)
				if err != nil {
					c.micro.Logger.LogError(c.micro.Name, err.Error())
				} else {
					data, err = v.ToJSON()
					if err != nil {
						c.micro.Logger.LogError(c.micro.Name, err.Error())
					}
				}
			} else {
				c.micro.Logger.LogError(c.micro.Name, fmt.Sprintf("Configuration %s did not have a ReflectType", name))
			}
		}
	}
	return
}

// SendConfigOnChannel will load the JSON based config file and publish it onto pubsub
func (c *context) sendConfigOnChannel(configtype string) (err error) {
	if c.configs != nil {
		configuration, exists := c.configs.Configurations[configtype]
		if exists {
			var configJSONData []byte
			configJSONData, err = ioutil.ReadFile(configuration.ConfigFilename)
			if err != nil {
				return err
			}
			if configJSONData != nil {
				v, err := c.configFromJSON(configtype, configJSONData)
				if err == nil {
					jsondata, err := v.ToJSON()
					if err == nil {
						c.micro.Logger.LogInfo(c.micro.Name, fmt.Sprintf("Publish %s on channel %s", string(jsondata), configuration.ChannelName))
						err = c.micro.Pubsub.Publish(configuration.ChannelName, string(jsondata))
					}
				}
			} else {
				err = fmt.Errorf("Configuration %s did not have JSON data or a ReflectType", configtype)
			}
		} else {
			err = fmt.Errorf("Configuration %s does not exist", configtype)
		}
	} else {
		err = fmt.Errorf("Haven't received configuration, so cannot send configuration requests")
	}
	return
}

func main() {
	register := []string{"config/config/", "config/request/", "config/presence/", "config/aqi/"}
	subscribe := []string{"config/config/", "config/request/"}

	m := microservice.New("config")
	m.RegisterAndSubscribe(register, subscribe)

	ctx := newContext()
	ctx.micro = m

	m.RegisterHandler("config/config/", func(m *microservice.Service, topic string, msg []byte) bool {
		config, err := configFromJSON(msg)
		if err == nil {
			m.Logger.LogInfo(m.Name, "received configuration")
			ctx.configs = config
			ctx.checkAllConfigurationFiles()
			ctx.initializeConfigFileWatcher()
		} else {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	m.RegisterHandler("config/request/", func(m *microservice.Service, topic string, msg []byte) bool {
		configname := string(msg)
		m.Logger.LogInfo(m.Name, "requested configuration for '"+configname+"'.")
		err := ctx.sendConfigOnChannel(configname)
		if err != nil {
			m.Logger.LogError(m.Name, err.Error())
		}
		return true
	})

	tickCount := 0
	m.RegisterHandler("tick/", func(m *microservice.Service, topic string, msg []byte) bool {
		if tickCount == 5 {
			tickCount = 0
			// Any config files updated ?
			ctx.updateConfigFileWatcher()
		} else {
			tickCount += 1
		}
		return true
	})

	m.Loop()
}
