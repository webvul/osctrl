package main

import (
	"fmt"
	"path/filepath"
	"plugin"
)

// Variables for plugin functions
var (
	logsDispatcher func(logging, logType string, params ...interface{})
)

// Function to load logging dispatcher plugin
func loadLoggingDispatcherPlugin() error {
	plugins, err := filepath.Glob("plugins/logging_dispatcher_plugin.so")
	if err != nil {
		return err
	}
	p, err := plugin.Open(plugins[0])
	if err != nil {
		return err
	}
	symbolLogsDispatcher, err := p.Lookup("LogsDispatcher")
	if err != nil {
		return err
	}
	var ok bool
	logsDispatcher, ok = symbolLogsDispatcher.(func(logging, logType string, params ...interface{}))
	if !ok {
		return fmt.Errorf("Plugin has no 'LogsDispatcher' function")
	}
	return nil
}