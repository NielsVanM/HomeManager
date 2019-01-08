package frame

import (
	"html/template"
	"io"
	"strings"

	"github.com/nielsvanm/homemanager/tools/log"
)

// TemplateFolder is the folder where the templates are located,
// typically injected by the main.go file
var TemplateFolder string

func init() {
	TemplateFolder = "./templates/"
}

// Page structure that keeps the data of a page, these are used for rendering
// pages.
type Page struct {
	Template *template.Template
	Context  map[string]interface{}
}

// NewPage creates a new page in memory, it tries to parse the parent and body
// html templates and if there is an error it returns no page and an error.
// If it was succesfull the function populates Context with an empty
// map[string]interface{}
func NewPage(pages []string) *Page {
	p := Page{}
	// Add template folder to pages
	for i := 0; i < len(pages); i++ {
		pages[i] = TemplateFolder + pages[i]
	}

	// Try to parse the files
	var err error
	p.Template, err = template.ParseFiles(pages...)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			log.Fatal("PageParser", "Can't find file "+err.Error())
		}
		log.Err("PageParser", "Failed to parse template "+err.Error())
	}
	p.Context = map[string]interface{}{}

	return &p
}

// AddContext adds context to the page, this is passed to the templates when
// rendering the page. It takes an key and value to be used as context.
func (p *Page) AddContext(key string, value interface{}) {
	p.Context[key] = value
}

// CleanContext deletes all of the current context.
func (p *Page) CleanContext() {
	p.Context = map[string]interface{}{}
}

// Render renders the page to the io.Writer that is passed to it. It includes
// the pages Context to load data into.
// After rendering the page it cleans the context for the next request
func (p *Page) Render(w io.Writer) {
	err := p.Template.Execute(w, p.Context)
	if err != nil {
		log.Err("PageParser", err.Error())
	}
	p.CleanContext()
}
