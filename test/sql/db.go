package main

import (
	"database/sql"
	"nextworks/nsm/internal/nbi"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	// log.SetReportCaller(true)  // add function name to logs

	log.Info("APP Started")

	db, err := sql.Open("mysql", "root:root@/nsmm")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Use the DB normally, execute the querys etc
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// INSERT DATA
	insert, err := db.Prepare("INSERT INTO gateways VALUES(?,'','','','','','') ") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	result, err := insert.Exec("prova3")
	if err != nil {
		log.Error(err)
	}
	log.Info(result)
	result, err = insert.Exec("prova3")
	if err != nil {
		log.Error(err)
	}
	log.Info(result)
	defer insert.Close() // Close the statement when we leave main() / the program terminates

	// RETRIEVE
	// Prepare statement for reading data
	var e nbi.GatewayConnectivity
	selectX, err := db.Prepare("SELECT * FROM gateways WHERE sliceId = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer selectX.Close()
	boolean := selectX.QueryRow("prova1").Scan(&e.SliceID, &e.NetworkID, &e.SubnetID, &e.RouterID, &e.InterfaceID, &e.FloatingIP, &e.VmGatewayID)
	log.Info(boolean)
	log.Info(e)

	// DELETE DATA
	delete, err := db.Prepare("DELETE FROM gateways WHERE sliceId = ? ") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer delete.Close()
	_, err = delete.Exec("prova3")
	if err != nil {
		log.Error(err)
	}
}
