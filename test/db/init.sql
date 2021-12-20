/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE resource_sets (
  id INT NOT NULL AUTO_INCREMENT,
  slice_id varchar(255) NOT NULL UNIQUE,
  `status` varchar(255) NOT NULL,
  vim_name varchar(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  gw_mgmt_ip varchar(255) DEFAULT NULL,
  gw_mgmt_port smallint DEFAULT NULL,
  gw_external_ip varchar(255) DEFAULT NULL,
  gw_exposed_nets varchar(255) DEFAULT NULL,
  gw_pub_key varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY id (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE networks (
  id INT NOT NULL AUTO_INCREMENT,
  resource_set_id INT NOT NULL,
  network_id varchar(255) DEFAULT NULL,
  network_name varchar(255) DEFAULT NULL,
  subnet_id varchar(255) DEFAULT NULL,
  subnet_name varchar(255) DEFAULT NULL,
  subnet_cidr varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY id (id),
  FOREIGN KEY (resource_set_id) REFERENCES resource_sets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE saps (
  id INT NOT NULL AUTO_INCREMENT,
  resource_set_id INT NOT NULL,
  network_id varchar(255) DEFAULT NULL,
  network_name varchar(255) DEFAULT NULL,
  subnet_id varchar(255) DEFAULT NULL,
  subnet_name varchar(255) DEFAULT NULL,
  subnet_cidr varchar(255) DEFAULT NULL,
  router_id varchar(255) DEFAULT NULL,
  router_name varchar(255) DEFAULT NULL,
  router_port_id varchar(255) DEFAULT NULL,
  floating_net_id varchar(255) DEFAULT NULL,
  floating_net_name varchar(255) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY id (id),
  FOREIGN KEY (resource_set_id) REFERENCES resource_sets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE connections (
  id INT NOT NULL AUTO_INCREMENT,
  resource_set_id INT NOT NULL,
  `role` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  peer_ip varchar(255) DEFAULT NULL,
  peer_port smallint DEFAULT NULL,
  allowed_ips varchar(255) NOT NULL,
  public_key varchar(255) NOT NULL,
  subnet_to_redirect varchar(255) NOT NULL,
  PRIMARY KEY (id),
  KEY id (id),
  FOREIGN KEY (resource_set_id) REFERENCES resource_sets(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;