-- CREATE DATABASE netwatch;

CREATE TABLE IF NOT EXISTS crawl_index (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    links TEXT[],
    content JSONB,
    html TEXT,
    headers JSONB,
    timestamp TIMESTAMPTZ DEFAULT NOW() 
);

-- Create an index on the url column for faster lookups
CREATE INDEX IF NOT EXISTS idx_crawl_index_url ON crawl_index (url);