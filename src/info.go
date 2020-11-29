package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type asamblea struct {
	Name string
	Date string
}

type actividad struct {
	Name string
	Year string
}

func (s *server) infoHandler(w http.ResponseWriter, r *http.Request) {
	s.Log(r, LOG_DEBUG, "Page /info")

	p := make(map[string]interface{})

	p["board"] = s.DBlistBoard()
	p["asambleas"] = []asamblea{}
	p["actividades"] = []actividad{}

	var files []os.FileInfo
	var err error

	files, err = ioutil.ReadDir("../katiuskas/public/Asambleas")
	if err == nil {
		for _, file := range files {
			name := file.Name()
			n := strings.TrimSuffix(name, ".pdf")
			d := strings.Split(n, "-")
			if len(d) != 4 {
				continue
			}
			p["asambleas"] = append(p["asambleas"].([]asamblea), asamblea{
				Name: file.Name(),
				Date: fmt.Sprintf("%s-%s-%s", d[3], d[2], d[1]),
			})
		}
	}

	files, err = ioutil.ReadDir("../katiuskas/public/Memorias-de-actividades")
	if err == nil {
		for _, file := range files {
			name := file.Name()
			n := strings.TrimSuffix(name, ".pdf")
			d := strings.Split(n, "-")
			if len(d) != 3 {
				continue
			}
			p["actividades"] = append(p["actividades"].([]actividad), actividad{
				Name: file.Name(),
				Year: d[2],
			})
		}
	}

	renderTemplate(w, r, "info", p)
}
