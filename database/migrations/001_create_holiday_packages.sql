-- Holiday Packages System Migration
-- This creates tables for managing holiday packages that use existing flight/helicopter schedules

-- 1. Holiday Packages Table
CREATE TABLE IF NOT EXISTS `holiday_packages` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL,
    `description` text,
    `package_type` enum('spiritual','wildlife','adventure','cultural') NOT NULL DEFAULT 'spiritual',
    `duration_days` int NOT NULL DEFAULT 1,
    `duration_nights` int NOT NULL DEFAULT 0,
    `price_per_person` decimal(10,2) NOT NULL,
    `max_passengers` int NOT NULL DEFAULT 6,
    `status` tinyint NOT NULL DEFAULT 1 COMMENT '1=Active, 0=Inactive',
    `image_url` varchar(500) DEFAULT NULL,
    `inclusions` json DEFAULT NULL COMMENT 'Array of included services',
    `exclusions` json DEFAULT NULL COMMENT 'Array of excluded services',
    `itinerary` json DEFAULT NULL COMMENT 'Day-wise itinerary',
    `terms_conditions` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_package_type` (`package_type`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Package Schedules Table (Links packages to flight/helicopter schedules)
CREATE TABLE IF NOT EXISTS `package_schedules` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `package_id` bigint unsigned NOT NULL,
    `schedule_type` enum('flight','helicopter') NOT NULL,
    `schedule_id` int NOT NULL COMMENT 'References flight_schedules.id or helicopter_schedules.id from Node.js backend',
    `sequence_order` int NOT NULL DEFAULT 1 COMMENT 'Order of this schedule in the package (1st leg, 2nd leg, etc.)',
    `day_number` int NOT NULL DEFAULT 1 COMMENT 'Which day of the package this schedule is on',
    `is_return` tinyint NOT NULL DEFAULT 0 COMMENT '1=Return journey, 0=Onward journey',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_package_id` (`package_id`),
    KEY `idx_schedule_type_id` (`schedule_type`, `schedule_id`),
    KEY `idx_sequence` (`package_id`, `sequence_order`),
    CONSTRAINT `fk_package_schedules_package` FOREIGN KEY (`package_id`) REFERENCES `holiday_packages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Package Bookings Table
CREATE TABLE IF NOT EXISTS `package_bookings` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `package_id` bigint unsigned NOT NULL,
    `booking_reference` varchar(20) NOT NULL UNIQUE,
    `pnr` varchar(10) NOT NULL UNIQUE,
    `guest_name` varchar(255) NOT NULL,
    `guest_email` varchar(255) NOT NULL,
    `guest_phone` varchar(20) NOT NULL,
    `num_passengers` int NOT NULL,
    `travel_date` date NOT NULL COMMENT 'Start date of the package',
    `total_amount` decimal(10,2) NOT NULL,
    `booking_status` enum('pending','confirmed','cancelled','completed') NOT NULL DEFAULT 'pending',
    `payment_status` enum('pending','paid','failed','refunded') NOT NULL DEFAULT 'pending',
    `payment_id` varchar(100) DEFAULT NULL,
    `payment_method` varchar(50) DEFAULT NULL,
    `special_requests` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_booking_reference` (`booking_reference`),
    UNIQUE KEY `uk_pnr` (`pnr`),
    KEY `idx_package_id` (`package_id`),
    KEY `idx_booking_status` (`booking_status`),
    KEY `idx_payment_status` (`payment_status`),
    KEY `idx_travel_date` (`travel_date`),
    CONSTRAINT `fk_package_bookings_package` FOREIGN KEY (`package_id`) REFERENCES `holiday_packages` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Package Passengers Table
