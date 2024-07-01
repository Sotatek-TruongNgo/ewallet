-- Create "transactions" table
CREATE TABLE "transactions" ("tx" character varying(36) NOT NULL, "from" character varying(36) NOT NULL, "to" character varying(36) NOT NULL, "amount" double precision NOT NULL, "nonce" character varying(36) NOT NULL, "timestamp" bigint NOT NULL, PRIMARY KEY ("tx"));
-- Create "users" table
CREATE TABLE "users" ("id" character varying(36) NOT NULL, "name" character varying(50) NOT NULL, "created_time" bigint NOT NULL, "updated_time" bigint NOT NULL, PRIMARY KEY ("id"));
-- Create "wallets" table
CREATE TABLE "wallets" ("address" character varying(36) NOT NULL, "user_id" character varying(36) NOT NULL, "balance" double precision NOT NULL, "created_time" bigint NOT NULL, "updated_time" bigint NOT NULL, PRIMARY KEY ("address"));
