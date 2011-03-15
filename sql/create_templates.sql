CREATE TABLE `templates` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Index` text NOT NULL,
  `Style` text NOT NULL,
  `Title` varchar(255) NOT NULL,
  `FromUrl` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8