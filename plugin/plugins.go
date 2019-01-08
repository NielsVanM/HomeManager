package plugin

import (
	"nvmtech.nl/homemanager/plugin/torrentplugin"
	"nvmtech.nl/homemanager/plugin/ytsamplugin"
)

// PluginList is a collection of all registered plugins
func init() {
	PluginManager.Plugins = []*Plugin{
		&Plugin{
			"YTSAMPlugin",
			"Pulls movies from YTS.AM and allows you to download them",
			"Entertainment",
			ytsamplugin.SetupDB,
			ytsamplugin.Tables,
			ytsamplugin.APIEndpoints,
			ytsamplugin.ViewEndpoints,
			ytsamplugin.GetMovies,
			[]string{"torrents"},
		},
		&Plugin{
			"TorrentPlugin",
			"Download torrents",
			"Internet",
			torrentplugin.SetupDB,
			torrentplugin.Tables,
			torrentplugin.APIEndpoints,
			torrentplugin.ViewEndpoints,
			torrentplugin.UpdateTorrents,
			[]string{"torrents", "downloads"},
		},
	}
}
