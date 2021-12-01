/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE openstack_resources (
  id INT NOT NULL AUTO_INCREMENT,
  network_vim_id varchar(255) DEFAULT NULL,
  network_vim_name varchar(255) DEFAULT NULL,
  subnet_vim_id varchar(255) DEFAULT NULL,
  subnet_vim_name varchar(255) DEFAULT NULL,
  subnet_cidr varchar(255) DEFAULT NULL,
  router_vim_id varchar(255) DEFAULT NULL,
  router_vim_name varchar(255) DEFAULT NULL,
  router_vim_port_id varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY id (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE gateways (
  id INT NOT NULL AUTO_INCREMENT,
  slice_id varchar(255) NOT NULL UNIQUE,
  `status` varchar(255) NOT NULL,
  vim_name varchar(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  vim_resource_id INT DEFAULT NULL,
  external_ip varchar(255) DEFAULT NULL,
  management_ip varchar(255) DEFAULT NULL,
  management_port smallint DEFAULT NULL,
  vpn_server_port smallint DEFAULT NULL,
  vpn_server_interface varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY id (id),
  FOREIGN KEY (vim_resource_id) REFERENCES openstack_resources(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE connections (
  id INT NOT NULL AUTO_INCREMENT,
  gateway_id INT NOT NULL,
  `role` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  remote_server_ip varchar(255) DEFAULT NULL,
  remote_server_port varchar(255) DEFAULT NULL,
  allowed_ips varchar(255) NOT NULL,
  public_key varchar(255) NOT NULL,
  subnet_to_redirect varchar(255) NOT NULL,
  PRIMARY KEY (id),
  KEY id (id),
  FOREIGN KEY (gateway_id) REFERENCES gateways(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;