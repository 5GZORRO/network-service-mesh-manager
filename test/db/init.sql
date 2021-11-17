DROP TABLE IF EXISTS `gateways`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gateways` (
  `sliceId` varchar(255) NOT NULL,
  `networkId` varchar(255) DEFAULT NULL,
  `subnetId` varchar(255) DEFAULT NULL,
  `routerId` varchar(255) DEFAULT NULL,
  `interfaceId` varchar(255) DEFAULT NULL,
  `floatingIp` varchar(255) DEFAULT NULL,
  `vmGatewayId` varchar(255) DEFAULT NULL,
  UNIQUE KEY `sliceId` (`sliceId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
