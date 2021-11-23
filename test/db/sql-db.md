# SQL with 5 tables:

- _gateways table_: list of requested gateways 
    ```json
    {
        "slice_id":"",  primary_key
        "status": "",  
        "network_id": "",
        "network_name": "", 
        "subnet_id": "",
        "subnet_name": "",  
        "subnet_cidr": "",
        "router_id": "",
        "router_name": "", 
        "port_id": "",
        "floating_ip_id": "",
        "gateway_vm_id": "",
        "gateway_role": "", [server, client in the VPN connection]
        "remote_endpoint": ""
    }
    ```