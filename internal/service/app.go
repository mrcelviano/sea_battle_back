package service

import "context"

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) InitBattlefield(ctx context.Context) ([][]int32, error) {
	var (
		length = 10
		width  = 10
	)

	battlefield := make([][]int32, 0, length)
	for i := range battlefield {
		widthBattlefield := make([]int32, 0, width)
		for j := range widthBattlefield {
			widthBattlefield[j] = -1
		}

		battlefield[i] = widthBattlefield
	}

	return battlefield, nil
}
