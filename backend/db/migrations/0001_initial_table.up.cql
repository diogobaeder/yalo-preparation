CREATE TABLE IF NOT EXISTS yalo.messages
(
    user    text,
    time    timestamp,
    message text,
    PRIMARY KEY ( user, time )
) WITH CLUSTERING ORDER BY (time DESC);