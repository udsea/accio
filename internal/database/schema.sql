-- Schema for Accio Profiles Database

-- Profiles table stores information about social media profiles
CREATE TABLE IF NOT EXISTS profiles (
    id INTEGER PRIMARY KEY,
    real_name TEXT NOT NULL,
    username TEXT NOT NULL,
    platform TEXT NOT NULL,
    profile_url TEXT NOT NULL,
    image_url TEXT,
    verified BOOLEAN DEFAULT FALSE,
    follower_count INTEGER,
    bio TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(username, platform)
);

-- Create indices for faster lookups
CREATE INDEX IF NOT EXISTS idx_real_name ON profiles(real_name);
CREATE INDEX IF NOT EXISTS idx_username ON profiles(username);
CREATE INDEX IF NOT EXISTS idx_platform ON profiles(platform);

-- Name parts table for better name matching
CREATE TABLE IF NOT EXISTS name_parts (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    name_part TEXT NOT NULL,
    part_type TEXT NOT NULL, -- 'first', 'middle', 'last', 'nickname'
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_name_part ON name_parts(name_part);

-- Known aliases table
CREATE TABLE IF NOT EXISTS aliases (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    alias TEXT NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_alias ON aliases(alias);

-- Platform-specific data
CREATE TABLE IF NOT EXISTS platform_data (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    data_key TEXT NOT NULL,
    data_value TEXT NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    UNIQUE(profile_id, data_key)
);

-- Search history for improving results
CREATE TABLE IF NOT EXISTS search_history (
    id INTEGER PRIMARY KEY,
    query TEXT NOT NULL,
    result_count INTEGER NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User feedback for improving accuracy
CREATE TABLE IF NOT EXISTS user_feedback (
    id INTEGER PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    feedback_type TEXT NOT NULL, -- 'correct', 'incorrect', 'missing'
    comment TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);