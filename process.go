package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/libangzhu/erc20Scaner/erc20abi/generated"
	"strings"
	"time"
)

// Process 业务处理模块
type Process struct {
	cli        *Client
	startPoint uint64
}

// Init 初始化模块
func (p *Process) Init() {
	if p.cli == nil {
		p.cli = new(Client)

	}
	p.cli.ConnectEth()
}

// Start 启动模块
func (p *Process) Start() {

	for {
		blockNum, err := p.cli.BlockNum()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if blockNum < p.startPoint {
			time.Sleep(time.Second)
			continue
		}

		block, err := p.cli.BlockByNumber(p.startPoint)
		if err != nil {
			fmt.Println("err:", err)
			continue
		}
		p.startPoint++
		err = p.ParaseBlock(block)
		if err != nil {
			fmt.Println("ParaseBlock err:", err)
			p.startPoint--
			continue
		}
	}

}

func (p *Process) Close() error {
	p.cli.CloseConnect()
	return nil
}

// ParaseBlock 解析区块
func (p *Process) ParaseBlock(block *types.Block) error {

	if block == nil {
		time.Sleep(time.Second)
		return fmt.Errorf("nil block")
	}

	txs := block.Transactions()

	fmt.Println("startPoint:", p.startPoint, "txsnum:", len(txs))
	for _, tx := range txs {

		if tx.To() == nil {
			//合约创建交易

			receipt, err := p.cli.TxReceipt(tx.Hash())
			if err != nil {
				return err
			}
			if receipt != nil {
				//log.Info("receipt:", receipt)
				if receipt.Status != types.ReceiptStatusSuccessful {
					fmt.Println("receipt status is ", receipt.Status)
					return err
				}
				//成功的交易
				//fmt.Println("evmaddress:", receipt.ContractAddress)
				//校验是否是ERC20的合约地址
				//jmb, err := json.Marshal(receipt)
				//if err != nil {
				//	return err
				//}
				//fmt.Println("receipt:", string(jmb))
				//check is or not ERC20 evm contract
				//decimals
				decimals, err := p.unPackageAbi("decimals", &receipt.ContractAddress)
				if err != nil {
					return err
				}
				cname, err := p.unPackageAbi("name", &receipt.ContractAddress)
				if err != nil {
					return err
				}
				symbol, err := p.unPackageAbi("symbol", &receipt.ContractAddress)
				if err != nil {
					return err
				}
				supply, err := p.unPackageAbi("totalSupply", &receipt.ContractAddress)
				if err != nil {
					return err
				}
				ut := time.Unix(int64(block.Time()), 0)
				cst, _ := time.LoadLocation("Asia/Shanghai")

				fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
				fmt.Println("contract address:", receipt.ContractAddress)
				fmt.Println("contract name:", cname)
				fmt.Println("contract symbol:", symbol)
				fmt.Println("contract totalSupply:", supply)
				fmt.Println("contract decimals:", decimals)
				fmt.Println("contract deploy time:", ut.In(cst))
				fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			}
		}
	}
	return nil
}

func (p *Process) unPackageAbi(methodName string, cAddress *common.Address) (interface{}, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(generated.ERC20ABI))
	if err != nil {
		panic(err)
	}
	abidata, err := parsedAbi.Pack(methodName)
	if err != nil {
		return nil, err
	}
	//通过节点查询返回evm查询结果
	var calldata ethereum.CallMsg
	calldata.Data = abidata
	calldata.To = cAddress
	callResult, err := p.cli.CallContract(calldata, nil)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = parsedAbi.UnpackIntoInterface(&result, methodName, callResult)
	if err != nil {
		return nil, err
	}

	//fmt.Println("result", result)
	return result, nil

}
