ALTER TABLE discount_vouchers ADD COLUMN transaction_id TEXT NOT NULL DEFAULT '';
ALTER TABLE slot_claim ADD COLUMN payment_session TEXT NOT NULL DEFAULT '';