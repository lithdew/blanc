package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/lithdew/casso"
)

type View interface {
	Init(app *App) error
	Render(screen tcell.Screen)
}

type App struct {
	solver *casso.Solver

	sw casso.Symbol
	sh casso.Symbol

	views []View
}

func NewApp(screen tcell.Screen) (*App, error) {
	app := &App{
		solver: casso.NewSolver(),
		sw:     casso.New(),
		sh:     casso.New(),
	}

	if err := app.solver.Edit(app.sw, casso.Strong); err != nil {
		return nil, err
	}
	if err := app.solver.Edit(app.sh, casso.Strong); err != nil {
		return nil, err
	}
	sw, sh := screen.Size()
	if err := app.Resize(sw, sh); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) Add(view View) error {
	if err := view.Init(a); err != nil {
		return fmt.Errorf("failed to init view: %w", err)
	}
	a.views = append(a.views, view)
	return nil
}

func (a *App) Resize(width, height int) error {
	if err := a.solver.Suggest(a.sw, float64(width)); err != nil {
		return err
	}
	if err := a.solver.Suggest(a.sh, float64(height)); err != nil {
		return err
	}
	return nil
}

func (a *App) Render(screen tcell.Screen) {
	if a.Width() == 0 || a.Height() == 0 {
		return
	}

	for _, view := range a.views {
		view.Render(screen)
	}

	screen.Show()
}

func (a *App) Val(symbol casso.Symbol) float64 { return a.solver.Val(symbol) }

func (a *App) Width() int  { return int(a.Val(a.sw)) }
func (a *App) Height() int { return int(a.Val(a.sh)) }
