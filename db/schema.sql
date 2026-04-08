-- Create "users" table
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "loans" table
CREATE TABLE `loans` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `principal` decimal(10,2) NOT NULL,
  `term` int NOT NULL,
  `interest` float NOT NULL,
  `total_amount` decimal(10,2) NOT NULL,
  `weekly_installment` decimal(10,2) NOT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "loan_billings" table
CREATE TABLE `loan_billings` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `loan_id` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` tinyint NOT NULL,
  `due_date` datetime NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_loan_billings_on_loan_id_and_due_date` (`loan_id`, `due_date`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "repayments" table
CREATE TABLE `repayments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `loan_billing_id` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `status` tinyint NOT NULL,
  `reference` varchar(255) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "accounts" table
CREATE TABLE `accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `balance` decimal(10,2) NOT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_accounts_on_user_id` (`user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "mutations" table
CREATE TABLE `mutations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `type` tinyint NOT NULL,
  `reference` varchar(255) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_mutations_on_account_id_and_type_and_reference` (`account_id`, `type`, `reference`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
