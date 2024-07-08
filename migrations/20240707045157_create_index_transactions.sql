-- Create index "transactions_timestamp_idx" to table: "transactions"
CREATE INDEX "transactions_timestamp_idx" ON "transactions" ("timestamp" DESC);
-- Create index "transactions_user_id_idx" to table: "transactions"
CREATE INDEX "transactions_user_id_idx" ON "transactions" ("user_id");
