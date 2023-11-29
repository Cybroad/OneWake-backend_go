package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// DBパス
const DBPATH = "./db/data.db"

// host構造体
type Host struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	HostName   string `json:"hostname"`
	IpAddress  string `json:"ipaddress"`
	MacAddress string `json:"macaddress"`
	OSName     string `json:"osname"`
}

type Response struct {
	Message string `json:"message"`
	Hosts   []Host `json:"hosts"`
}

func getRregisteredAllHosts(c echo.Context) error {
	// レコード数
	var cnt int

	db, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// レコード数を取得
	cntCmd := "SELECT COUNT(id) FROM hosts"
	cntRow, err := db.Query(cntCmd)
	if err != nil {
		panic(err.Error())
	}
	cntRow.Scan(&cnt)

	// レコードが0件の場合は空配列を返す
	if cnt == 0 {
		var res Response
		res.Message = "success"
		res.Hosts = []Host{}

		return c.JSON(http.StatusOK, res)
	} else {
		getCmd := "SELECT id, name, hostName, ipAddress, macAddress, osName FROM hosts"
		getRow, err := db.Query(getCmd)
		if err != nil {
			panic(err.Error())
		}

		var hosts []Host
		for getRow.Next() {
			var host Host
			if err := getRow.Scan(&host.Id, &host.Name, &host.HostName); err != nil {
				panic(err.Error())
			}
			hosts = append(hosts, host)
		}

		fmt.Println(hosts)
		return c.JSON(http.StatusOK, hosts)
	}
}

func addHost(c echo.Context) error {
	return c.String(http.StatusOK, "add host")
}

func updateHostInfo(c echo.Context) error {
	return c.String(http.StatusOK, "update host info")
}

func deleteHost(c echo.Context) error {
	return c.String(http.StatusOK, "delete host")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// DBの初期化
	db, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		panic(err.Error())
	}

	cmd := "CREATE TABLE IF NOT EXISTS hosts (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, hostName TEXT, ipAddress TEXT, macAddress TEXT, osName TEXT)"
	_, err = db.Exec(cmd)
	if err != nil {
		panic(err.Error())
	}

	db.Close()

	// 動作確認用
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "it works!")
	})

	// ルート定義
	g := e.Group("/api/v1") // ルートグループ
	g.GET("/hosts", getRregisteredAllHosts)
	g.POST("/host/add", addHost)
	g.PUT("/host/update", updateHostInfo)
	g.DELETE("/host/delete", deleteHost)

	e.Logger.Fatal(e.Start(":8080"))
}
