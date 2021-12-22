package nsm

import (
	"nextworks/nsm/internal/vim"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// deleteResources is a goroutine to delete in an async way all the network resources
func deleteResources(database *gorm.DB, vim *vim.VimDriver, res *ResourceSet) {
	time.Sleep(time.Second * 5)
	// TODO Delete resources from VIM
	// (*vim).DeleteNetwork()

	// if removal from VIM is OK then delete it from DB
	result := database.Delete(&res)
	if result.Error != nil {
		log.Error("Error deleting network resource set with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}

// TODO configureGateway is a goroutine to configure the VM gateway, using
// an HTTP client
func configureGateway(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 10)
	// TODO configure VM gateway

	// update the state
	res.Status = READY
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}

// TODO resetGateway is a goroutine to configure the VM gateway, using
// an HTTP client
func resetGateway(database *gorm.DB, res *ResourceSet) {
	time.Sleep(time.Second * 10)
	// TODO reset VM gateway

	// update the state of the gateway to WAIT_FOR
	res.Status = WAIT_FOR_GATEWAY_CONFIG
	res.Gateway = Gateway{}
	result := database.Save(&res)
	if result.Error != nil {
		log.Error("Error updating resource set status with ID: ", res.ID, " and slice-id: ", res.SliceId)
	}
}
