CREATE TABLE `blogs` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(255) NOT NULL,
  `Url` varchar(255) NOT NULL,
  `Template` int(11) NOT NULL,
  `Keywords` varchar(255) NOT NULL,
  `Description` varchar(255) NOT NULL,
  `Slogan` varchar(255) NOT NULL,
  `Server` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8