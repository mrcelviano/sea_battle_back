package main

import (
	"context"
	"fmt"
	"github.com/mrcelviano/sea_battle_back/internal/service"
	"log"
)

func main() {
	app := service.NewApp()
	ctx := context.Background()

	battlefield, err := app.InitBattlefield(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	battlefield, err = app.AutomaticPlacement(ctx, battlefield)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func printBattlefield(battlefield [][]int32) {
	for _, withValues := range battlefield {
		for _, withValue := range withValues {
			fmt.Print(withValue)
		}
		fmt.Println()
	}
}
