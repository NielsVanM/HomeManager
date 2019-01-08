package main

import (
	"flag"
	"os"
	"strings"

	"github.com/nielsvanm/homemanager/database"
	"github.com/nielsvanm/homemanager/frame"
	"github.com/nielsvanm/homemanager/plugin"
	"github.com/nielsvanm/homemanager/tools/log"
	"github.com/nielsvanm/homemanager/views"
)

var serverPort = 8080
var server *frame.WebServer
var db *database.DB

func init() {
	// General app setup
	// Setup database for webapp
	db = database.NewDB("postgres", "SuperSecure8", "homemanager", "127.0.0.1", 5432)
	db.Connect()

	// Setup pluginmanager and plugins
	plugin.PluginManager.DB = db
	plugin.PluginManager.Setup()
	database.Database = db

	// Create database tables
	queries := plugin.PluginManager.GetSetupQueries()
	db.CreateTables(queries)

	// Create server manager
	server = frame.NewWebServer()

	// Setup global enpoints
	server.RegisterEndpoint("/", views.DashboardView)
	server.RegisterEndpoint("/stats/", views.StatisticsView)
	server.RegisterEndpoint("/stats/processorcount/", views.ProcessorCountView)
	server.RegisterEndpoint("/stats/memory/", views.MemoryStatView)
	server.RegisterEndpoint("/stats/plugincount/", views.PluginCountView)
	server.RegisterEndpoint("/stats/logsize/", views.LogSizeView)
	server.RegisterEndpoint("/database/", views.DatabaseView)
	server.RegisterEndpoint("/database/create/{pluginname}/", views.CreateTablesView)
	server.RegisterEndpoint("/database/drop/{pluginname}/", views.DropTablesView)

	// Setup plugin endpoints
	server.AddEndpoints(
		plugin.PluginManager.GetEndpoints(),
	)
}

func main() {
	str := flag.String("runplugin", "", "--runplugin <pluginname>")

	flag.Parse()

	if *str != "" {
		RunSinglePlugin(*str, db)
	}

	// Start webserver
	server.Run(serverPort)
}

// RunSinglePlugin runs a single plugin for one iteration
func RunSinglePlugin(name string, db *database.DB) {

	for _, plugin := range plugin.PluginManager.Plugins {
		if strings.ToLower(plugin.Name) == name {
			batches := plugin.Main()

			for _, batch := range batches {
				db.ExecBatch(batch)
			}

			os.Exit(0)
		}
	}
	log.Err("SinglePuginRunner", "Failed to find plugin with name "+name)
	os.Exit(-1)
}
