from flask import Flask, request, Response, request
from ipaddress import IPv4Network
import json
import requests
import sys


# monkey.patch_all()
app = Flask(__name__)

@app.route('/launch', methods=['POST'])
def launch():
    response = Response()
    req = request.data.decode("utf-8")
    req = json.loads(req)
    ip_range = req["ip_range"]
    net_interface = req["net_interface"]
    port = req["port"]
    print("Starting VPN Server on interface ",net_interface,"on port: ",port, " with IP Range: ",ip_range)
    response.status_code = 200
    return response

        
@app.route('/get_configuration', methods=['GET'])
def getconfig():
    response = Response()
    print("Get configuration")
    data = {
        "did": "did:5gzorro:dummy12345",
        "public_key": "svs",
        "IP_range": "vreg",
        "vpn_port": ""
    }
    return json.dumps(data)


@app.route('/connect_to_VPN', methods=['POST'])
def connect():
    response = Response()
    req = request.data.decode("utf-8")
    req = json.loads(req)
    ip_address_server = req["ip_address_server"]
    port_server = req["port_server"]
    IP_range_to_redirect = req["IP_range_to_redirect"]
    print("Connecting to peer:  ",ip_address_server," ",port_server, " ",IP_range_to_redirect)
    response.status_code = 200
    return response


@app.route('/disconnect_to_VPN', methods=['POST'])
def disconnect():
    response = Response()
    req = request.data.decode("utf-8")
    req = json.loads(req)
    ip_address_server = req["ip_address_server"]
    port_server = req["port_server"]

    print("Disconnecting from peer:  ",ip_address_server," ",port_server)
    response.status_code = 200
    return response


def launch_server_REST(port):
    # api.add_resource(launch, '/launch')
    # api.add_resource(get_configuration, '/get_configuration')
    # api.add_resource(connect_to_VPN, '/connect_to_VPN')
    # api.add_resource(disconnect_to_VPN, '/disconnect_to_VPN')
    app.run(host='0.0.0.0', port=port, debug=True)


if __name__ == "__main__":
    if len(sys.argv)!=2:
        print("Usage: python3 app_api.py [port]")
    else:
        port=int(sys.argv[1])
        print("Starting Server on port: ",port)
        launch_server_REST(port)


