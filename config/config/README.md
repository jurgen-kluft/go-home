# Config service

This package implements a small microservice that distributes JSON configuration
files onto pubsub channels and watches those files for changes.

Purpose
- Provide a single place to register named configuration files and the
	pubsub channel each configuration should be published on.
- Load configuration files, validate/convert them via the top-level
	`github.com/jurgen-kluft/go-home/config` parsers and publish the resulting
	JSON onto the configured channel.
- Watch registered configuration files for filesystem changes and publish
	updated configuration automatically.

Files
- [config.go](config/config.go): main microservice logic. Registers handlers
	for `config/config` (load master config), `config/request` (request single
	config) and `tick` (drive the file watcher). Contains helpers to parse the
	master configuration, register channels and publish config payloads.
- [configwatch.go](config/configwatch.go): lightweight file watcher and event
	detection used by the service. It exposes simple event constants such as
	`CREATED`, `MODIFIED`, `DELETED`, `NOEXIST`, etc.

Master config format
The service expects a master configuration JSON describing all managed
configurations. The JSON structure is:

```
{
	"configurations": {
		"aqi": {
			"name": "aqi",
			"filename": "config/aqi.config.json",
			"channel": "config/aqi"
		},
		"automation": {
			"name": "automation",
			"filename": "config/automation.config.json",
			"channel": "config/automation"
		}
	}
}
```

Each entry's `filename` should point to a JSON file on disk. The service uses
the `config` package's specific parsers (for example `AqiConfigFromJSON`,
`AutomationConfigFromJSON`, etc.) to turn those files into typed config objects
before publishing.

Behavior
- When the service receives the master JSON on the `config/config` channel it
	loads and registers all configured files, validates/parses each file and
	registers their pubsub channels.
- When a message is published to `config/request` with a configuration name
	(e.g. `"aqi"`) the service will publish the most recent JSON for that
	configuration on the configured channel (for the example above: `config/aqi`).
- The service periodically checks the watched files (via the `tick` handler).
	When a file is created/modified/deleted the watcher will emit an event and
	the service will re-load and re-publish the updated configuration.

Usage
- Build/run the service from the repository root:

```bash
go run ./config/config
```

- Provide the master configuration by publishing the JSON to the
	`config/config` channel (using whatever pubsub/microservice tooling you
	have in your environment). Once received, the service will register and
	start publishing the individual configurations on their channels.

Development notes
- Parsers and per-config validation live in the top-level `config` package and
	are used by `configFromJSON` in [config.go](config/config.go).
- The watcher is intentionally simple and suitable for services that can be
	polled at a reasonable interval (the service registers a `tick` handler to
	drive updates).

If you want changes or additional examples (e.g. concrete per-config JSON
examples), tell me which config type you'd like documented and I will add it.

