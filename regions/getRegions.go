package regions

import (
	"Eve-Dabblings/globals"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// returns body which is a []byte
func getRegions() ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, reqErr := http.NewRequest("GET", globals.EsiDomain+"/latest/universe/regions", nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Add("Content-Type", "application/json")
	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received %v", resp.Status)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println("Error reading body")
		fmt.Println(readErr)
	}
	resp.Body.Close()

	return body, nil

}

// returns the universe region IDs
func GetRegionIDs() ([]int32, error) {
	regionResp, callErr := getRegions()
	if callErr != nil {
		return nil, callErr
	}

	// simply unmarshal the response and return the []int32, error
	return func() ([]int32, error) {
		var regionIDs []int32 // response only contains a list of int32
		err := json.Unmarshal(regionResp, &regionIDs)
		if err != nil {
			return nil, err
		}
		return regionIDs, nil
	}()

}
