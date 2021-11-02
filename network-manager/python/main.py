import openstack

# Global params
project_id='7953babdca974e7ab44cc6c69f093956'
is_shared=False
admin_state=True
availability_zone=['nova']

# Initialize and turn on debug logging
openstack.enable_logging(debug=True, http_debug=True)

# Initialize connection with OS
conn = openstack.connect(cloud='local')

# Function to create a network and a subnet
# TODO create subnet
def create_network(conn, network_name, cidr):
    print("create_network: "+network_name)
    network = conn.network.create_network(name=network_name,
                                            is_shared=is_shared,
                                            is_admin_state_up=admin_state,
                                            project_id=project_id,
                                            availability_zone_hints=availability_zone)
    return network

# Function to retrieve a network by its name
def retrieve_network(conn, network_name):
    print("retrieve_network: "+network_name)
    
    gen = conn.network.networks(name=network_name,is_shared=False)
    nets = list(gen)
    if len(nets) == 0:
        print("No network with name "+network_name+" found")
        return None

    elif len(nets) == 1:
        print("Network with name "+network_name+" found")
        return nets[0]
    else:
        print("More than one network with name "+network_name+" found")
        return None

# TODO: explicit delete of subnet
def delete_network(conn, network_name):
    print("delete_network: "+network_name)

    gen = conn.network.networks(name=network_name,is_shared=False)
    nets = list(gen)
    if len(nets) == 0:
        print("No network with name "+network_name+" found")
        return False

    elif len(nets) == 1:
        print("Network with name "+network_name+" found")
        conn.network.delete_network(nets[0],ignore_missing=True)
        return True
    else:
        print("More than one network with name "+network_name+" found")
        return False


print("")
print("Retrieve_network")
network = create_network(conn,"test_network","")
print(network.to_dict())
network = retrieve_network(conn,"test_network")
print(network.to_dict())
result = delete_network(conn,"test_network")
print(result)

conn.close()

