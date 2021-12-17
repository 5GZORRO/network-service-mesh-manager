# SQL with 5 tables:
- _slices_: list of requested gateways 
```json
{
    "id": "",               primary_key
    "slice_id":"",          unique
    "status": "",  
    "vim_name": "",              different VIMs
    "created_at": "",
    "updated_at": "",
    vim gateways info
    # configuration
    "gw_mgmt_ip": "",
    "gw_mgmt_port": "",
    "gw_external_ip": "",
    "gw_exposed_nets":"",
    # VPN configuration
    "gw_vpn_port": "",
    "gw_vpn_interface": ""
}
```


- _networks_
```json
{   
    "id": "",
    "resource_set_id": "",
    "network_vim_id": "",
    "network_vim_name": "",
    "subnet_vim_id": "",
    "subnet_vim_name": "",
    "subnet_cidr": ""
}
```
- _saps_
```json
{   
    "id": "",
    "resource_set_id": "",
    "network_id": "",
    "network_name": "",
    "subnet_id": "",
    "subnet_name": "",
    "subnet_cidr": "",
    "router_id": "",
    "router_name": "",
    "router_id": "",
    "floating_net_name": "",
    "floating_net_id": ""
}
```

- _connections table_: list of connections
```json
{
    "id": "",
    "resource_set_id": "",
    "role": "",                 [server, client in the VPN connection]
    "status":"",
    "created_at":"",
    "peer_ip": "", 
    "peer_port": "",
    "allowed_ips": "",
    "public_key": "",
    "subnet_to_redirect": "",
    "created-at": ""
}
```

- _vims_
```json
{   
    "id": "",
    "vim_name": "",
    "vim_type": "",
    "endpoint": "",
    "username": "",
    "password": "",
    "domain": "",
    "tenant": "",

}