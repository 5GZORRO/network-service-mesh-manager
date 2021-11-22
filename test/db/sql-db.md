# SQL with 5 tables:

- _gateways table_: list of requested gateway connectivity
    ```json
    {
        "slice_id":"",  primary_key
        "status": "",  ["creating","ready","configuring","running","zombie"]
        "network_id": "", 
        "subnet_id": "",  
        "router_id": "", 
        "port_id": "",
        "floating_ip_id": "",
        "gateway_vm_id": ""

    }
    ```

- _networks table_ 
    ```json
    {
        "network_id": "",  
        "network_name": ""  
    }
    ```
- _subnets table_ 
    ```json
    {
        "subnet_id": "",  
        "subnet_name": "",  
        "subnet_cidr": ""
    }
    ```
- _routers table_ --> 
    ```json
    {
        "router_id": "",  
        "router_name": ""  
    }
    ```
- _floating_ips table_
    ```json
    {
        "id": "",
        "slice_id": "",
        "floating_ip": "",
        "status": ""
    }
    ```
- _gateway_vms table_
    ```json
    {
        "id": "",
        "slice_id": "",
        "vm_id": "",
        "status": "" ["ready","configured","running"]
    }
    ```