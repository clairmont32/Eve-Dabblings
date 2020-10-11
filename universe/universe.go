package universe

import (
	"Eve-Dabblings/globals"
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// response struct for system information
type SystemInfo struct {
	Systems []struct {
		SystemID   int32  `json:"id"`
		SystemName string `json:"name"`
	} `json:"systems"`
}

type InvInfo struct {
	Type []struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"inventory_types"`
}

// returns body which is a []byte
func searchUniverse(data []byte) ([]byte, error) {
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
		log.Printf("Received %v", resp.Status)
		log.Infoln(resp.Header)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Infoln("Error reading body")
		log.Infoln(readErr)
	}
	resp.Body.Close()

	return body, nil

}

// submits a single system name and provides the ID
func SearchSystem(systemName string) (*SystemInfo, error) {
	data, marshErr := func() ([]byte, error) {
		systemBytes, marshErr := json.Marshal([]string{systemName})
		if marshErr != nil {
			return nil, marshErr
		}
		log.Infoln(string(systemBytes))
		return systemBytes, nil
	}()

	if marshErr != nil {
		return &SystemInfo{}, marshErr
	}

	results, searchErr := searchUniverse(data)
	if searchErr != nil {
		return &SystemInfo{}, searchErr
	}

	// unmarshal into SystemInfo
	regionInfo, unmarshErr := func() (*SystemInfo, error) {
		var regionData *SystemInfo
		err := json.Unmarshal(results, &regionData)
		if err != nil {
			log.Infoln(string(results))
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
