package views

import (
	"net/http"

	"github.com/nielsvanm/homemanager/frame"
)

// DashboardView is the main index page of the site
func DashboardView(w http.ResponseWriter, r *http.Request) {
	dashboardPage := frame.NewPage([]string{"base.html", "dashboard/dashboard.html"})

	dashboardPage.Render(w)
}
