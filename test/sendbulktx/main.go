package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bytom/blockchain/pseudohsm"
	"github.com/bytom/blockchain/txbuilder"
	pack "github.com/bytom/test/sendbulktx/txPackage"
)

// 1、build-transaction sign-transaction submit-transaction 发送交易，回去到返回值tx_id
// 2、根据tx_id 查询交易是否存在 get-transaction
// 3、发送的交易总数，成功的总数，失败的总数以及那些tx_id失败

func main() {
	acctNum := flag.Int("acctNum", 10, "Number of created accounts")
	btmNum := flag.Int("btmNum", 10000, "Number of btm to send trading accounts")
	thdNum := flag.Int("thdNum", 5, "goroutine num")
	txBtmNum := flag.Int("txBtmNum", 10, "Number of transactions btm")
	sendAcct := flag.String("sendAcct", "0CHHJNM3G0A02", "who send btm")
	flag.Parse()

	controlPrograms := make([]string, *acctNum)
	txidChan := make(chan string)
	// create key
	param := []string{"alice", "123"}
	fmt.Println("*****************create key start*****************")
	var xpub pseudohsm.XPub
	resp, b := pack.SendReq(pack.CreateKey, param)
	if !b {
		resp, b := pack.SendReq(pack.ListKeys, param)
		if !b {
			os.Exit(1)
		}
		dataList, _ := resp.([]interface{})
		for _, item := range dataList {
			pack.RestoreStruct(item, &xpub)
			if strings.EqualFold(xpub.Alias, param[0]) {
				break
			}

		}
	} else {
		pack.RestoreStruct(resp, &xpub)
	}
	fmt.Println("*****************create key end*****************")

	fmt.Println("*****************create account start*****************")
	for i := 0; i < *acctNum; i++ {
		// create account
		name := fmt.Sprintf("alice%d", i)
		param = []string{name, xpub.XPub.String()}
		_, b = pack.SendReq(pack.CreateAccount, param)
		// create receiver
		param = []string{name}
		resp, b = pack.SendReq(pack.CreateReceiver, param)
		if !b {
			os.Exit(1)
		}
		var recv txbuilder.Receiver
		pack.RestoreStruct(resp, &recv)
		recvText, _ := recv.ControlProgram.MarshalText()
		controlPrograms[i] = string(recvText)
	}
	fmt.Println("*****************create account end*****************")

	threadTxNum := *btmNum / (*thdNum * *txBtmNum)
	txBtm := fmt.Sprintf("%d", *txBtmNum*2000)
	fmt.Println("*****************send tx start*****************")
	// send btm to account
	for i := 0; i < *thdNum; i++ {
		go pack.Sendbulktx(threadTxNum, txBtm, *sendAcct, controlPrograms, txidChan)
	}
	var txid string
	fail := 0
	sucess := 0
	//以读写方式打开文件，如果不存在，则创建
	file, error := os.OpenFile("./txid.txt", os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		fmt.Println(error)
	}
	for {
		select {
		case txid = <-txidChan:
			if strings.EqualFold(txid, "") {
				fail++
			} else {
				sucess++
				file.WriteString(txid)
				file.WriteString("\n")
			}
		default:
			if (sucess + fail) >= (*thdNum * threadTxNum) {
				file.Close()
				os.Exit(0)
			}
			time.Sleep(time.Second * 2)
		}
	}
}
