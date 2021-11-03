import openstack
import openstack_functions as op
import logging
from flask import Flask, Response, make_response, request


logging.basicConfig(format='[%(name)s] %(asctime)s - %(message)s', level=logging.DEBUG)
module_name = "NSMM"
logger = logging.getLogger(module_name)


# Initialize and turn on debug logging
openstack.enable_logging(debug=True, http_debug=True)

# Initialize connection with OS
conn = openstack.connect(cloud='local')


app = Flask(__name__)



@app.route('/network', methods=['POST', 'DELETE', 'GET'])
def networks():
    response = Response()
    if request.method == 'POST':
        response.status_code = 200
        network_name = request.args.get('name', None)
        if network_name is not None:
            data = request.get_json()
            network_nameB = data.get("name")
            network_cidr = data.get("cidr")
            if network_name == network_nameB:
                logger.info("create_network : name: "+network_name+" cidr: "+network_cidr)
                # TODO call to create
                result = op.create_network(conn,network_name,cidr=network_cidr)
                if not result:
                    response.status_code = 404
            else:
                response.status_code = 400
        else:
            response.status_code = 400
        return response
    elif request.method == 'DELETE':
        logger.debug("delete_network")
        network_name = request.args.get('name', None)
        if network_name is not None:
            result = op.delete_network(conn,network_name)
            if result:
                response.status_code = 200
            else:
                response.status_code = 404
        else:
            response.status_code = 400
        return response
    elif request.method == 'GET':
        logger.debug("retrieve_network")
        network_name = request.args.get('name', None)
        if network_name is not None:
            network = op.retrieve_network(conn,network_name)
            if network is not None:
                return network 
            else:
                response.status_code = 404
        else:
            response.status_code = 400
        return response

app.run(host='0.0.0.0', port=8080, debug=True)
# conn.close()