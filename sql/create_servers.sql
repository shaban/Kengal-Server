CREATE TABLE `servers` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `IP` varchar(255) NOT NULL,
  `Vendor` varchar(255) NOT NULL,
  `Type` varchar(255) NOT NULL,
  `Item` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8