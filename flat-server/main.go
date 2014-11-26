package main

import (
  "flag"
  "fmt"
  "html/template"
  "io"
  "log"
  "net/http"
  "path/filepath"
)

var (
  //assetsDir = flag.String("assets", "github.com/longda/flat/flat-server/assets", "assets directory")
  assetsDir = flag.String("assets", "assets", "assets directory")
  httpAddr = flag.String("http", ":8080", "Listen for HTTP connections on this address")
)

func main() {
  if err := parseHtmlTemplates([][]string{
      {"index.html"},
    }); err != nil {
    log.Fatal(err)
  }

  mux := http.NewServeMux()
  mux.HandleFunc("/", serveHome)

  if err := http.ListenAndServe(*httpAddr, mux); err != nil {
    log.Fatal(err)
  }
}

func serveHome(resp http.ResponseWriter, req *http.Request) {
  executeTemplate(resp, "index.html", map[string]interface{}{})
}

var templates = map[string]interface { Execute(io.Writer, interface {}) error}{}

func joinTemplateDir(base string, files[]string) []string {
  result := make([]string, len(files))
  for i := range files {
    result[i] = filepath.Join(base, "templates", files[i])
  }

  return result
}

func parseHtmlTemplates(sets [][]string) error {
  for _, set := range sets {
    templateName := set[0]
    t := template.New("")
    if _, err := t.ParseFiles(joinTemplateDir(*assetsDir, set)...); err != nil {
      return err
    }

    templates[templateName] = t
  }
  return nil
}

func executeTemplate(resp http.ResponseWriter, name string, data interface{}) {
  t := templates[name]
  if t == nil {
    fmt.Errorf("Template %s not found", name)
  }

  t.Execute(resp, data);
}