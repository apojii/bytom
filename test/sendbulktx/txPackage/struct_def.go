package txPackage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bytom/crypto/ed25519/chainkd"
	"github.com/bytom/util"
)

type keyIns struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type Reveive struct {
	AccountID    string `json:"account_id"`
	AccountAlias string `json:"account_alias"`
}

type account struct {
	RootXPubs   []chainkd.XPub         `json:"root_xpubs"`
	Quorum      int                    `json:"quorum"`
	Alias       string                 `json:"alias"`
	Tags        map[string]interface{} `json:"tags"`
	AccessToken string                 `json:"access_token"`
}

type asset struct {
	RootXPubs   []chainkd.XPub         `json:"root_xpubs"`
	Quorum      int                    `json:"quorum"`
	Alias       string                 `json:"alias"`
	Tags        map[string]interface{} `json:"tags"`
	Definition  map[string]interface{} `json:"definition"`
	AccessToken string                 `json:"access_token"`
}

func RestoreStruct(data interface{}, out interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if ok != true {
		fmt.Println("invalid type assertion")
		os.Exit(util.ErrLocalParse)
	}

	rawData, err := json.MarshalIndent(dataMap, "", "  ")
	if err != nil {
		fmt.Println(err)
		os.Exit(util.ErrLocalParse)
	}
	json.Unmarshal(rawData, out)
}
