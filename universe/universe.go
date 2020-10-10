package universe

import (
	"Eve-Dabblings/globals"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// response struct

//TODO: make this match the API response
type Region struct {
	Systems []struct {
		SystemID   int64
		SystemName string
	}
}

// returns body which is a []byte
func searchUniverse(data []byte) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, reqErr := http.NewRequest("POST", globals.EsiDomain+"/latest/universe/ids", bytes.NewReader([]byte("[\"kamio\"]"))) // TODO: figure out how to pass the values so the API likes it
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
		fmt.Println(resp.Header)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println("Error reading body")
		fmt.Println(readErr)
	}
	resp.Body.Close()

	return body, nil

}

func SearchSystem(regionName string) (*Region, error) {
	data, marshErr := func() ([]byte, error) {
		systemBytes, marshErr := json.Marshal(regionName)
		if marshErr != nil {
			return nil, marshErr
		}
		fmt.Println(string(systemBytes))
		return systemBytes, nil
	}()

	if marshErr != nil {
		return &Region{}, marshErr
	}

	results, searchErr := searchUniverse(data)
	if searchErr != nil {
		return &Region{}, searchErr
	}

	//TODO: this raises a panic- Region.Systems of type struct { SystemID int32; SystemName string }
	regionInfo, unmarshErr := func() (*Region, error) {
		var regionData *Region
		err := json.Unmarshal(results, &regionData)
		if err != nil {
			return nil, err
		}
		return regionData, nil
	}()

	if unmarshErr != nil {
		log.Error("Error unmarshalling response from server")
		return nil, unmarshErr
	}

	return regionInfo, nil
}
