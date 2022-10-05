CREATE TABLE `users` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `password` VARCHAR(100),
  `status` TINYINT(1),
  `role` TINYINT(1),
  `created_at` DATETIME,
  `updated_at` DATETIME
);

CREATE TABLE `user_auths` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` INT,
  `token` VARCHAR(200),
  `expired_at` DATETIME,
  `login_at` DATETIME,
  `created_at` DATETIME,
  `updated_at` DATETIME
);

CREATE TABLE `shops` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` INT,
  `created_at` DATETIME,
  `updated_at` DATETIME
);

CREATE TABLE `products` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `shop_id` INT,
  `name` VARCHAR(255),
  `price` double,
  `qty` INT,
  `created_at` DATETIME,
  `updated_at` DATETIME
);

CREATE TABLE `transactions` (
  `id` INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `product_id` INT,
  `amount` double,
  `user_id` INT,
  `status` INT,
  `qty` INT,
  `created_at` DATETIME,
  `updated_at` DATETIME
);

ALTER TABLE `user_auths` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `shops` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `products` ADD FOREIGN KEY (`shop_id`) REFERENCES `shops` (`id`);
