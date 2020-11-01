CREATE TABLE IF NOT EXISTS trip_messages
(
    trip_id    uuid primary key not null,
    message_id bigint           not null
);

-- each message should only be in here once, since a message can't be about more than one trip
CREATE UNIQUE INDEX trip_messages_message_id_unique_idx ON trip_messages (message_id);
