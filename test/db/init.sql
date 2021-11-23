DROP TABLE IF EXISTS `gateways`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateways` (
  `slice_id` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `network_id` varchar(255) DEFAULT NULL,
  `network_name` varchar(255) DEFAULT NULL,
  `subnet_id` varchar(255) DEFAULT NULL,
  `subnet_name` varchar(255) DEFAULT NULL,
  `subnet_cidr` varchar(255) DEFAULT NULL,
  `router_id` varchar(255) DEFAULT NULL,
  `router_name` varchar(255) DEFAULT NULL,
  `port_id` varchar(255) DEFAULT NULL,
  `floating_ip` varchar(255) DEFAULT NULL,
  `vm_gateway_id` varchar(255) DEFAULT NULL,
  `gateway_role` varchar(255) DEFAULT NULL,
  `remote_endpoint` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`slice_id`),
  KEY `slice_id` (`slice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

