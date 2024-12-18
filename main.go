package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	start := os.Args[1]
	iStart, err := strconv.Atoi(start)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// 连接到以太坊节点
	//client, err := ethclient.Dial("https://rpc.ankr.com/eth")
	client, err := ethclient.Dial("http://119.147.213.60:38545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	// 连接到SQLite数据库
	db, err := sql.Open("sqlite3", "ens.db") // 替换为你的数据库文件路径
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// count
	// 要查询的表名
	tableName := "reg_logs" // 替换为你的表名

	// 执行查询
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	// 打印记录总数
	fmt.Printf("Total number of records in %s: %d\n", tableName, count)

	// modify balance for all records
	for i := iStart; i <= count; i++ {
		// get owner
		query := "SELECT owner FROM reg_logs WHERE lid = ?" // 替换为你的表名和字段名
		var owner string
		err = db.QueryRow(query, i).Scan(&owner) // 使用lid的值作为查询参数
		if err != nil {
			log.Fatal(err)
		}

		// get balance
		bal := getBalance(client, common.HexToAddress(owner))
		fmt.Printf("lid: %d, owner: %s,balance: %s\n", i, owner, bal)

		// update balance
		// 准备SQL语句
		updateStmt := `UPDATE reg_logs SET balance = ? WHERE lid= ?` // 替换为你的表名和字段名
		balanceValue := bal                                          // 替换为你要设置的值
		lid := i
		// 执行SQL语句
		result, err := db.Exec(updateStmt, balanceValue, lid)
		if err != nil {
			log.Fatal(err)
		}
		// 获取影响的行数
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Rows affected: %d\n", rowsAffected)

		fmt.Println()
	}
}

// get balance
func getBalance(client *ethclient.Client, address common.Address) string {

	// 指定要查询的账户地址
	//address := common.HexToAddress("0x29bE0ceE87c6b18C2e264313A9EaD66c083B894F") // 替换为你要查询的地址

	// 查询余额
	balance, err := client.BalanceAt(context.Background(), address, nil) // nil 表示最新区块
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	// 打印余额
	//fmt.Printf("Balance of address %s: %s ETH\n", address.Hex(), balance.String())

	return balance.String()
}
