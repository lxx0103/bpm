# Getting Started
1.	make a config file

        cp config.toml.example config.toml

2.	database init

        CREATE TABLE `user_auths` (
            `id` int NOT NULL AUTO_INCREMENT,
            `user_id` int NOT NULL DEFAULT '0',
            `auth_type` tinyint NOT NULL,
            `identifier` varchar(255) NOT NULL,
            `credential` varchar(255) NOT NULL,
            `created_at` datetime NOT NULL,
            `updated_at` datetime NOT NULL,
            PRIMARY KEY (`id`),
            UNIQUE KEY `login` (`auth_type`,`identifier`) USING BTREE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 

        CREATE TABLE `users` (
            `id` int NOT NULL AUTO_INCREMENT,
            `name` varchar(64) NOT NULL,
            `email` varchar(255) DEFAULT NULL,
            `role_id` int DEFAULT '0',
            `gender` tinyint DEFAULT '0',
            `created_at` datetime NOT NULL,
            `updated_at` datetime NOT NULL,
            PRIMARY KEY (`id`),
            UNIQUE KEY `email` (`email`)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 

3.	Run

        go run main.go

