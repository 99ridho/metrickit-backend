CREATE TABLE `metadata` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `version` varchar(255) NOT NULL DEFAULT '',
  `build_number` varchar(255) NOT NULL DEFAULT '',
  `device_type` varchar(255) NOT NULL DEFAULT '',
  `os` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `app_launch_time_first_draw` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `metadata_id` int(11) unsigned NOT NULL,
  `range_start` double NOT NULL,
  `range_end` double NOT NULL,
  `frequency` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_fd_metadata_id` (`metadata_id`),
  CONSTRAINT `fk_fd_metadata_id` FOREIGN KEY (`metadata_id`) REFERENCES `metadata` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `app_launch_time_resume` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `metadata_id` int(11) unsigned NOT NULL,
  `range_start` double NOT NULL,
  `range_end` double NOT NULL,
  `frequency` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_r_metadata_id` (`metadata_id`),
  CONSTRAINT `fk_r_metadata_id` FOREIGN KEY (`metadata_id`) REFERENCES `metadata` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `app_signpost` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `metadata_id` int(11) unsigned NOT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `category` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `fk_sp_metadata_id` (`metadata_id`),
  CONSTRAINT `fk_sp_metadata_id` FOREIGN KEY (`metadata_id`) REFERENCES `metadata` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `app_signpost_interval` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `signpost_id` int(11) unsigned NOT NULL,
  `average_memory` double NOT NULL,
  `cumulative_cpu_time` double NOT NULL,
  `cumulative_logical_writes` double NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_spi_signpost_interval` (`signpost_id`),
  CONSTRAINT `fk_spi_signpost_interval` FOREIGN KEY (`signpost_id`) REFERENCES `app_signpost` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `app_signpost_histogram` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `signpost_interval_id` int(11) unsigned NOT NULL,
  `range_start` double NOT NULL,
  `range_end` double NOT NULL,
  `frequency` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_sph_signpost_intv_id` (`signpost_interval_id`),
  CONSTRAINT `fk_sph_signpost_intv_id` FOREIGN KEY (`signpost_interval_id`) REFERENCES `app_signpost_interval` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;