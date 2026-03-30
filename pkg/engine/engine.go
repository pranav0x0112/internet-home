package engine

import (
	"bytes"
	"html/template"
	"math/rand"
	"os"
	fp "path/filepath"
	"strings"
	"time"

	"github.com/anna-ssg/anna/v3/pkg/logger"
	"github.com/anna-ssg/anna/v3/pkg/parser"
)

// DeepDataMerge This struct holds all the ssg data
type DeepDataMerge struct {
	// Templates stores the template data of all the pages of the site
	// Access the data for a particular page by using the relative path to the file as the key
	Templates map[template.URL]parser.TemplateData

	// Templates stores the template data of all tag sub-pages of the site
	Tags map[template.URL]parser.TemplateData

	// K-V pair storing all templates corresponding to a particular tag in the site
	TagsMap map[template.URL][]parser.TemplateData

	// Stores data parsed from layout/config.yml
	LayoutConfig parser.LayoutConfig

	// Templates stores the template data of all collection sub-pages of the site
	Collections map[template.URL]parser.TemplateData

	// K-V pair storing all templates corresponding to a particular collection in the site
	CollectionsMap map[template.URL][]parser.TemplateData

	// K-V pair storing the template layout name for a particular collection in the site
	CollectionsSubPageLayouts map[template.URL]string

	// Stores the index generated for search functionality
	JSONIndex map[template.URL]JSONIndexTemplate
}

type Engine struct {
	// Stores the merged ssg data
	DeepDataMerge DeepDataMerge

	// Common logger for all engine functions
	ErrorLogger *logger.Logger

	// The path to the directory being rendered
	SiteDataPath string
}

type PageData struct {
	DeepDataMerge DeepDataMerge

	PageURL       template.URL
	Image         string
	GalleryImages []string
	IsHome        bool
}

// JSONIndexTemplate This structure is solely used for storing the JSON index
type JSONIndexTemplate struct {
	CompleteURL template.URL
	Frontmatter parser.Frontmatter
	Tags        []string
}

/*
RenderPage
fileOutPath - stores the parent directory to store rendered files, usually `site/`

pagePath - stores the path to write the given page without the prefix directory
Eg: site/content/posts/file1.html to be passed as posts/file1.html

template - stores the HTML templates parsed from the layout/ directory

templateStartString - stores the name of the template to be passed to ExecuteTemplate()
*/
func (e *Engine) RenderPage(fileOutPath string, pagePath template.URL, template *template.Template, templateStartString string) {
	outputPath := fileOutPath + "rendered/" + string(pagePath)
	if strings.HasPrefix(string(pagePath), "blogs/") && strings.HasSuffix(string(pagePath), ".html") && string(pagePath) != "blogs/index.html" {
		blogPath := strings.TrimSuffix(string(pagePath), ".html")
		if err := os.MkdirAll(fileOutPath+"rendered/"+blogPath, 0750); err != nil {
			e.ErrorLogger.Fatal(err)
		}
		outputPath = fileOutPath + "rendered/" + blogPath + "/index.html"
		_ = os.Remove(fileOutPath + "rendered/" + string(pagePath))
	} else if strings.Contains(string(pagePath), "/") {
		if string(pagePath) == "blogs/index.html" {
			_ = os.RemoveAll(fileOutPath + "rendered/blogs/index")
		}
		// Creating subdirectories if the filepath contains '/'
		splitPaths := strings.Split(string(pagePath), "/")
		filename := splitPaths[len(splitPaths)-1]
		pagePathWithoutFilename, _ := strings.CutSuffix(string(pagePath), filename)

		err := os.MkdirAll(fileOutPath+"rendered/"+pagePathWithoutFilename, 0750)
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}
	}
	var buffer bytes.Buffer

	pageData := PageData{
		DeepDataMerge: e.DeepDataMerge,
		PageURL:       pagePath,
		Image:         e.DeepDataMerge.Templates[pagePath].Frontmatter.Image,
	}
	if string(pagePath) == "index.html" {
		pageData.IsHome = true
	}

	if string(pagePath) == "gallery/index.html" {
		root, _ := os.Getwd()
		galleryPath := fp.Join(root, "assets", "images", "gallery")
		entries, err := os.ReadDir(galleryPath)
		if err != nil {
			e.ErrorLogger.Fatal(err)
		}

		allowedExts := map[string]struct{}{
			".jpg":  {},
			".jpeg": {},
			".png":  {},
			".webp": {},
		}

		images := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			ext := strings.ToLower(fp.Ext(entry.Name()))
			if _, ok := allowedExts[ext]; !ok {
				continue
			}

			images = append(images, entry.Name())
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(images), func(i, j int) {
			images[i], images[j] = images[j], images[i]
		})

		pageData.GalleryImages = images
	}

	// Storing the rendered HTML file to a buffer
	err := template.ExecuteTemplate(&buffer, templateStartString, pageData)
	if err != nil {
		e.ErrorLogger.Println("Error at path: ", pagePath)
		e.ErrorLogger.Fatal(err)
	}

	// Flushing data from the buffer to the disk
	err = os.WriteFile(outputPath, buffer.Bytes(), 0666)
	if err != nil {
		e.ErrorLogger.Fatal(err)
	}
}
