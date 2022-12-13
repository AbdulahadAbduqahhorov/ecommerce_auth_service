CREATE TYPE "user_type_t" AS ENUM ('SUPERADMIN', 'ADMIN', 'USER');

CREATE TABLE IF NOT EXISTS "user" (
	"id" CHAR(36) PRIMARY KEY,
	"full_name" VARCHAR(255)  NOT NULL,
	"login" VARCHAR(255) UNIQUE NOT NULL,
	"phone" VARCHAR UNIQUE,
    "email" VARCHAR UNIQUE,
	"password" VARCHAR(255) NOT NULL ,
	"user_type" user_type_t NOT NULL ,
	"created_at" TIMESTAMP DEFAULT now() NOT NULL,
	"updated_at" TIMESTAMP
);

  
