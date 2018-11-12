package robot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

// Mc match configuration

var wg sync.Mutex

var mc *MatchConfig

// GetMatchConfig get match config
func GetMatchConfig() *MatchConfig {
	wg.Lock()
	defer wg.Unlock()
	if mc == nil {
		configFile := "application.json"
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatalf("can't read config:%v", err)
		}
		json.Unmarshal(file, &mc)
	}
	return mc
}
