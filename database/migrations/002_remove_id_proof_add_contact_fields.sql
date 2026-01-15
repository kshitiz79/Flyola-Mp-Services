-- Migration: Remove ID proof fields and add contact fields to passengers table
-- Date: 2026-01-14
-- Description: Remove id_proof_type and id_proof_number fields, add email and phone fields for primary passenger

-- Remove ID proof fields from package_passengers table
ALTER TABLE `package_passengers` 
DROP COLUMN `id_proof_type`,
DROP COLUMN `id_proof_number`;

-- Add contact fields for primary passenger
ALTER TABLE `package_passengers` 
ADD COLUMN `email` varchar(255) DEFAULT NULL COMMENT 'Email for primary passenger (contact person)',
ADD COLUMN `phone` varchar(20) DEFAULT NULL COMMENT 'Phone for primary passenger (contact person)',
ADD COLUMN `is_primary` tinyint NOT NULL DEFAULT 0 COMMENT '1=Primary passenger (contact person), 0=Additional passenger';

-- Add index for primary passenger lookup
ALTER TABLE `package_passengers` 
ADD KEY `idx_is_primary` (`booking_id`, `is_primary`);

-- Update existing records to mark first passenger as primary
UPDATE `package_passengers` p1
SET `is_primary` = 1
WHERE `id` = (
    SELECT MIN(`id`) 
    FROM (SELECT `id`, `booking_id` FROM `package_passengers`) p2 
    WHERE p2.`booking_id` = p1.`booking_id`
);