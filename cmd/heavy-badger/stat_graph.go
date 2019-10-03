//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/campoy/tools/imgcat"
	"github.com/wcharczuk/go-chart"

	"github.com/insolar/insolar/insolar"
)

type Grapher interface {
	Add(insolar.Pulse, float64)
	Draw()
}

type StubDrawer struct{}

func (d StubDrawer) Add(insolar.Pulse, float64) {}

func (d StubDrawer) Draw() {}

type ConsoleGraph struct {
	x []time.Time
	y []float64
}

func (g *ConsoleGraph) Add(x insolar.Pulse, y float64) {
	g.x = append(g.x, pulseTime(x))
	g.y = append(g.y, y)
}

func (g *ConsoleGraph) Draw() {
	enc, err := imgcat.NewEncoder(
		os.Stdout,
		imgcat.Inline(true),
		imgcat.Width(imgcat.Percent(100)),
	)
	if err != nil {
		fatalf("%s\n", err)
	}

	fmt.Print("Press 'Enter' to show size graph...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')

	b := g.genImage()
	if err := enc.Encode(bytes.NewReader(b)); err != nil {
		fatalf("image encode failed: %v", err)
	}
}

func (g *ConsoleGraph) genImage() []byte {
	mainSeries := chart.TimeSeries{
		Name: "DB Size Rate (Mb)",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorBlue,
			FillColor:   chart.ColorBlue.WithAlpha(100),
		},
		XValues: g.x,
		YValues: g.y,
	}
	width := 1280
	if len(g.y) > width {
		width = len(g.y)
	}
	graph := chart.Chart{
		// Log:    chart.NewLogger(),
		Width:  width,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},
		YAxis: chart.YAxis{
			Name:  "Values size",
			Style: chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
			// ValueFormatter: func(v interface{}) string {
			// 	return fmt.Sprintf("%d bytes", int(v.(float64)))
			// },
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeHourValueFormatter,
			GridMajorStyle: chart.Style{
				// Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
		},
		Series: []chart.Series{
			mainSeries,
		},
	}

	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}

	var b bytes.Buffer
	err := graph.Render(chart.PNG, &b)
	if err != nil {
		fatalf("render failed: %v", err)
	}
	return b.Bytes()
}
