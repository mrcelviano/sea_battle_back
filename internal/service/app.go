package service

import (
	"context"
	"crypto/rand"
	"github.com/mrcelviano/sea_battle_back/internal/model"
	"math/big"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) InitBattlefield(ctx context.Context) ([][]int32, error) {
	battlefield := make([][]int32, model.LengthBattlefield)
	for i := range battlefield {
		widthBattlefield := make([]int32, model.WidthBattlefield)
		for j := range widthBattlefield {
			widthBattlefield[j] = model.EmptyCell
		}

		battlefield[i] = widthBattlefield
	}

	return battlefield, nil
}

func (a *App) AutomaticPlacement(ctx context.Context, battlefield [][]int32) ([][]int32, error) {
	for _, ship := range model.CharacteristicsShips {
		for i := int32(0); i < ship.Count; i++ {
			battlefield, err := a.SetShipToBattlefield(ctx, ship.Length, battlefield)
			if err != nil {
				return battlefield, err
			}
		}
	}

	return battlefield, nil
}

func (a *App) SetShipToBattlefield(ctx context.Context, lengthShip int32, battlefield [][]int32) ([][]int32, error) {
	position, err := getPositionType(ctx, lengthShip)
	if err != nil {
		return nil, err
	}

	switch position {
	case model.VerticalPosition:
		battlefield, err = setShipFromVerticalPosition(ctx, lengthShip, battlefield)
		if err != nil {
			return nil, err
		}
	case model.HorizontalPosition:
		battlefield, err = setShipFromHorizontalPosition(ctx, lengthShip, battlefield)
		if err != nil {
			return nil, err
		}
	case model.NotSetPosition:
		battlefield, err = setShipFromNotSetPosition(ctx, battlefield)
		if err != nil {
			return nil, err
		}
	}

	return battlefield, nil
}

func setShipFromNotSetPosition(ctx context.Context, battlefield [][]int32) ([][]int32, error) {
	for {
		xC, yC, err := getCoordinate(ctx)
		if err != nil {
			return nil, err
		}

		isSetShip := checkCoordinate(ctx, xC, yC, battlefield)
		if !isSetShip {
			continue
		}

		battlefield[yC][xC] = model.ShipCell

		if xC-1 >= 0 {
			battlefield[yC][xC-1] = model.LockedCell
		}
		if xC+1 <= model.WidthBattlefield-1 {
			battlefield[yC][xC+1] = model.LockedCell
		}

		if yC-1 >= 0 {
			battlefield[yC-1][xC] = model.LockedCell
			if xC-1 >= 0 {
				battlefield[yC-1][xC-1] = model.LockedCell
			}
			if xC+1 <= model.WidthBattlefield-1 {
				battlefield[yC-1][xC+1] = model.LockedCell
			}
		}

		if yC+1 <= model.LengthBattlefield-1 {
			battlefield[yC+1][xC] = model.LockedCell
			if xC-1 >= 0 {
				battlefield[yC+1][xC-1] = model.LockedCell
			}
			if xC+1 <= model.WidthBattlefield-1 {
				battlefield[yC+1][xC+1] = model.LockedCell
			}
		}

		return battlefield, nil
	}
}

func setShipFromHorizontalPosition(ctx context.Context, lengthShip int32, battlefield [][]int32) ([][]int32, error) {
	for {
		xC, yC, err := getCoordinate(ctx)
		if err != nil {
			return nil, err
		}

		var (
			coordinatesShipCell       = make([]model.CoordinateCell, 0, lengthShip)
			coordinatesLockedCellCell = make([]model.CoordinateCell, 0)
			startCoordinate           = xC
			finishCoordinate          int32
		)

		for i := xC; i < lengthShip+xC; i++ {
			if i == startCoordinate {
				finishCoordinate = i
				continue
			}

			if i > model.WidthBattlefield-1 {
				startCoordinate -= 1
				continue
			}

			finishCoordinate = i
		}

		isNotSetShip := false
		for i := startCoordinate; i <= finishCoordinate; i++ {
			isSetShip := checkCoordinate(ctx, i, yC, battlefield)
			if !isSetShip {
				isNotSetShip = true
			}

			coordinatesShipCell = append(coordinatesShipCell, model.CoordinateCell{
				X: i,
				Y: yC,
			})

			if i == startCoordinate {
				if i-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i - 1,
						Y: yC,
					})
				}

				if yC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC - 1,
					})
					if i-1 >= 0 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: i - 1,
							Y: yC - 1,
						})
					}
				}

				if yC+1 <= model.LengthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC + 1,
					})
					if i-1 >= 0 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: i - 1,
							Y: yC + 1,
						})
					}
				}
			}

			if i > startCoordinate && i < finishCoordinate {
				if yC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC - 1,
					})
				}

				if yC+1 <= model.LengthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC + 1,
					})
				}
			}

			if i == finishCoordinate {
				if i+1 <= model.WidthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i + 1,
						Y: yC,
					})
				}

				if yC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC - 1,
					})
					if i+1 <= model.WidthBattlefield-1 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: i + 1,
							Y: yC - 1,
						})
					}
				}

				if yC+1 <= model.LengthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: i,
						Y: yC + 1,
					})
					if i+1 <= model.WidthBattlefield-1 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: i + 1,
							Y: yC + 1,
						})
					}
				}
			}
		}

		if isNotSetShip {
			continue
		}

		for _, coordinateShip := range coordinatesShipCell {
			battlefield[coordinateShip.Y][coordinateShip.X] = model.ShipCell
		}

		for _, coordinateLockedCellCell := range coordinatesLockedCellCell {
			battlefield[coordinateLockedCellCell.Y][coordinateLockedCellCell.X] = model.LockedCell
		}

		return battlefield, nil
	}
}

