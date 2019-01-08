package views

import (
	"net/http"

	"github.com/nielsvanm/homemanager/database"

	"github.com/gorilla/mux"
	"github.com/nielsvanm/homemanager/frame"
	"github.com/nielsvanm/homemanager/plugin"
)

// DatabaseView is an overview and management page for the database tables
func DatabaseView(w http.ResponseWriter, r *http.Request) {
	dbPage := frame.NewPage([]string{"base.html", "database/dashboard.html"})

	// Get all plugins and add it to context
	plugins := plugin.PluginManager.Plugins
	dbPage.AddContext("plugins", plugins)

	dbPage.Render(w)
}

func CreateTablesView(w http.ResponseWriter, r *http.Request) {
	pluginName := mux.Vars(r)["pluginname"]

	plug := plugin.PluginManager.GetPlugin(pluginName)
	if plug == nil {
		w.Write([]byte("Failed to find the plugin"))
		return
	}

	for _, query := range plug.SetupDatabase {
		database.Database.Exec(query)
	}
}

func DropTablesView(w http.ResponseWriter, r *http.Request) {
	pluginName := mux.Vars(r)["pluginname"]

	plug := plugin.PluginManager.GetPlugin(pluginName)
	if plug == nil {
		w.Write([]byte("Failed to find the plugin"))
		return
	}

	queries := plug.DropTables()

	for _, query := range queries {
		database.Database.Exec(query)
	}
}
