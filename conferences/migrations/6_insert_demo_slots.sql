INSERT INTO conference_slot (name, description,
    COST, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public, conference_id)
    VALUES ('General Admision - Early Bird', 'Access to community area and talks, does not inlude access to paid workshops.', 30000, 100, '2021-06-09 17:00:00+00', '2021-06-19 17:00:00+00', '2020-12-03 17:00:00+00', '2021-02-03 21:00:00+00', TRUE, 1), ('Pre-Conference Workshop - Looking Grumpy', 'This ticked does not imply access to the conference.', 30000, 100, '2021-06-09 17:00:00+00', '2021-06-19 17:00:00+00', '2020-12-03 17:00:00+00', '2021-02-03 21:00:00+00', TRUE, 1), ('Sponsor Access', 'This ticked Is not for sale.', 30000, 100, '2021-06-09 17:00:00+00', '2021-06-19 17:00:00+00', '2020-12-03 17:00:00+00', '2021-02-03 21:00:00+00', FALSE, 1);

INSERT INTO conference_slot (name, description,
    COST, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public, conference_id, depends_on)
    VALUES ('Conference Free Workshop - Looking Grumpy with an editor', 'This ticked requires general access to the conference.', 30000, 100, '2021-06-09 17:00:00+00', '2021-06-19 17:00:00+00', '2020-12-03 17:00:00+00', '2021-02-03 21:00:00+00', TRUE, 1, (
            SELECT
                id FROM conference_slot
            WHERE
                name = 'General Admision - Early Bird'
                AND conference_id = 1 LIMIT 1));

