// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/insolar/insolar/insolar"
)

type Grapher interface {
	Add(insolar.Pulse, float64)
	Draw()
}

type StubDrawer struct{}

func (d StubDrawer) Add(insolar.Pulse, float64) {}

func (d StubDrawer) Draw() {}

type webGraph struct {
	Title string

	CurveType           string
	HorizontalAxisTitle string
	VerticalAxisTitle   string
	DataHeaders         []string
	Data                []tmplData
}

type tmplData struct {
	XValue  string
	YValues []float64
}

// based on https://developers-dot-devsite-v2-prod.appspot.com/chart/interactive/docs/gallery/linechart.html
// https://jsfiddle.net/api/post/library/pure/
var graphTmpl = template.Must(template.New("graphHtml").Parse(`
<html>
  <head>
    <script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
    <script type="text/javascript">
      google.charts.load('current', {'packages':['corechart']});
      google.charts.setOnLoadCallback(drawChart);

      function drawChart() {
        var data = google.visualization.arrayToDataTable([
		  // headers
          [ {{ range .DataHeaders }}'{{ . }}',{{ end }} ],
		  // values
		  {{ range .Data }}
          [{{ .XValue }}, {{ range .YValues }} {{ . }}, {{ end }} ],
		  {{ end }}
        ]);

        var options = {
          title: '{{ .Title }}',
          curveType: '{{ .CurveType }}',
          legend: { position: 'bottom' },
		  hAxis: { title: '{{ .HorizontalAxisTitle }}' },
          vAxis: { title: '{{ .VerticalAxisTitle }}' }
        };

        var chart = new google.visualization.LineChart(document.getElementById('chart'));

        chart.draw(data, options);
      }
    </script>
  </head>
  <body>
    <div id="chart" style="width: 900px; height: 500px"></div>
  </body>
</html>
`))

func makeWebDrawFile(tCtx webGraph) (*os.File, error) {
	// to avoid creating and cleaning temporary files on template errors
	var b bytes.Buffer
	err := graphTmpl.ExecuteTemplate(&b, "graphHtml", tCtx)
	if err != nil {
		return nil, fmt.Errorf("template failed: %v", err)
	}

	f, err := ioutil.TempFile("", "heavy_badger_web_report_*.html")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(f, &b)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (g *webGraph) Add(x insolar.Pulse, y float64) {
	g.Data = append(g.Data, tmplData{
		XValue:  pulseTime(x).String(),
		YValues: []float64{y},
	})
}

func (g *webGraph) Draw() {
	tmpFile, err := makeWebDrawFile(*g)
	if err != nil {
		fatalf("failed to make web draw: %v", err)
	}

	fmt.Println("saves report html file in", tmpFile.Name())
	cmd := exec.Command("open", "--wait-apps", tmpFile.Name())

	fin := &finalizersHolder{}
	fin.add(func() error {
		fmt.Println("\nremove", tmpFile.Name())
		return os.Remove(tmpFile.Name())
	})
	done := fin.onSignals(syscall.SIGINT, syscall.SIGTERM)

	err = cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fatalf("\nopen command failed: %v", err)
		}
		fmt.Println("\nopen command finished:", err)
	}
	<-done
}
