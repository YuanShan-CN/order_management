-- Migration: Add user_id to orders and shops tables
-- Run this migration if you have existing data

-- Create users table if not exists
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(50) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  INDEX `idx_users_username` (`username`),
  INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Add user_id to orders table (allow NULL for existing data, then update)
ALTER TABLE `orders` ADD COLUMN `user_id` INT UNSIGNED NOT NULL DEFAULT 1 AFTER `id`;
ALTER TABLE `orders` ADD INDEX `idx_orders_user_id` (`user_id`);

-- Add user_id to shops table (allow NULL for existing data, then update)
ALTER TABLE `shops` ADD COLUMN `user_id` INT UNSIGNED NOT NULL DEFAULT 1 AFTER `id`;
ALTER TABLE `shops` ADD INDEX `idx_shops_user_id` (`user_id`);

-- Remove unique constraint from shops.name (now unique per user)
-- Note: If you need to keep shop names unique across all users, skip this
-- ALTER TABLE `shops` DROP INDEX `name`;