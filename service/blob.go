package service

import (
	"blob-index/model/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (srv *Service) GetMints(owner string, page, pageSize int) (resp *http.GetMintsResponse, err error) {
	resp = &http.GetMintsResponse{}

	if owner == "" {
		resp.TotalCount = globalTotalCount
	} else {
		totalCount, err := srv.d.GetMintCountByOwner(owner)
		if err != nil {
			log.WithError(err).Error("db: failed to get inscription count by owner")
			return resp, err
		}
		resp.TotalCount = totalCount
	}
	mints, err := srv.d.GetMintsByOwner(owner, page, pageSize)
	if err != nil {
		log.WithError(err).Error("db: failed to get inscriptions by owner")
		return
	}

	for _, mint := range mints {
		resp.Mints = append(resp.Mints, http.Mint{
			TxHash:       mint.TxHash,
			BlockNumber:  mint.Block,
			Owner:        mint.Owner,
			Amount:       1000,
			BlobGasPrice: mint.BlobGasPrice,
			GasFee:       mint.TotalGas,
			Time:         mint.BlockTime,
		})
	}
	return
}

func (s *Service) GetTotalCount() int {
	return globalTotalCount
}

var globalTotalCount = 0

func (s *Service) loopGetTotalCount() {
	for {
		globalTotalCount, _ = s.d.GetMintCount()
		time.Sleep(time.Second * 10)
	}
}
