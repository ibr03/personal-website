package builder

import (
	"io"
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/adrg/frontmatter"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Meta  map[string]interface{}
	Content template.HTML
}

func Build() error {
	// clean build/
	os.RemoveAll("build")
	os.MkdirAll("build", 0755)

	// load template
	tmpl, err := template.ParseFiles("templates/base.html")
	if err != nil {
		return err
	}

	// read content dir
	files, err := os.ReadDir("pages")
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".md" {
			continue
		}

		inputPath := filepath.Join("pages", file.Name())
		data, err := os.ReadFile(inputPath)

		if err != nil {
			return err
		}

		meta, body, err := parseMarkdown(data)
		if err != nil {
			return err
		}
		html := markdownToHTML(body)

		page := Page{
			Meta:    meta,
			Content: template.HTML(html),
		}

		// output path
		name := strings.TrimSuffix(file.Name(), ".md")
		var outDir string
		if name == "index" {
			outDir = "build"
		} else {
			outDir = filepath.Join("build", name)
		}
		os.MkdirAll(outDir, 0755)

		outFile := filepath.Join(outDir, "index.html")
		f, _ := os.Create(outFile)
		defer f.Close()

		err = tmpl.Execute(f, page)
		if err != nil {
			return err
		}
	}

	// copy static files
	err = copyDir("static", "build")
	if err != nil {
		return err
	}

	return nil
}

func parseMarkdown(input []byte) (map[string]interface{}, string, error) {
	var meta map[string]interface{}

	body, err := frontmatter.Parse(bytes.NewReader(input), &meta)
	if err != nil {
		return nil, "", err
	}

	return meta, string(body), nil
}

func markdownToHTML(md string) string {
	var buf bytes.Buffer
	goldmark.Convert([]byte(md), &buf)
	return buf.String()
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		
		// copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}