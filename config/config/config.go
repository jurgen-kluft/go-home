package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ctx.log.AddEntry("config")
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
}

func configFromJSON(data []byte) (*configs, error) {
	c := &configs{}
	err := json.Unmarshal(data, c)
	return c, err
}

func (c *context) configFromJSON(configname string, jsondata []byte) (config.Config, error) {
	var ci config.Config
	var err error
	c.log.LogInfo("config", fmt.Sprintf("configuration %s, FromJSON", configname))
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
			c.log.LogError("config", err.Error())
		} else {
			if data != nil {
				v, err := c.configFromJSON(name, data)
				if err != nil {
					c.log.LogError("config", err.Error())
				} else {
					data, err = v.ToJSON()
					if err != nil {
						c.log.LogError("config", err.Error())
					}
				}
			} else {
				c.log.LogError("config", fmt.Sprintf("Configuration %s did not have a ReflectType", name))
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
						c.log.LogInfo("config", fmt.Sprintf("Publish %s on channel %s", string(jsondata), configuration.ChannelName))
						err = c.pubsub.Publish(configuration.ChannelName, string(jsondata))
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
	ctx := newContext(config.PubSubCfg)

	for {
		connected := true
		for connected {
			register := []string{"config/config/", "config/request/", "config/presence/", "config/aqi/"}
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
							config, err := configFromJSON(msg.Payload())
							if err == nil {
								ctx.log.LogInfo("config", "received configuration")
								ctx.configs = config
								ctx.checkAllConfigurationFiles()
								ctx.initializeConfigFileWatcher()
							} else {
								ctx.log.LogError("config", err.Error())
							}
						} else if topic == "config/request/" {
							configname := string(msg.Payload())
							ctx.log.LogInfo("config", "requested configuration for '"+configname+"'.")
							err := ctx.sendConfigOnChannel(configname)
							if err != nil {
								ctx.log.LogError("config", err.Error())
							}
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
				ctx.log.LogError("config", err.Error())
			}
		}

		time.Sleep(5 * time.Second)
	}
}
