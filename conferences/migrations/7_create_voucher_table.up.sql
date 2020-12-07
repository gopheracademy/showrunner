CREATE TABLE discount_vouchers (
    voucher_id uuid PRIMARY KEY,
    valid_from timestamptz NOT NULL,
    valid_to timestamptz NOT NULL,
    discount_percentage integer NOT NULL DEFAULT 0,
    discount_percentage_max_amount_cents integer NOT NULL DEFAULT 0,
    discount_amount_cents integer NOT NULL DEFAULT 0,
    spent boolean NOT NULL DEFAULT FALSE,
    conference_id integer NOT NULL REFERENCES conference(id),
    backing_payment_id integer REFERENCES payment_method_conference_discount (id)
);

