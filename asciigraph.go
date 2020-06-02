package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/asciigraph"
	"github.com/lithdew/blanc/layout"
)

type ASCIIGraph struct {
	lastW  int // last rendered width
	lastH  int // last rendered height
	series []float64
	text   Text
}

func NewASCIIGraph() ASCIIGraph {
	g := ASCIIGraph{series: []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}}
	//g.text.SetWrap(true)
	return g
}

func (g *ASCIIGraph) SetStyle(style tcell.Style) { g.text.SetStyle(style) }
func (g *ASCIIGraph) SetSeries(series []float64) { g.series = series }

func (g *ASCIIGraph) Draw(s tcell.Screen, r layout.Rect) {
	if r.W != g.lastW || r.H != g.lastH {
		g.text.SetText(asciigraph.Plot(g.series, asciigraph.Width(r.W), asciigraph.Height(r.H)))
		g.lastW = r.W
		g.lastH = r.H
	}

	g.text.Draw(s, r)
}
