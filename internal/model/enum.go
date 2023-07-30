package model

type Ship string

const (
	BattleshipShip = Ship("Battleship")
	CruisersShip   = Ship("Cruisers")
	DestroyersShip = Ship("Destroyers")
	TorpedoShip    = Ship("Torpedo")
)

var CharacteristicsShips = map[Ship]CharacteristicShip{
	BattleshipShip: {Length: 4, Count: 1},
	CruisersShip:   {Length: 3, Count: 2},
	DestroyersShip: {Length: 2, Count: 3},
	TorpedoShip:    {Length: 1, Count: 4},
}

const (
	EmptyCell  = 0
	LockedCell = 1
	ShipCell   = 2
)

type PositionType int

var (
	NotSetPosition     PositionType = 0
	VerticalPosition   PositionType = 1
	HorizontalPosition PositionType = 2
)

const (
	LengthBattlefield = 10
	WidthBattlefield  = 10
)
