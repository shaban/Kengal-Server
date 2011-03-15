CREATE TABLE `articles` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(255) NOT NULL,
  `Rubric` int(11) NOT NULL,
  `Text` text NOT NULL,
  `Teaser` text NOT NULL,
  `Blog` int(11) NOT NULL,
  `Keywords` varchar(255) NOT NULL,
  `Description` varchar(255) NOT NULL,
  `Date` datetime NOT NULL,
  `Url` varchar(255) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=10011 DEFAULT CHARSET=utf8