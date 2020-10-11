package main

import (
	"Eve-Dabblings/regions"
	"Eve-Dabblings/universe"
	//	log "github.com/sirupsen/logrus"
	"fmt"
)

func getRegionIDs() {
	regionIDs, err := regions.GetRegionIDs()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(regionIDs)

}

func getUniverseID(name string) {
	data, err := universe.GetUniverseID(name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}

func main() {

	//globals.SetupLogging("debug")
	getUniverseID("The Citadel")

}
