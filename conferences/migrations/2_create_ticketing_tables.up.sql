CREATE TABLE slot_claim (
    id BIGSERIAL PRIMARY KEY,
    conference_slot_id BIGINT, 
    ticket_id VARCHAR(100) CONSTrAINT ticket_id_is_unique UNIQUE,
    redeemed BOOLEAN,
    FOREIGN KEY(conference_slot_id) REFERENCES conference_slot(id)
);

CREATE TABLE attendee (
    id BIGSERIAL PRIMARY KEY,
    email BIGINT, 
    coc_accepted BOOLEAN
);

CREATE TABLE attendee_to_slot_claims (
    attendee_id BIGINT,
    slot_claim_id BIGINT CONSTRAINT slot_claim_id_is_unique UNIQUE,
    FOREIGN KEY (attendee_id) REFERENCES attendee(id),
    FOREIGN KEY (slot_claim_id) REFERENCES slot_claim(id)
);

CREATE TABLE claim_payment (
    id BIGSERIAL PRIMARY KEY,
    invoice TEXT -- just in case we need to store the whole thing.
);

CREATE TABLE payment_method_money (
    id BIGSERIAL PRIMARY KEY,
    amount INTEGER,
    ref VARCHAR(250)
);

CREATE TABLE payment_method_money_to_claim_payment (
    payment_method_money_id BIGINT,
    claim_payment_id BIGINT,
    FOREIGN KEY (payment_method_money_id) REFERENCES payment_method_money(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_credit_note (
    id BIGSERIAL PRIMARY KEY,
    amount INTEGER,
    detail VARCHAR(250)
);

CREATE TABLE payment_method_credit_note_to_claim_payment (
    payment_method_credit_note_id BIGINT,
    claim_payment_id BIGINT,
    FOREIGN KEY (payment_method_credit_note_id) REFERENCES payment_method_credit_note(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_conference_discount (
    id BIGSERIAL PRIMARY KEY,
    amount INTEGER,
    detail VARCHAR(250)
);

CREATE TABLE payment_method_conference_discount_to_claim_payment (
    payment_method_conference_discount_id BIGINT,
    claim_payment_id BIGINT,
    FOREIGN KEY (payment_method_conference_discount_id) REFERENCES payment_method_conference_discount(id),
    FOREIGN KEY (claim_payment_id) REFERENCES claim_payment(id)
);

