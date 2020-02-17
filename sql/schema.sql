CREATE DATABASE IF NOT EXISTS `monzo_plus_plus` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `monzo_plus_plus`;

-- MySQL dump 10.13  Distrib 5.7.29, for Win64 (x86_64)
--
-- Host: %    Database: monzo_plus_plus
-- ------------------------------------------------------
-- Server version	8.0.18
 /*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;

/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;

/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;

/*!40101 SET NAMES utf8 */;

/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;

/*!40103 SET TIME_ZONE='+00:00' */;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;

/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;

/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `user_tokens`
--

DROP TABLE IF EXISTS `user_tokens`;

/*!40101 SET @saved_cs_client     = @@character_set_client */;

/*!40101 SET character_set_client = utf8 */;


CREATE TABLE `user_tokens`
  (`user_id` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
                                                  `access_token` text COLLATE utf8_unicode_ci NOT NULL,
                                                                                              `refresh_token` text COLLATE utf8_unicode_ci NOT NULL,
                                                                                                                                           `expires` datetime NOT NULL,
                                                                                                                                                              `token_type` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
                                                                                                                                                                                                               PRIMARY KEY (`user_id`), CONSTRAINT `fk_user_token`
   FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE) ENGINE=InnoDB DEFAULT
CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;

/*!40101 SET @saved_cs_client     = @@character_set_client */;

/*!40101 SET character_set_client = utf8 */;


CREATE TABLE `users` (`id` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
                                                                `monzo_id` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
                                                                                                                `state_token` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
                                                                                                                                                                   PRIMARY KEY (`id`), UNIQUE KEY `monzo_id_UNIQUE` (`monzo_id`),
                                                                                                                                                                                                  UNIQUE KEY `stateToken_UNIQUE` (`state_token`)) ENGINE=InnoDB DEFAULT
CHARSET=utf8 COLLATE=utf8_unicode_ci;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping events for database 'monzo_plus_plus'
--
 --
-- Dumping routines for database 'monzo_plus_plus'
--
/*!50003 DROP PROCEDURE IF EXISTS `create_new_user` */;

/*!50003 SET @saved_cs_client      = @@character_set_client */ ;

/*!50003 SET @saved_cs_results     = @@character_set_results */ ;

/*!50003 SET @saved_col_connection = @@collation_connection */ ;

/*!50003 SET character_set_client  = utf8mb4 */ ;

/*!50003 SET character_set_results = utf8mb4 */ ;

/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;

/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;

/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION' */ ;


DELIMITER ;

;


CREATE DEFINER=`user`@`%` PROCEDURE `create_new_user`(in id varchar(128), in state_token varchar(128)) BEGIN
INSERT INTO users (`id`,`monzo_id`, `state_token`)
VALUES (id,
        '',
        state_token); END ;

;


DELIMITER ;

/*!50003 SET sql_mode              = @saved_sql_mode */ ;

/*!50003 SET character_set_client  = @saved_cs_client */ ;

/*!50003 SET character_set_results = @saved_cs_results */ ;

/*!50003 SET collation_connection  = @saved_col_connection */ ;

/*!50003 DROP PROCEDURE IF EXISTS `get_user_by_state_token` */;

/*!50003 SET @saved_cs_client      = @@character_set_client */ ;

/*!50003 SET @saved_cs_results     = @@character_set_results */ ;

/*!50003 SET @saved_col_connection = @@collation_connection */ ;

/*!50003 SET character_set_client  = utf8mb4 */ ;

/*!50003 SET character_set_results = utf8mb4 */ ;

/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;

/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;

/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION' */ ;


DELIMITER ;

;


CREATE DEFINER=`user`@`%` PROCEDURE `get_user_by_state_token`(in state_token varchar(128)) BEGIN
SELECT id,
       `state_token`
FROM users
WHERE `state_token` = state_token; END ;

;


DELIMITER ;

/*!50003 SET sql_mode              = @saved_sql_mode */ ;

/*!50003 SET character_set_client  = @saved_cs_client */ ;

/*!50003 SET character_set_results = @saved_cs_results */ ;

/*!50003 SET collation_connection  = @saved_col_connection */ ;

/*!50003 DROP PROCEDURE IF EXISTS `update_user` */;

/*!50003 SET @saved_cs_client      = @@character_set_client */ ;

/*!50003 SET @saved_cs_results     = @@character_set_results */ ;

/*!50003 SET @saved_col_connection = @@collation_connection */ ;

/*!50003 SET character_set_client  = utf8mb4 */ ;

/*!50003 SET character_set_results = utf8mb4 */ ;

/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;

/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;

/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION' */ ;


DELIMITER ;

;


CREATE DEFINER=`user`@`%` PROCEDURE `update_user`(IN user_id VARCHAR(128), IN monzo_id VARCHAR(128), IN access_token TEXT, IN refresh_token TEXT, IN expires DATETIME, IN token_type VARCHAR(45)) BEGIN DECLARE hasToken BOOL;
SET hasToken = EXISTS
  (SELECT *
   FROM user_tokens
   WHERE `user_id` = user_id);
UPDATE users
SET `monzo_id` = monzo_id
WHERE `id` = user_id; IF hasToken THEN
  UPDATE user_tokens
  SET `access_token` = access_token,
      `refresh_token` = refresh_token,
      `expires` = expires,
      `token_type` = token_type
  WHERE `user_id` = user_id; ELSE
    INSERT INTO user_tokens (`user_id`, `access_token`, `refresh_token`, `expires`, `token_type`)
  VALUES (user_id,
          access_token,
          refresh_token,
          expires,
          token_type); END IF; END ;

;


DELIMITER ;

/*!50003 SET sql_mode              = @saved_sql_mode */ ;

/*!50003 SET character_set_client  = @saved_cs_client */ ;

/*!50003 SET character_set_results = @saved_cs_results */ ;

/*!50003 SET collation_connection  = @saved_col_connection */ ;

/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;

/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;

/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;

/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2020-02-14 17:10:26
