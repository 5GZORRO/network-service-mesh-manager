package main

import (
	"nextworks/nsm/internal/nbi"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	// log.SetReportCaller(true)  // add function name to logs

	log.Info("APP Started")

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// 	db, err := sql.Open("mysql", "root:root@/nsmm")
	dsn := "root:root@tcp(127.0.0.1:3306)/nsmm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	gc := nbi.GatewayConnectivity{SliceID: "prova1"}
	log.Info("Create ", gc)
	result := db.Create(&gc) // pass pointer of data to Create
	log.Info(result.Error)
	log.Info(result.RowsAffected)

	gc1 := nbi.GatewayConnectivity{SliceID: "prova2"}
	log.Info("Create ", gc1)
	result = db.Create(&gc1) // pass pointer of data to Create
	log.Info(result.Error)
	log.Info(result.RowsAffected)

	// Retrieve
	log.Info("Retrieve by primary key ", gc1)
	var gc2 nbi.GatewayConnectivity
	result = db.First(&gc2, "slice_id = ?", "prova2")
	log.Info(result.Error)
	log.Info(result.RowsAffected)
	log.Info("Retrieved ", gc2)

	// Update
	log.Info("Update ", gc)
	gc.NetworkID = "mod1"
	gc.SubnetID = "mod1"
	gc.RouterID = "mod1"
	gc.InterfaceID = "mod1"
	gc.FloatingIP = "mod1"
	gc.VmGatewayID = "mod1"
	result = db.Save(&gc).Debug()
	log.Info(result.Error)
	log.Info(result.RowsAffected)

	// Delete
	log.Info("Delete ", gc1)
	result = db.Delete(&gc1)
	log.Info(result.Error)
	log.Info(result.RowsAffected)

	// // Delete with additional conditions
	// db.Where("name = ?", "jinzhu").Delete(&email)
	// // DELETE from emails where id = 10 AND name = "jinzhu";

}
