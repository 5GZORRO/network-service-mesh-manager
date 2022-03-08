package nsm

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// LoadAssociationFromDB loads from DB all the active connection of a resource set
func LoadConnectionAssociationFromDB(database *gorm.DB, netres *ResourceSet) error {
	var conns []Connection

	// Retrieve connections association
	a := database.Model(&netres).Association("Connections")
	result := a.Find(&conns)
	if result != nil {
		return result
	}
	netres.Connections = conns
	return nil
}

// SearchConnectionAssociationFromDB searches a Connection by ConnectionID of a resource set
func SearchConnectionAssociationFromDB(database *gorm.DB, netres *ResourceSet, cid int) (*Connection, error) {
	var conn []Connection

	// Retrieve connections association
	result := database.Model(&netres).Where("id = " + fmt.Sprint(cid)).Association("Connections").Find(&conn)
	if result != nil {
		log.Error("Error retrieving connections for resource set ", netres.ID)
		return nil, result
	}
	if len(conn) == 0 {
		log.Error("Error no connections found for resource set ", netres.ID)
		return nil, ErrNoConnection
	}
	return &conn[0], nil
}
