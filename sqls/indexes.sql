DROP INDEX IF EXISTS uq_channels_channel_id CASCADE;
DROP INDEX IF EXISTS uq_videos_video_id CASCADE;
DROP INDEX IF EXISTS idx_videos_channel_id CASCADE;
DROP INDEX IF EXISTS idx_chats_video_id CASCADE;
DROP INDEX IF EXISTS idx_chats_channel_id CASCADE;

CREATE UNIQUE INDEX uq_channels_channel_id ON channels(channel_id);
CREATE UNIQUE INDEX uq_videos_video_id ON videos(video_id);
CREATE INDEX idx_videos_channel_id ON videos(channel_id);
CREATE INDEX idx_chats_video_id ON chats(video_id);
CREATE INDEX idx_chats_channel_id ON chats(channel_id);
