package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"
	"path/filepath"
)

var (
	configDir = flag.String("config-dir", "config", "Path to the directory where configuration is stored")
)

func main() {
	flag.Parse()
	glog.Infoln("Loading configuration file")
	configPaths := []string{filepath.Join(*configDir, "config.json")}
	var config *configuration
	var err error
	for _, cfgpth := range configPaths {
		config, err = loadConfig(cfgpth)
		if err != nil && os.IsNotExist(err) {
			// The configuration file we just attempted to load
			// does not exist. No issue; we will just move on
			// to the next one.
			continue
		} else if err != nil {
			glog.Fatalln(err)
		} else if err == nil {
			break
		}
	}
	if config == nil {
		// *WELP*, it looks like we couldn't find a configuration file.
		glog.Fatalln("No config.json or config.yml file found")
	}

	glog.Infoln("Loading process configuration files")
	processConfigDir := filepath.Join(*configDir, "process.d")
	processConfigs := make([]string, 0)
	globs := []string{filepath.Join(processConfigDir, "*.json")}
	for _, pat := range globs {
		matches, err := filepath.Glob(pat)
		if err != nil {
			glog.Fatalln(err)
		}
		processConfigs = append(processConfigs, matches...)
	}
	processes := make([]*process, len(processConfigs))
	for i, pconfig := range processConfigs {
		if proc, err := loadProcess(pconfig, *config); err != nil {
			glog.Fatalln(err)
		} else {
			processes[i] = proc
		}
	}

	glog.Infoln("Starting processes")
	for _, proc := range processes {
		glog.V(1).Infoln("|-", proc.Name)
		go func(p *process) {
			if err := p.start(); err != nil {
				glog.Warningf("%s: %v", p.Name, err)
			}
		}(proc)
	}

	if config.HTTP.Enabled {
		glog.Infoln("Starting HTTP server; listening on", config.HTTP.ListenAddr)
		http.HandleFunc("/", handleHTTPStatus)
		if err := http.ListenAndServe(config.HTTP.ListenAddr, nil); err != nil {
			glog.Fatalln(err)
		}
	}
	glog.Infoln("Shutting down")
}
