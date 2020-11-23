CREATE TABLE slot_claim (
    id SERIAL PRIMARY KEY,
    conference_slot_id INT REFERENCES conference_slot(id), 
    ticket_id UUID CONSTRAINT ticket_id_is_unique UNIQUE,
    redeemed BOOLEAN
);

CREATE TABLE attendee (
    id SERIAL PRIMARY KEY,
    email TEXT, 
    coc_accepted BOOLEAN
);

CREATE TABLE attendee_to_slot_claims (
    attendee_id INT REFERENCES attendee(id),
    slot_claim_id INT CONSTRAINT slot_claim_id_is_unique UNIQUE,
    FOREIGN KEY (slot_claim_id) REFERENCES slot_claim(id)
);

CREATE TABLE claim_payment (
    id SERIAL PRIMARY KEY,
    invoice TEXT 
);

CREATE TABLE payment_method_money (
    id SERIAL PRIMARY KEY,
    amount INTEGER,
    ref TEXT
);

CREATE TABLE payment_method_money_to_claim_payment (
    payment_method_money_id INT UNIQUE REFERENCES payment_method_money(id),
    claim_payment_id INT REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_credit_note (
    id SERIAL PRIMARY KEY,
    amount INTEGER,
    detail TEXT
);

CREATE TABLE payment_method_credit_note_to_claim_payment (
    payment_method_credit_note_id INT UNIQUE REFERENCES payment_method_credit_note(id),
    claim_payment_id INTREFERENCES claim_payment(id) 
);

CREATE TABLE payment_method_conference_discount (
    id SERIAL PRIMARY KEY,
    amount INTEGER,
    detail TEXT
);

CREATE TABLE payment_method_conference_discount_to_claim_payment (
    payment_method_conference_discount_id INT UNIQUE REFERENCES payment_method_conference_discount(id),
    claim_payment_id INT REFERENCES claim_payment(id)
);

