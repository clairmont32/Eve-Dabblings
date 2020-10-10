package main

import (
	"Eve-Dabblings/regions"
	"fmt"
)

func main() {

	regionIDs, err := regions.GetRegionIDs()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(regionIDs)

}
