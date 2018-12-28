CREATE DATABASE `prepaid-card` IF NOT EXISTS;
# Dump of table account
# ------------------------------------------------------------

DROP TABLE IF EXISTS `account`;

CREATE TABLE `account` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `account_holder_name` varchar(255) NOT NULL DEFAULT '',
  `card_number` int(11) NOT NULL,
  `sort_code` int(11) NOT NULL,
  `status` enum('INACTIVE','ACTIVE','DISABLED') NOT NULL DEFAULT 'INACTIVE',
  `balance` decimal(13,2) DEFAULT '0.00',
  `address` varchar(255) DEFAULT NULL,
  `postcode` varchar(100) DEFAULT NULL,
  `date_of_birth` date DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

# Dump of table blocked_transactions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `blocked_transactions`;

CREATE TABLE `blocked_transactions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` varchar(255) DEFAULT NULL,
  `amount` decimal(13,2) DEFAULT NULL,
  `reason` text,
  `blocked_at` datetime DEFAULT NULL,
  `transaction_id` varchar(255) DEFAULT NULL,
  `card_number` varchar(255) DEFAULT NULL,
  `status` enum('BLOCKED','CAPTURED','REVERSED') NOT NULL DEFAULT 'BLOCKED',
  `original_transaction_id` int(11) DEFAULT NULL,
  `balance` decimal(13,2) DEFAULT NULL,
  `parent_transaction_id` varchar(255) DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table merchants
# ------------------------------------------------------------

DROP TABLE IF EXISTS `merchants`;

CREATE TABLE `merchants` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` varchar(255) DEFAULT NULL,
  `status` enum('ACTIVE','INACTIVE') DEFAULT 'ACTIVE',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `merchants` WRITE;
/*!40000 ALTER TABLE `merchants` DISABLE KEYS */;

INSERT INTO `merchants` (`id`, `merchant_id`, `status`)
VALUES
	(2,'123456789','ACTIVE'),
	(4,'123456788','ACTIVE');

/*!40000 ALTER TABLE `merchants` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table topups
# ------------------------------------------------------------

DROP TABLE IF EXISTS `topups`;

CREATE TABLE `topups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `amount` decimal(13,2) NOT NULL,
  `topup_at` datetime NOT NULL,
  `card_number` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table transactions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `transactions`;

CREATE TABLE `transactions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `card_number` varchar(255) NOT NULL DEFAULT '',
  `blocked_transaction_id` varchar(255) DEFAULT NULL,
  `merchant_id` int(11) NOT NULL,
  `amount` double(13,2) NOT NULL,
  `status` enum('CAPTURED','REFUND') NOT NULL DEFAULT 'CAPTURED',
  `balance` decimal(13,3) DEFAULT NULL,
  `transaction_id` varchar(255) NOT NULL DEFAULT '',
  `captured_at` datetime DEFAULT NULL,
  `parent_transaction_id` varchar(255) DEFAULT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;