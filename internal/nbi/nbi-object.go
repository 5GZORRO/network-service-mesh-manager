package nbi

import (
	"errors"
	"nextworks/nsm/internal/openstackclient"

	log "github.com/sirupsen/logrus"
)

// Env object should contain the environment used by the functions
// associated to REST API, such as a ConnectionPool to a DB or a
// Provider to OpenStackAPI
type Env struct {
	Client *openstackclient.OpenStackClient
	DB     []GatewayConnectivity
}

// AddSliceConnectivity add a new GatewayConnectivity object in the DB slice
// each GatewayConnectivity is uniquely identified by the SliceID
func (env *Env) AddGatewayConnectivityInDB(sliceID string) (*GatewayConnectivity, error) {
	// check if the slice already exysts
	_, slice, _ := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if slice != nil {
		log.Error("A slice with name ", sliceID, " already exists in the DB")
		return nil, errors.New("slice already exists")
	}
	// add the new slice
	gc := GatewayConnectivity{SliceID: sliceID}
	env.DB = append(env.DB, gc)
	log.Info("Inserted a new slice with name ", sliceID, " [DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
	return &gc, nil
}

func (env *Env) RetrieveGatewayConnectivityFromDB(sliceID string) (int, *GatewayConnectivity, error) {
	for i := range env.DB {
		if env.DB[i].SliceID == sliceID {
			return i, &env.DB[i], nil
		}
	}
	return -1, nil, errors.New("slice not found")
}

func (env *Env) RemoveGatewayConnectivityFromDB(sliceID string) (*GatewayConnectivity, error) {
	index, slice, err := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if err != nil {
		return nil, errors.New("slice not found")
	}

	ret := make([]GatewayConnectivity, 0)
	ret = append(ret, env.DB[:index]...)
	env.DB = append(ret, env.DB[index+1:]...)
	log.Info("Deleted slice with name ", sliceID, " [DB len: ", len(env.DB), " capacity: ", cap(env.DB), "] \nDB:", env.DB)
	return slice, nil
}

func (env *Env) RemoveAllGatewayConnectivitiesFromDB() {
	env.DB = []GatewayConnectivity{}
}

func (env *Env) UpdateGatewayConnectivityInDB(sliceID string, privnetID string, subnetID string, routerID string, portID string) (*GatewayConnectivity, error) {
	_, slice, err := env.RetrieveGatewayConnectivityFromDB(sliceID)
	if err != nil {
		return nil, errors.New("slice not found")
	}
	slice.PrivNetID = privnetID
	slice.SubnetID = subnetID
	slice.RouterID = routerID
	slice.InterfaceID = portID
	log.Info("Updated ", sliceID, " slice: ", slice, " \nDB:", env.DB)
	return slice, nil
}
