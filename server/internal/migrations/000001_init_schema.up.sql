 -- Create players table
CREATE TABLE IF NOT EXISTS players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- games_count INTEGER DEFAULT 0,
    -- wins_count INTEGER DEFAULT 0,
    -- captain_games INTEGER DEFAULT 0,
    -- captain_wins INTEGER DEFAULT 0,
    -- role_stats JSONB DEFAULT '{
    --     "carry": {"games_count": 0, "wins_count": 0},
    --     "mid": {"games_count": 0, "wins_count": 0},
    --     "offlane": {"games_count": 0, "wins_count": 0},
    --     "pos4": {"games_count": 0, "wins_count": 0},
    --     "pos5": {"games_count": 0, "wins_count": 0}
    -- }'::jsonb,
    games_played JSONB DEFAULT '[]'::jsonb
);

-- Create games table
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    winner VARCHAR(10) NOT NULL CHECK (winner IN ('RADIANT', 'DIRE')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create game_players table
CREATE TABLE IF NOT EXISTS game_players (
    game_id UUID REFERENCES games(id),
    player_id UUID REFERENCES players(id),
    team VARCHAR(10) NOT NULL CHECK (team IN ('RADIANT', 'DIRE')),
    role VARCHAR(10) NOT NULL CHECK (role IN ('carry', 'mid', 'offlane', 'pos4', 'pos5')),
    is_captain BOOLEAN NOT NULL DEFAULT false,
    is_winner BOOLEAN NOT NULL DEFAULT false,
    PRIMARY KEY (game_id, player_id)
); 