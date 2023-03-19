CREATE TABLE IF NOT EXISTS messages
(
    user    text,
    time    timestamp,
    message text,
    direction text,
    PRIMARY KEY ( user, time )
) WITH CLUSTERING ORDER BY (time DESC);