package nbi

import (
	"errors"
	"nextworks/nsm/internal/openstackclient"

	log "github.com/sirupsen/logrus"
)

// GatewayConnectivity represent all connectivity information for an interdomain connection:
// for each slice, I assume a priv_net with a subnet
// router connected to the floating net, so with a port
// a gateway vm
// a floating ip
type GatewayConnectivity struct {
	SliceID     string `json:"sliceID" binding:"required"`
	PrivNetID   string `json:"networkID" binding:"required"`
	SubnetID    string `json:"subnetID" binding:"required"`
	RouterID    string `json:"routerID" binding:"required"`
	InterfaceID string `json:"interfaceID"`
	FloatingIP  string `json:"floatingIP"`
	VmGatewayID string `json:"vmID"`
}

// Env object should contain the environment used by the functions
// associated to REST API, such as a ConnectionPool to a DB or a
// Provider to OpenStackAPI
type Env struct {
	Client *openstackclient.OpenStackClient
	DB     []GatewayConnectivity
}

// AddGatewayConnectivityInDB add a new GatewayConnectivity object in the DB
// each GatewayConnectivity is uniquely identified by the SliceID
func (env *Env) AddGatewayConnectivityInDB(sliceID string) (*GatewayConnectivity, error) {
	// check if the gateway connectivity object for sliceID already exysts
	_, gateway, _ := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if gateway != nil {
		log.Error("A GatewayConnectivity for the slice with sliceID ", sliceID, " already exists in the DB")
		return nil, errors.New("Slice with sliceID " + sliceID + " has already a Gateway")
	}
	// add the new gatewayconnectivity obj
	gc := GatewayConnectivity{SliceID: sliceID}
	env.DB = append(env.DB, gc)
	log.Info("Inserted a new GatewayConnectivity for slice with sliceID ", sliceID, " \n[DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
	return &gc, nil
}

func (env *Env) RetrieveGatewayConnectivityFromDB(sliceID string) (int, *GatewayConnectivity, error) {
	for i := range env.DB {
		if env.DB[i].SliceID == sliceID {
			return i, &env.DB[i], nil
		}
	}
	return -1, nil, errors.New("GatewayConnectivty information for slice with SliceID " + sliceID + " not found")
}

func (env *Env) RemoveGatewayConnectivityFromDB(sliceID string) (*GatewayConnectivity, error) {
	index, gc, err := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if err != nil {
		log.Error("Error during removal of GatewayConnectivity info of slice with sliceID: ", sliceID)
		return nil, err
	}

	ret := make([]GatewayConnectivity, 0)
	ret = append(ret, env.DB[:index]...)
	env.DB = append(ret, env.DB[index+1:]...)
	log.Info("Deleted GatewayConnectivity for the slice with sliceID ", sliceID, "\n[DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
	return gc, nil
}

func (env *Env) RemoveAllGatewayConnectivitiesFromDB() {
	env.DB = []GatewayConnectivity{}
	log.Info("Deleted the GatewayConnectivity info for all the slices")
}

func (env *Env) UpdateGatewayConnectivityInDB(sliceID string, privnetID string, subnetID string, routerID string, portID string) (*GatewayConnectivity, error) {
	_, gc, err := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if err != nil {
		return nil, errors.New("GatewayConnectivty information for slice with SliceID " + sliceID + " not found")
	}
	gc.PrivNetID = privnetID
	gc.SubnetID = subnetID
	gc.RouterID = routerID
	gc.InterfaceID = portID
	log.Info("Updated GatewayConnectivity for slice with sliceID ", sliceID, " gateway: ", gc, " \nDB:", env.DB)
	return gc, nil
}