func setShipFromVerticalPosition(ctx context.Context, lengthShip int32, battlefield [][]int32) ([][]int32, error) {
	for {
		xC, yC, err := getCoordinate(ctx)
		if err != nil {
			return nil, err
		}

		var (
			coordinatesShipCell       = make([]model.CoordinateCell, 0, lengthShip)
			coordinatesLockedCellCell = make([]model.CoordinateCell, 0)
			startCoordinate           = yC
			finishCoordinate          int32
		)

		for i := yC; i < lengthShip+yC; i++ {
			if i == startCoordinate {
				finishCoordinate = i
				continue
			}

			if i > model.LengthBattlefield-1 {
				startCoordinate -= 1
				continue
			}

			finishCoordinate = i
		}

		isNotSetShip := false
		for i := startCoordinate; i <= finishCoordinate; i++ {
			isSetShip := checkCoordinate(ctx, xC, i, battlefield)
			if !isSetShip {
				isNotSetShip = true
			}

			coordinatesShipCell = append(coordinatesShipCell, model.CoordinateCell{
				X: xC,
				Y: i,
			})

			if i == startCoordinate {
				if i-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC,
						Y: i - 1,
					})
				}

				if xC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC - 1,
						Y: i,
					})
					if i-1 >= 0 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: xC - 1,
							Y: i - 1,
						})
					}
				}

				if xC+1 <= model.WidthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC + 1,
						Y: i,
					})
					if i-1 >= 0 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: xC + 1,
							Y: i - 1,
						})
					}
				}
			}

			if i > startCoordinate && i < finishCoordinate {
				if xC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC - 1,
						Y: i,
					})
				}

				if xC+1 <= model.WidthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC + 1,
						Y: i,
					})
				}
			}

			if i == finishCoordinate {
				if i+1 <= model.LengthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC,
						Y: i + 1,
					})
				}

				if xC-1 >= 0 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC - 1,
						Y: i,
					})
					if i+1 <= model.LengthBattlefield-1 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: xC - 1,
							Y: i + 1,
						})
					}
				}

				if xC+1 <= model.WidthBattlefield-1 {
					coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
						X: xC + 1,
						Y: i,
					})
					if i+1 <= model.LengthBattlefield-1 {
						coordinatesLockedCellCell = append(coordinatesLockedCellCell, model.CoordinateCell{
							X: xC + 1,
							Y: i + 1,
						})
					}
				}
			}
		}

		if isNotSetShip {
			continue
		}

		for _, coordinateShip := range coordinatesShipCell {
			battlefield[coordinateShip.Y][coordinateShip.X] = model.ShipCell
		}

		for _, coordinateLockedCellCell := range coordinatesLockedCellCell {
			battlefield[coordinateLockedCellCell.Y][coordinateLockedCellCell.X] = model.LockedCell
		}

		return battlefield, nil
	}
}

func getPositionType(ctx context.Context, lengthShip int32) (model.PositionType, error) {
	if lengthShip == 1 {
		return model.NotSetPosition, nil
	}

	r, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return 0, err
	}

	if r.Int64() <= 50 {
		return model.HorizontalPosition, nil
	}

	return model.VerticalPosition, nil
}

func getCoordinate(ctx context.Context) (x, y int32, err error) {
	xC, err := rand.Int(rand.Reader, big.NewInt(model.WidthBattlefield-1))
	if err != nil {
		return 0, 0, err
	}

	yC, err := rand.Int(rand.Reader, big.NewInt(model.LengthBattlefield-1))
	if err != nil {
		return 0, 0, err
	}

	return int32(xC.Int64()), int32(yC.Int64()), nil
}

func checkCoordinate(ctx context.Context, xC, yC int32, battlefield [][]int32) bool {
	if battlefield[yC][xC] == model.ShipCell || battlefield[yC][xC] == model.LockedCell {
		return false
	}

	if xC-1 >= 0 && battlefield[yC][xC-1] == model.ShipCell {
		return false
	}

	if xC+1 <= model.WidthBattlefield-1 && battlefield[yC][xC+1] == model.ShipCell {
		return false
	}

	if yC-1 >= 0 && battlefield[yC-1][xC] != model.ShipCell {
		if xC-1 >= 0 && battlefield[yC-1][xC-1] == model.ShipCell {
			return false
		}

		if xC+1 <= model.WidthBattlefield-1 && battlefield[yC-1][xC+1] == model.ShipCell {
			return false
		}
	}

	if yC+1 <= model.LengthBattlefield-1 && battlefield[yC+1][xC] != model.ShipCell {
		if xC-1 >= 0 && battlefield[yC+1][xC-1] == model.ShipCell {
			return false
		}

		if xC+1 <= model.WidthBattlefield-1 && battlefield[yC+1][xC+1] == model.ShipCell {
			return false
		}
	}

	return true
}
