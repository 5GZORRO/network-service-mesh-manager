DROP TABLE IF EXISTS `gateways`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateways` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `slice_id` varchar(255) NOT NULL UNIQUE,
  `status` varchar(255) NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  `network_id` varchar(255) DEFAULT NULL,
  `subnet_id` varchar(255) DEFAULT NULL,
  `router_id` varchar(255) DEFAULT NULL,
  `external_ip` varchar(255) DEFAULT NULL,
  `management_ip` varchar(255) DEFAULT NULL,
  `management_port` smallint DEFAULT NULL,
  `vpn_port` varchar(255) DEFAULT NULL,
  `vpn_interface` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `networks` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `vim_network_id` varchar(255) DEFAULT NULL,
  `vim_network_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `subnets` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `vim_subnet_id` varchar(255) DEFAULT NULL,
  `vim_subnet_name` varchar(255) DEFAULT NULL,
  `subnet_cidr` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `routers` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `vim_router_id` varchar(255) DEFAULT NULL,
  `vim_router_name` varchar(255) DEFAULT NULL,
  `vim_router_port_id` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
