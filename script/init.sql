CREATE TABLE "users" (
 "id" character varying(36) NOT NULL,
 "name" character varying(50) NOT NULL,
 "created_time" bigint NOT NULL,
 "updated_time" bigint NOT NULL,
 PRIMARY KEY ("id")
);

CREATE TABLE "wallets" (
 "address" character varying(36) NOT NULL,
 "user_id" character varying(36) NOT NULL,
 "balance" double precision NOT NULL,
 "created_time" bigint NOT NULL,
 "updated_time" bigint NOT NULL,
 PRIMARY KEY ("address")
);

CREATE INDEX wallets_user_id_idx ON wallets(user_id);

CREATE TABLE "transactions" (
 "tx" character varying(36) NOT NULL,
 "user_id" character varying(36) NOT NULL,
 "from_address" character varying(36) NOT NULL,
 "to_address" character varying(36) NOT NULL,
 "amount" double precision NOT NULL,
 "nonce" character varying(36) NOT NULL,
 "timestamp" bigint NOT NULL,
 PRIMARY KEY ("tx")
);

CREATE INDEX transactions_user_id_idx ON transactions(user_id);
CREATE INDEX transactions_timestamp_idx ON transactions(timestamp DESC);