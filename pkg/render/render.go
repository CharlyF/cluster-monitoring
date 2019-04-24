package render

import (
	"bytes"
	"encoding/json"
	"io"
	"fmt"
	"html/template"
	"github.com/CharlyF/cluster-monitoring/util"
	"path/filepath"
)

var (
	here, _        = util.Folder()
	fmap           = Fmap()
	templateFolder string
)

func init() {
	templateFolder = filepath.Join("/go/src/github.com/CharlyF/cluster-monitoring", "templates")
}

func FormatData(data []byte) (string, error) {
	var b= new(bytes.Buffer)

	stats := make(map[string]interface{})
	json.Unmarshal(data, &stats)
	renderData(b, stats)
	return b.String(), nil
}


func renderData(w io.Writer, stats map[string]interface{}) {
	t := template.Must(template.New("data.tmpl").Funcs(fmap).ParseFiles(filepath.Join(templateFolder, "data.tmpl")))
	err := t.Execute(w, stats)
	if err != nil {
		fmt.Println(err)
	}
}
