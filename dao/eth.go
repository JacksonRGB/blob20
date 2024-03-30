package dao

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

func (d *Dao) GetBlockHeight(behindBlock ...int) (height int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	n, err := d.ethClient.BlockNumber(ctx)
	if len(behindBlock) > 0 {
		n -= uint64(behindBlock[0])
		if n < 0 {
			n = 0
		}
	}
	return int(n), err
}

func (d *Dao) GetLatestBockHash() (hash string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	block, err := d.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return
	}
	return block.Hash().Hex(), nil
}

func (d *Dao) GetBlockTime(height int) (timestamp int, err error) {
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		block, err := d.ethClient.BlockByNumber(ctx, big.NewInt(int64(height)))
		if err == nil {
			return int(block.Time()), nil
		}
	}
	return
}

type rpcBlock struct {
	Timestamp    hexutil.Big      `json:"timestamp"`
	Hash         common.Hash      `json:"hash"`
	Transactions []rpcTransaction `json:"transactions"`
}
type rpcTransaction struct {
	tx *types.Transaction
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetBatchedBlockTransactions(startHeight, endHeight int) (blockTimes []int, txss []types.Transactions, err error) {
	reqs := make([]rpc.BatchElem, 0)
	resps := make([]*rpcBlock, 0)
	for i := startHeight; i <= endHeight; i++ {
		res := &rpcBlock{}
		resps = append(resps, res)
		reqs = append(reqs, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeUint64(uint64(i)), true},
			Result: &res,
		})
	}
	for i := 0; i < 10; i++ {
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		err = d.ethRPC.BatchCallContext(ctx, reqs)
		if err != nil {
			log.WithError(err).Error("get batch block txs")
			continue
		}

		for j := 0; j < len(reqs); j++ {
			if reqs[j].Error != nil {
				log.WithError(err).Error("get batch block txs")
				continue
			}
		}
		break
	}

	if err != nil {
		return
	}

	for _, resp := range resps {
		blockTimes = append(blockTimes, int(resp.Timestamp.ToInt().Int64()))
		rawTxs := make([]*types.Transaction, 0)
		for _, tx := range resp.Transactions {
			rawTxs = append(rawTxs, tx.tx)
		}
		txss = append(txss, rawTxs)

	}
	return
}

func (d *Dao) GetBlockTransactions(height int) (blockTime int, txs types.Transactions, err error) {
	for i := 0; i < 10; i++ {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		block, err := d.ethClient.BlockByNumber(ctx, big.NewInt(int64(height)))
		if err != nil {
			log.WithError(err).Error("get block txs")
			continue
		}
		return int(block.Time()), block.Transactions(), nil
	}
	return
}

func (d *Dao) GetTransactionReceipt(txid string) (receipt *types.Receipt, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.ethClient.TransactionReceipt(ctx, common.HexToHash(txid))
}
