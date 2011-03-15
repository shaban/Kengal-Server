CREATE TABLE `rubrics` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(255) NOT NULL,
  `Url` varchar(255) NOT NULL,
  `Keywords` varchar(255) NOT NULL,
  `Description` varchar(255) NOT NULL,
  `Blog` int(11) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1005 DEFAULT CHARSET=utf8