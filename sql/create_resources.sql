CREATE TABLE `resources` (
  `Name` varchar(255) NOT NULL,
  `Template` int(11) NOT NULL,
  `Data` mediumblob NOT NULL,
  PRIMARY KEY (`Name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8