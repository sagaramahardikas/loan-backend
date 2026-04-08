-- Modify "accounts" table
ALTER TABLE `accounts` ADD INDEX `idx_accounts_on_user_id` (`user_id`);
-- Modify "loan_billings" table
ALTER TABLE `loan_billings` ADD INDEX `idx_loan_billings_on_loan_id_and_due_date` (`loan_id`, `due_date`);
-- Modify "mutations" table
ALTER TABLE `mutations` ADD INDEX `idx_mutations_on_account_id_and_type_and_reference` (`account_id`, `type`, `reference`);
