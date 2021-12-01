package vim

import (
	"nextworks/nsm/internal/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VimDriverList struct {
	VimList map[string]VimDriver
}

func (vims *VimDriverList) Exists(name string) bool {
	return vims.VimList[name] != nil
}

func (vims *VimDriverList) GetVim(name string) *VimDriver {
	vim := vims.VimList[name]
	return &vim
}

func (vims *VimDriverList) addVim(name string, vim VimDriver) {
	if vims.VimList[name] == nil {
		vims.VimList[name] = vim
	} else {
		log.Error("Vim ", name, " already exists")
	}
}

func newVimDriverList() *VimDriverList {
	return &VimDriverList{VimList: make(map[string]VimDriver)}
}

// TODO to be complete - also reading from DB
// and save new vim in DB
func InizializeVims(db *gorm.DB, vimConfigs []config.VimConfigurations) *VimDriverList {
	log.Info("Initializing vims...")
	vimList := newVimDriverList()

	// TODO first read from DB
	// then read config file
	for _, configVim := range vimConfigs {
		log.Debug("Vim config ", configVim)
		err := config.CheckVimParams(configVim)
		if err != nil {
			log.Error(err)
		} else {
			// log.Info("Type: ", configVim.Type)
			switch configVim.Type {
			case string(Openstack):
				openstackclient := NewOpenStackDriver(configVim.IdentityEndpoint, configVim.Username, configVim.Password, configVim.TenantID, configVim.DomainID)
				log.Trace("Loaded vim: ", openstackclient)
				vimList.addVim(configVim.Name, openstackclient)
			case string(Kubernetes):
				// TODO
			default:
				log.Error(config.ErrWrongVimType.Error())
			}
		}
	}

	log.Info("All vims loaded")
	return vimList
}
