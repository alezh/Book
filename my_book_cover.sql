/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50714
Source Host           : localhost:3306
Source Database       : booklist

Target Server Type    : MYSQL
Target Server Version : 50714
File Encoding         : 65001

Date: 2018-02-08 09:57:20
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for my_book_cover
-- ----------------------------
DROP TABLE IF EXISTS `my_book_cover`;
CREATE TABLE `my_book_cover` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  `author` varchar(32) DEFAULT NULL,
  `catalog_url` varchar(255) DEFAULT NULL,
  `status` tinyint(1) DEFAULT NULL,
  `desc` varchar(255) DEFAULT NULL,
  `cover_img` varchar(255) DEFAULT NULL,
  `new_chapter` varchar(32) DEFAULT NULL,
  `sort` varchar(32) DEFAULT NULL,
  `favorite` int(11) DEFAULT NULL,
  `hits` int(11) DEFAULT NULL,
  `created` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `updated` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
