from openstack.exceptions import OpenStackCloudException

# Global params
project_id='7953babdca974e7ab44cc6c69f093956'
is_shared=False
admin_state=True
availability_zone=['nova']

# Function to create a network and a subnet
def create_network(conn, network_name, cidr):
    print("create_network: "+network_name)
    try:
        network = conn.network.create_network(name=network_name,
                                            is_shared=is_shared,
                                            is_admin_state_up=admin_state,
                                            project_id=project_id,
                                            availability_zone_hints=availability_zone)
        if network != None:
            print("create_network ID: "+network.id)
            subnet_name=network_name+"_subnet"
        subnet=create_subnet(conn,subnet_name,network_id=network.id,cidr=cidr)
        return True
    except OpenStackCloudException:
        return False

def create_subnet(conn, subnet_name, network_id, cidr):
    print("create_subnet: "+subnet_name)
    subnet = conn.network.create_subnet(name=subnet_name,
                                            ip_version=4,
                                            network_id=network_id,
                                            project_id=project_id,
                                            is_dhcp_enabled=True,
                                            cidr=cidr)
    return subnet

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

def delete_subnet(conn, subnet_name):
    gen = conn.network.subnets(name=subnet_name,project_id=project_id)
    subns = list(gen)
    if len(subns) == 0:
        print("No network with name "+subnet_name+" found")
        return False

    elif len(subns) == 1:
        print("Network with name "+subnet_name+" found")
        print(">> SUBNET <<")
        print(subns[0].to_dict())
        conn.network.delete_subnet(subns[0],ignore_missing=True)
        return True
    else:
        print("More than one network with name "+subnet_name+" found")
        return False
