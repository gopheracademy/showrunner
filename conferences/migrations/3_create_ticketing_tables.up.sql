ALTER TABLE conference_slot ADD COLUMN depends_on INT REFERENCES conference_slot(id);


CREATE TABLE attendee (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL , 
    coc_accepted BOOLEAN NOT NULL
); 

CREATE TABLE slot_claim (
    id SERIAL PRIMARY KEY,
    attendee_id INT NOT NULL REFERENCES attendee(id),
    conference_slot_id INT NOT NULL REFERENCES conference_slot(id) , 
    ticket_id UUID NOT NULL UNIQUE,
    redeemed BOOLEAN NOT NULL
);


CREATE TABLE claim_payment (
    id SERIAL PRIMARY KEY,
    invoice TEXT 
);

CREATE TABLE payment_method_money (
    id SERIAL PRIMARY KEY,
    amount INTEGER NOT NULL,
    ref TEXT NOT NULL
);

CREATE TABLE payment_method_money_to_claim_payment (
    payment_method_money_id INT NOT NULL UNIQUE REFERENCES payment_method_money(id),
    claim_payment_id INT NOT NULL REFERENCES claim_payment(id)
);

CREATE TABLE payment_method_credit_note (
    id SERIAL PRIMARY KEY,
    amount INTEGER NOT NULL,
    detail TEXT NOT NULL
);

CREATE TABLE payment_method_credit_note_to_claim_payment (
    payment_method_credit_note_id INT NOT NULL UNIQUE REFERENCES payment_method_credit_note(id),
    claim_payment_id INT NOT NULL REFERENCES claim_payment(id) 
);

CREATE TABLE payment_method_conference_discount (
    id SERIAL PRIMARY KEY,
    amount INTEGER  NOT NULL,
    detail TEXT NOT NULL
);

CREATE TABLE payment_method_conference_discount_to_claim_payment (
    payment_method_conference_discount_id INT  NOT NULL UNIQUE REFERENCES payment_method_conference_discount(id),
    claim_payment_id INT NOT NULL REFERENCES claim_payment(id)
);

