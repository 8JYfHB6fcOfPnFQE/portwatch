// Package config provides loading and validation of portwatch configuration
// files written in YAML.
//
// A configuration file may look like:
//
//	scan_interval: 15s
//	port_range:
//	  from: 1
//	  to: 10000
//	output:
//	  format: json
//	  file: /var/log/portwatch.log
//	rules_file: /etc/portwatch/rules.yaml
//
// Call [Load] to read a file, or [Default] to obtain sensible defaults
// without reading any file.
package config
