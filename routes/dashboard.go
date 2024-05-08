package routes

import (
	"CyberDefenseEd/QuadDB/database"
	"CyberDefenseEd/QuadDB/util"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func RenderTemplate(w http.ResponseWriter, templateName string, title string, data interface{}) {
	tmpl, err := template.ParseFiles(
		"./dashboard/templates/base.html",
		"./dashboard/templates/pages/"+templateName,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]interface{}{
		"title": title,
		"data":  data,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SetupDashboardRoutes(router *gin.Engine, dataDir string, aesKey []byte) {
	router.GET("/", func(c *gin.Context) {
		databases := make(map[string]*database.Database)

		dbFiles, err := filepath.Glob(filepath.Join(dataDir, "*.qdb"))
		if err != nil {
			util.Error("Failed to index qdb files.")
			http.Error(c.Writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for _, dbFile := range dbFiles {
			dbName := strings.TrimSuffix(filepath.Base(dbFile), ".qdb")
			db := database.LoadDB(dbFile, aesKey)
			databases[dbName] = db
		}

		RenderTemplate(c.Writer, "home.html", "Dashboard", gin.H{"databases": databases})
	})
}
