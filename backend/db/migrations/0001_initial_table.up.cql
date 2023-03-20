CREATE TABLE IF NOT EXISTS messages
(
    user    text,
    time    timestamp,
    message text,
    direction text,
    PRIMARY KEY ( user, time )
) WITH
    CLUSTERING ORDER BY (time DESC)
    AND compaction = {
        'class': 'TimeWindowCompactionStrategy',
        'compaction_window_unit': 'DAYS',
        'compaction_window_size': 1,
        'max_threshold': 1000000,
        'expired_sstable_check_frequency_seconds': 86400
    }
    AND default_time_to_live = 31536000;