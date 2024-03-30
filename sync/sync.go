package sync

import (
	"blob-index/config"
	"blob-index/dao"
	dbmodel "blob-index/model/db"
	"bytes"
	"flag"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

var (
	LastBlockKey = "last_block"
	BehindBlock  = 1
	PrefixData   = common.Hex2Bytes("646174613a3b72756c653d65736970362c")
	BatchSize    = 4
	jobCh        = make(chan bool, 5)
)

type Sync struct {
	c *config.Config
	d *dao.Dao
	sync.Mutex
}

var manualSync = flag.Int("sync", 0, "sync block height")

func New(_c *config.Config, _d *dao.Dao) (sync *Sync) {
	sync = &Sync{
		c: _c,
		d: _d,
	}
	return sync
}

func (s *Sync) Start() {
	lastHeight, err := s.d.GetStorageHeight(LastBlockKey)
	if err != nil {
		log.WithError(err).Error("get last block height")
		return
	}

	if lastHeight != 1 {
		// 数据库里保存的是已完成的区块, 再次同步时+1
		lastHeight++
	}

	if *manualSync > 0 {
		lastHeight = *manualSync
	}

	log.WithField("height", lastHeight).Info("last sync block height")

	var latestHeight int
	var beginHeight = lastHeight
	var endHeight = beginHeight + BatchSize

	for {
		latestHeight, err = s.d.GetBlockHeight(BehindBlock)
		if err != nil {
			log.WithError(err).Error("get latest block height")
			return
		}
		if (latestHeight-BatchSize)-beginHeight < BatchSize+1 {
			time.Sleep(10 * time.Second)
			continue
		}

		s.syncBatchBlock(beginHeight, endHeight)
		// for i := beginHeight; i <= endHeight; i++ {
		// 	s.syncBlock(i)
		// }

		if err = s.d.SetStorageHeight(LastBlockKey, endHeight); err != nil {
			log.WithError(err).Error("set last block height")
		}
		beginHeight = endHeight + 1
		endHeight = beginHeight + BatchSize
		log.WithFields(log.Fields{
			"begin height":  beginHeight,
			"end height":    endHeight,
			"latest height": latestHeight,
		}).Info("sync block height")
	}

}

func (s *Sync) syncBatchBlock(start, end int) {
	blockTimes, txss, err := s.d.GetBatchedBlockTransactions(start, end)
	if err != nil {
		log.WithError(err).Error("get batch block transactions")
		return
	}

	for i := range txss {
		for j := range txss[i] {
			jobCh <- true
			go s.processTx(start+i, blockTimes[i], txss[i][j])
		}
	}
}

func (s *Sync) processTx(blockNumber, blockTime int, tx *types.Transaction) {
	defer func() {
		<-jobCh
	}()
	if s.processMint(blockNumber, blockTime, tx) {
		return
	}
}

func (s *Sync) syncBlock(blockNumber int) {
	blockTime, txs, err := s.d.GetBlockTransactions(blockNumber)
	if err != nil {
		log.WithError(err).Error("get block transactions")
		return
	}
	for _, tx := range txs {
		if s.processMint(blockNumber, blockTime, tx) {
			continue
		}
	}
	return
}

func (s *Sync) processMint(block, blockTime int, tx *types.Transaction) (isMint bool) {
	if !bytes.HasPrefix(tx.Data(), PrefixData) {
		// esip6
		return
	}

	if len(tx.BlobHashes()) != 1 {
		// 必须只有一笔
		return
	}

	targetHash := common.HexToHash("01ee8325bc5607a16dd64ff2bcbec7d596b170f31def52615abf6b3f25ceb5a5")
	if tx.BlobHashes()[0].String() != targetHash.String() {
		return
	}

	receipt, err := s.d.GetTransactionReceipt(tx.Hash().String())
	if err != nil {
		log.WithField("txid", tx.Hash().String()).WithError(err).Error("get transaction receipt")
		return
	}

	txGas := uint(receipt.GasUsed) * uint(receipt.EffectiveGasPrice.Int64())
	blobGas := uint(receipt.BlobGasUsed) * uint(receipt.BlobGasPrice.Int64())

	mint := &dbmodel.Mint{
		TxHash:       tx.Hash().String(),
		Owner:        strings.ToLower(tx.To().String()),
		Block:        block,
		BlockTime:    blockTime,
		BlobGas:      blobGas,
		BlobGasPrice: uint(receipt.BlobGasPrice.Int64()),
		TxGas:        txGas,
		TotalGas:     txGas + blobGas,
	}

	_, err = s.d.CreateMint(mint)
	if err != nil {
		log.WithError(err).Error("create mint")
		return
	}
	log.WithField("owner", mint.Owner).WithField("txid", tx.Hash().String()).Info("create mint")

	return true
}
