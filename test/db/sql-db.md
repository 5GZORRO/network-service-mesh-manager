# SQL with 5 tables:

- _gateways table_: list of requested gateways 
    ```json
    {
        "id": "",               primary_key
        "slice_id":"",          unique
        "status": "",  
        "vim": "",              different VIMs
        "created-at": "",
        "updated-at": "",
        gateways info
        "network_id": "",
        "subnet_id": "",
        "router_id": "",
        # configuration
        "management_ip": "",
        "management_port": "",
        "external_ip": "",
        # VPN configuration
        "server_ip": "",        should be external_ip
        "server_port": "",
        "server_interface": ""
    }
    ```


- _networks_
    ```json
    {   
        "id": "",
        "vim_network_id": "",
        "vim_network_name": "",
        "created-at": "",
        "updated-at": ""
    }
    ```

- _subnets_
    ```json
    {   
        "id": "",
        "vim_subnet_id": "",
        "vim_subnet_name": "", 
        "subnet_cidr": "",
        "created-at": "",
        "updated-at": ""
    }
    ```

- _routers_
    ```json
    {   
        "id": "",
        "vim_router_id": "",
        "vim_router_name": "",
        "vim_port_id": "",
        "created-at": "",
        "updated-at": ""
    }
    ```

- _connections table_: list of connections
    ```json
    {
        "id": "",
        "gateway_id": "",
        "slice_id": "",
        "role": "",                 [server, client in the VPN connection]
        "remote_peer": "", 
        "remote_port": "",
        "subnet_to_redirect": "",
        "connection_status": "",
        "created-at": "",
        "updated-at": ""
    }
    ```