CREATE TABLE IF NOT EXISTS `package_passengers` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `booking_id` bigint unsigned NOT NULL,
    `title` enum('Mr','Mrs','Ms','Dr','Master','Miss') NOT NULL,
    `first_name` varchar(100) NOT NULL,
    `last_name` varchar(100) NOT NULL,
    `age` int NOT NULL,
    `gender` enum('Male','Female','Other') NOT NULL,
    `passenger_type` enum('Adult','Child','Infant') NOT NULL DEFAULT 'Adult',
    `id_proof_type` varchar(50) DEFAULT NULL,
    `id_proof_number` varchar(100) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_booking_id` (`booking_id`),
    KEY `idx_passenger_type` (`passenger_type`),
    CONSTRAINT `fk_package_passengers_booking` FOREIGN KEY (`booking_id`) REFERENCES `package_bookings` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Package Schedule Bookings Table (Links package bookings to individual schedule bookings in Node.js backend)
CREATE TABLE IF NOT EXISTS `package_schedule_bookings` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `package_booking_id` bigint unsigned NOT NULL,
    `package_schedule_id` bigint unsigned NOT NULL,
    `node_booking_id` int DEFAULT NULL COMMENT 'References bookings.id or helicopter_bookings.id from Node.js backend',
    `booking_type` enum('flight','helicopter') NOT NULL,
    `booking_date` date NOT NULL,
    `seat_assignments` json DEFAULT NULL COMMENT 'Array of seat assignments for passengers',
    `booking_status` enum('pending','confirmed','cancelled') NOT NULL DEFAULT 'pending',
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_package_booking_id` (`package_booking_id`),
    KEY `idx_package_schedule_id` (`package_schedule_id`),
    KEY `idx_node_booking_id` (`node_booking_id`),
    KEY `idx_booking_type` (`booking_type`),
    CONSTRAINT `fk_package_schedule_bookings_package_booking` FOREIGN KEY (`package_booking_id`) REFERENCES `package_bookings` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_package_schedule_bookings_package_schedule` FOREIGN KEY (`package_schedule_id`) REFERENCES `package_schedules` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample holiday packages data
INSERT INTO `holiday_packages` (`title`, `description`, `package_type`, `duration_days`, `duration_nights`, `price_per_person`, `max_passengers`, `inclusions`, `itinerary`) VALUES
('Maihar VIP Darshan - Jabalpur Return', '🔱 VIP Helicopter Tour to Maa Sharda Devi Temple with same day return from Jabalpur', 'spiritual', 1, 0, 16000.00, 6, 
'["Jabalpur → Maihar → Jabalpur helicopter service", "Fast, safe & time-saving travel", "Maihar helipad to temple AC car transfer (to & fro)", "Maa Sharda Devi VIP Darshan (No queue, special arrangements)", "Special Prasad from temple", "Comfortable, secure & luxury journey", "Professional services & experienced staff"]',
'[{"day": 1, "title": "Jabalpur to Maihar VIP Darshan", "activities": ["Departure from Jabalpur", "Helicopter flight to Maihar", "AC car transfer to temple", "VIP Darshan at Maa Sharda Devi Temple", "Special Prasad collection", "Return helicopter flight to Jabalpur"], "duration": "2 Hr 30 Min"}]'),

('Maihar VIP Darshan - Chitrakoot', '🔱 VIP Helicopter Tour from Chitrakoot to Maa Sharda Temple', 'spiritual', 1, 0, 10000.00, 6,
'["Chitrakoot helicopter service", "Maa Sharda Temple VIP Darshan & special prasad", "Full AC luxury cab (arrival & departure)"]',
'[{"day": 1, "title": "Chitrakoot to Maihar VIP Darshan", "activities": ["Departure from Chitrakoot", "Helicopter flight to Maihar", "AC luxury cab to temple", "VIP Darshan at Maa Sharda Devi Temple", "Special Prasad collection", "Return journey"], "duration": "1 Hr 30 Min"}]'),

('Bandhavgarh Wildlife Safari', '🐅 1 Night / 2 Days Wildlife Helicopter Tour with jungle safari', 'wildlife', 2, 1, 25000.00, 6,
'["Helicopter travel (Jabalpur ⇄ Bandhavgarh)", "1 Night stay in Bandhavgarh", "All meals (Lunch, Dinner & Breakfast)", "Jungle Safari with all necessary permits", "Cab transfers (Helipad ⇄ Resort ⇄ Safari Gate)"]',
'[{"day": 1, "title": "Jabalpur to Bandhavgarh", "activities": ["Helicopter departure from Jabalpur", "Arrival at Bandhavgarh", "Check-in at resort", "Lunch", "Evening jungle safari", "Dinner at resort"], "duration": "Full Day"}, {"day": 2, "title": "Bandhavgarh Safari & Return", "activities": ["Early morning jungle safari", "Breakfast at resort", "Check-out", "Helicopter return to Jabalpur"], "duration": "Half Day"}]'),

('Bandhavgarh & Kanha Wildlife Tour', '🐅 2 Nights / 3 Days Complete Wildlife Experience', 'wildlife', 3, 2, 40000.00, 6,
'["Helicopter travel (Jabalpur → Bandhavgarh → Kanha → Jabalpur)", "2 Nights stay (Bandhavgarh & Kanha)", "All meals (Lunch, Dinner & Breakfast)", "Jungle Safari with permits", "Cab transfers (Helipad ⇄ Resort ⇄ Safari Gate)"]',
'[{"day": 1, "title": "Jabalpur to Bandhavgarh", "activities": ["Helicopter departure from Jabalpur", "Arrival at Bandhavgarh", "Check-in at resort", "Lunch", "Evening jungle safari", "Dinner"], "duration": "Full Day"}, {"day": 2, "title": "Bandhavgarh to Kanha", "activities": ["Morning jungle safari", "Breakfast", "Helicopter transfer to Kanha", "Check-in at Kanha resort", "Lunch", "Evening safari", "Dinner"], "duration": "Full Day"}, {"day": 3, "title": "Kanha Safari & Return", "activities": ["Early morning safari", "Breakfast", "Check-out", "Helicopter return to Jabalpur"], "duration": "Half Day"}]');