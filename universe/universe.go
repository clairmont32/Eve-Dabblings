package universe

import (
	"Eve-Dabblings/globals"
	"bytes"
	"encoding/json"
	"fmt"
	//log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// response struct for system and inventory items
type ResponseInfo struct {
	SystemInfo
	InvInfo
}

// stuct for systems, if systems exists
type SystemInfo struct {
	System []struct {
		SystemID int32  `json:"id"`
		Name     string `json:"name"`
	} `json:"systems,omitempty"`
}

// stuct for inventory items, if inventory_types exists
type InvInfo struct {
	Inventory []struct {
		TypeID int32  `json:"id"`
		Name   string `json:"name"`
	} `json:"inventory_types,omitempty"`
}

type RegionInfo struct {
	Region []struct {
		RegionID int32  `json:"id"`
		Name     string `json:"name"`
	} `json:"regions"`
}

// returns body which is a []byte
func callUniverseEndpoint(data []byte) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, reqErr := http.NewRequest("POST", globals.EsiDomain+"/latest/universe/ids", bytes.NewReader(data))
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Add("Content-Type", "application/json")
	resp, doErr := client.Do(req)
	if doErr != nil {
		return nil, doErr
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Received %v", resp.Status))
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

// submits a single system name and provides the ID
func GetUniverseID(searchName string) (*ResponseInfo, error) {
	data, marshErr := func() ([]byte, error) {
		systemBytes, marshErr := json.Marshal([]string{searchName})
		if marshErr != nil {
			return nil, marshErr
		}
		fmt.Println(string(systemBytes))
		return systemBytes, nil
	}()
	if marshErr != nil {
		return nil, marshErr
	}

	results, searchErr := callUniverseEndpoint(data)
	if searchErr != nil {
		return nil, searchErr
	}

	// unmarshal into SystemInfo
	regionInfo, unmarshErr := func() (*ResponseInfo, error) {
		var responseData *ResponseInfo
		err := json.Unmarshal(results, &responseData)
		if err != nil {
			fmt.Println(string(results))
			return nil, err
		}
		return responseData, nil
	}()

	if unmarshErr != nil {
		fmt.Println("Error unmarshalling response from server")
		return nil, unmarshErr
	}

	return regionInfo, nil
}
