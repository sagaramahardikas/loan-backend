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
