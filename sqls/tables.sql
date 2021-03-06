DROP TABLE IF EXISTS channels CASCADE;
CREATE TABLE channels (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  channel_id TEXT NOT NULL DEFAULT '',
  name TEXT NOT NULL DEFAULT '',
  image_url TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

DROP TABLE IF EXISTS videos CASCADE;
CREATE TABLE videos (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  video_id TEXT NOT NULL DEFAULT '',
  channel_id TEXT NOT NULL DEFAULT '',
  title TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  length_seconds BIGINT NOT NULL DEFAULT 0,
  view_count BIGINT NOT NULL DEFAULT 0,
  average_rating NUMERIC NOT NULL DEFAULT 0.0,
  thumbnail_url TEXT NOT NULL DEFAULT '',
  category TEXT NOT NULL DEFAULT '',
  is_private BOOLEAN NOT NULL DEFAULT FALSE,
  publish_date TIMESTAMPTZ,
  upload_date TIMESTAMPTZ,
  live_started_at TIMESTAMPTZ,
  live_ended_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  chat_id TEXT NOT NULL DEFAULT '',
  author_channel_id TEXT NOT NULL DEFAULT '',
  video_id TEXT NOT NULL DEFAULT '',
  type TEXT NOT NULL DEFAULT '',
  timestamp TEXT NOT NULL DEFAULT '',
  timestamp_usec BIGINT NOT NULL DEFAULT 0,
  message_elements JSONB NOT NULL DEFAULT '[]',
  purchase_amount NUMERIC DEFAULT 0.0,
  currency_unit TEXT NOT NULL DEFAULT '',
  super_chat_context JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

DROP TABLE IF EXISTS badges CASCADE;
CREATE TABLE badges (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  chat_id TEXT NOT NULL DEFAULT '',
  badge_type TEXT NOT NULL DEFAULT '',
  image_url TEXT NOT NULL DEFAULT '',
  label TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
