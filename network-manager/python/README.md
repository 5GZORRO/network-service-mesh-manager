# Network Manager Module
First version of __NetworkManager__ module implemented using Python and https://docs.openstack.org/openstacksdk/latest/ and Flask.

It exposes 3 API to interact with OpenStack and create/retrieve/delete networks and subnets:
```
GET network?name={name}
POST network?name={name} + body
DELETE network?name={name}
```

These API are described in the Postman collection `NSMM - NetworkManager`

# Run
Install required dependencies
```
pip install openstacksdk
pip install flask
```

```
python -m openstack version
> OpenstackSDK Version 0.59.0
```

## Run 
```
python3 main.py
```