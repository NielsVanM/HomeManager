package views

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"nvmtech.nl/homemanager/plugin"
	"nvmtech.nl/homemanager/tools"
	"nvmtech.nl/homemanager/tools/log"

	"nvmtech.nl/homemanager/frame"
)

func StatisticsView(w http.ResponseWriter, r *http.Request) {
	// Construct page
	page := frame.NewPage([]string{"base.html", "dashboard/statistics.html"})

	// Render page
	page.Render(w)
}

// ProcessorCountView returns a integer representing the amount of cores available
func ProcessorCountView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(strconv.Itoa(
		runtime.NumCPU(),
	)))
}

// MemoryStatView returns the human readable memory count
func MemoryStatView(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	usage := tools.ByteCountDecimal(int64(m.Alloc))

	w.Write([]byte(usage))
}

// PluginCountView returns the amount of plugins that are registered
func PluginCountView(w http.ResponseWriter, r *http.Request) {
	pluginCount := len(plugin.PluginManager.Plugins)

	w.Write(
		[]byte(strconv.Itoa(
			pluginCount,
		)),
	)
}

// LogSizeView writes the size of log.txt to the responsewriter
func LogSizeView(w http.ResponseWriter, r *http.Request) {
	f, err := os.Stat("./log.txt")
	if err != nil {
		log.Warn("Statistics", "Failed to access log.txt", err.Error())
	}

	size := f.Size()
	readableSize := tools.ByteCountDecimal(size)

	w.Write([]byte(readableSize))
}
