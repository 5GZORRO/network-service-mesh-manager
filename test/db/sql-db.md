# SQL with 5 tables:

- _gateways table_: list of requested gateways 
    ```json
    {
        "id": "",               primary_key
        "slice_id":"",          unique
        "status": "",  
        "vim_name": "",              different VIMs
        "created_at": "",
        "updated_at": "",
        vim gateways info
        "vim_resource_id": "",
        # configuration
        "management_ip": "",
        "management_port": "",
        "external_ip": "",
        # VPN configuration
        "vpn_server_port": "",
        "vpn_server_interface": ""
    }
    ```


- __openstack_resources__
    ```json
    {   
        "id": "",
        "network_vim_id": "",
        "network_vim_name": "",
        "subnet_vim_id": "",
        "subnet_vim_name": "",
        "subnet_cidr": "",
        "router_vim_id": "",
        "router_vim_name": "",
        "router_vim_port_id": ""
    }

- _connections table_: list of connections
    ```json
    {
        "id": "",
        "gateway_id": "",
        "role": "",                 [server, client in the VPN connection]
        "remote_peer": "", 
        "remote_port": "",
        "subnet_to_redirect": "",
        "connection_status": "",
        "created-at": "",
        "updated-at": ""
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