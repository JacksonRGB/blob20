package dao

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type blob20Inscription struct {
	Protocol string `json:"protocol"`
	Token    struct {
		Operation string `json:"operation"`
		Ticker    string `json:"ticker"`
		Amount    int    `json:"amount"`
	} `json:"token"`
}

func (d *Dao) CheckAttachment(txid string) (ok bool, err error) {
	resp, err := http.DefaultClient.Get(fmt.Sprintf("https://api-v2.ethscriptions.com/ethscriptions/%s/attachment", txid))
	if err != nil {
		log.WithError(err).Error("get attachment")
		return
	}
	defer resp.Body.Close()

	insc := &blob20Inscription{}

	err = json.NewDecoder(resp.Body).Decode(insc)
	if err != nil {
		log.WithError(err).Error("decode attachment")
		return
	}

	if insc.Protocol != "blob20" ||
		insc.Token.Operation != "mint" ||
		insc.Token.Amount != 1000 ||
		insc.Token.Ticker != "BLOB" {
		return
	}
	ok = true
	return

}
