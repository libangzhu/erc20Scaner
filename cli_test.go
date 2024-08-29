package main

import (
	"github.com/ethereum/go-ethereum/common"

	"testing"
)

var (
	cli = new(Client)
	_   = cli.ConnectEth()
)

func TestClient_BlockNum(t *testing.T) {

	num, err := cli.BlockNum()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("current block number:", num)
}

func TestClient_BlockByNumber(t *testing.T) {
	num, err := cli.BlockNum()
	if err != nil {
		t.Error(err)
		return
	}
	num = 34562329
	block, err := cli.BlockByNumber(num)
	if err != nil {
		t.Error(err)
		return
	}
	txs := block.Transactions()
	t.Log("tx.num:", len(txs))
	for _, tx := range txs {

		if tx.To() == nil {
			//合约创建交易
			t.Log("tx.Hash:", tx.Hash(), "tx.to:", tx.To(), "tx.value:", tx.Value(), "tx.data", common.Bytes2Hex(tx.Data()))
		}

	}
}

func Test_ParaseBlock(t *testing.T) {
	p := new(Process)
	p.Init()
	num := 34562329
	block, err := p.cli.BlockByNumber(uint64(num))
	if err != nil {
		t.Error(err)
		return
	}
	err = p.ParaseBlock(block)
	if err != nil {
		t.Error(err)
		return
	}
}
