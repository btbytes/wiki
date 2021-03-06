package wiki

import (
	"html/template"
	"net/http"

	"github.com/peterhellberg/wiki/db"
	"github.com/russross/blackfriday"
	"github.com/zenazn/goji/web"
)

func (w *Wiki) Show(c web.C, rw http.ResponseWriter, r *http.Request) {
	name := w.getPageName(c.URLParams["name"])

	w.DB().View(func(tx *db.Tx) error {
		p, err := tx.Page(name)

		if err != nil {
			p.Text = []byte(emptyPageString)
		}

		vars := map[string]interface{}{
			"Path": "/" + string(p.Name) + "/edit",
			"Text": BytesAsHTML(ParsedMarkdown(p.Text)),
		}

		t := template.Must(template.New("show").Parse(showTpl))
		t.Execute(rw, vars)

		return nil
	})
}

func (w *Wiki) RedirectToShow(c web.C, rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/"+c.URLParams["name"], 302)
}

func BytesAsHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

func ParsedMarkdown(b []byte) []byte {
	return blackfriday.MarkdownCommon(b)
}
