DROP INDEX IF EXISTS uq_channels_channel_id CASCADE;
DROP INDEX IF EXISTS uq_videos_video_id CASCADE;
DROP INDEX IF EXISTS uq_chats_author_channel_id_video_id_timestamp_usec CASCADE;
DROP INDEX IF EXISTS uq_badges_chat_id CASCADE;
DROP INDEX IF EXISTS idx_videos_channel_id CASCADE;
DROP INDEX IF EXISTS idx_chats_author_channel_id CASCADE;
DROP INDEX IF EXISTS idx_chats_video_id CASCADE;
DROP INDEX IF EXISTS idx_chats_created_at CASCADE;
DROP INDEX IF EXISTS idx_badges_chat_id_at CASCADE;

CREATE UNIQUE INDEX uq_channels_channel_id ON channels(channel_id);
CREATE UNIQUE INDEX uq_videos_video_id ON videos(video_id);
CREATE UNIQUE INDEX uq_chats_author_channel_id_video_id_timestamp_usec ON chats(author_channel_id, video_id, timestamp_usec);
CREATE UNIQUE INDEX uq_badges_chat_id ON badges(chat_id);
CREATE INDEX idx_videos_channel_id ON videos(channel_id);
CREATE INDEX idx_chats_author_channel_id ON chats(author_channel_id);
CREATE INDEX idx_chats_video_id ON chats(video_id);
CREATE INDEX idx_chats_created_at ON chats(created_at);
CREATE INDEX idx_badges_chat_id_at ON chats(chat_id);
