package plugin

import (
	"fmt"
	"os"
	"strings"

	"nvmtech.nl/homemanager/database"
	"nvmtech.nl/homemanager/frame"

	"nvmtech.nl/homemanager/tools"
	"nvmtech.nl/homemanager/tools/log"
)

// PluginManager is the app-wide plugin management object
var PluginManager = Manager{[]*Plugin{}, nil}

// Plugin is a type that represents actions that have to be executed on the server
// it wraps logic and settings that specify it's behaviour
type Plugin struct {
	// Plugin information
	Name        string
	Description string
	Category    string

	// Database queries
	SetupDatabase []string
	Tables        []string

	// Endpoints
	APIEndpoints  []*frame.Endpoint
	ViewEndpoints []*frame.Endpoint

	// Main Function
	Main func() []database.BatchQuery

	// Data dirs
	DataDirs []string
}

// Setup adds the name of the plugin at any %s that is provided
// by the plugin
func (p *Plugin) Setup() {
	// Add name to database queries
	for i := 0; i < len(p.SetupDatabase); i++ {
		p.SetupDatabase[i] = p.AddPluginNameToQuery(p.SetupDatabase[i])
	}
	for i := 0; i < len(p.Tables); i++ {
		p.Tables[i] = p.AddPluginNameToQuery(p.Tables[i])
	}

	// Add /api/{pluginname}/ to api endpoints
	newEndpoints := []*frame.Endpoint{}
	for _, endp := range p.APIEndpoints {
		endp.URL = "/api/" + strings.ToLower(p.Name) + endp.URL
		newEndpoints = append(newEndpoints, endp)
	}
	p.APIEndpoints = newEndpoints

	// Add {pluginname}/ to api endpoints
	newEndpoints = []*frame.Endpoint{}
	for _, endp := range p.ViewEndpoints {
		endp.URL = "/" + strings.ToLower(p.Name) + endp.URL
		newEndpoints = append(newEndpoints, endp)
	}
	p.ViewEndpoints = newEndpoints

	// Create data dirs
	for _, dir := range p.DataDirs {
		err := os.MkdirAll(p.GetDir(dir), os.ModePerm)
		if err != nil {
			if strings.Contains(err.Error(), "exists") {
				continue
			}
			log.Warn(p.Name, "Failed to create folder", err.Error())
		}
	}
}

// AddPluginNameToQuery replaces all the %s in the query with the plugin
// name
func (p *Plugin) AddPluginNameToQuery(query string) string {
	query = strings.Replace(query, "%s", "%[1]s", -1)
	return fmt.Sprintf(query, p.Name)
}

// DropTables provides a batch query for dropping the tables
func (p *Plugin) DropTables() []string {
	queries := []string{}

	for _, table := range p.Tables {
		query := `DROP TABLE %s CASCADE;`
		queries = append(
			queries,
			fmt.Sprintf(query, table),
		)
	}

	return queries
}

// GetDir returns the path of a plugin folder
func (p *Plugin) GetDir(name string) string {
	return "./__data/" + strings.ToLower(p.Name) + "/" + name
}

// Manager is a management object for the plugin struct
type Manager struct {
	Plugins []*Plugin
	DB      *database.DB
}

// Setup runs initial functionality of plugins to ensure they are operationalIn
func (m *Manager) Setup() {
	for _, plugin := range m.Plugins {
		plugin.Setup()
		log.Info("PluginManager", "Finished setting up "+plugin.Name)
	}
}

// GetSetupQueries returns the queries necessarry to setup the database
func (m *Manager) GetSetupQueries() []string {
	allQueries := []string{}

	for _, plugin := range m.Plugins {
		allQueries = append(allQueries, plugin.SetupDatabase...)
	}

	return allQueries
}

// GetEndpoints returns a list of all the endpoints any plugin has registered
func (m *Manager) GetEndpoints() []*frame.Endpoint {
	endpoints := []*frame.Endpoint{}

	for _, plugin := range m.Plugins {
		endpoints = append(endpoints, plugin.APIEndpoints...)
		endpoints = append(endpoints, plugin.ViewEndpoints...)
	}

	return endpoints
}

// GetCategories returns a list of categories that are registered by the plugin
// it also removes the duplicates from the list
func (m *Manager) GetCategories() []string {
	categories := []string{}

	for _, plugin := range m.Plugins {
		if tools.IsInList(plugin.Category, categories) {
			continue
		}
		categories = append(categories, plugin.Category)
	}

	return categories
}

func (m *Manager) RunPlugins() {
	for _, plugin := range m.Plugins {
		batches := plugin.Main()
		for _, batch := range batches {
			m.DB.ExecBatch(batch)
		}
	}
}

func (m *Manager) GetPlugin(pluginName string) *Plugin {
	for _, plugin := range m.Plugins {
		if plugin.Name == pluginName {
			return plugin
		}
	}

	return nil
}
