INSERT INTO channels (
  channel_id,
  name,
  image_url,
  created_at,
  updated_at
) VALUES (
  'CHANNEL_A',
  'A channel',
  'https://placehold.jp/64x64.png?text=dummy',
  NOW(),
  NOW()
), (
  'CHANNEL_B',
  'B channel',
  'https://placehold.jp/64x64.png?text=dummy',
  NOW(),
  NOW()
);

INSERT INTO videos (
  video_id,
  channel_id,
  created_at,
  updated_at
) VALUES (
  'VIDEO_A',
  'CHANNEL_A',
  NOW(),
  NOW()
), (
  'VIDEO_B',
  'CHANNEL_B',
  NOW(),
  NOW()
);

INSERT INTO chats (
  author_channel_id,
  chat_id,
  video_id,
  type,
  timestamp,
  timestamp_usec,
  message_elements,
  purchase_amount,
  currency_unit,
  super_chat_context,
  created_at,
  updated_at
) VALUES (
  'CHANNEL_A',
  'CHAT_A',
  'VIDEO_B',
  'chat',
  '12:34',
  '1234567890123456',
  '[{"type":"text","text":"hello"}]',
  0,
  '',
  '{}',
  NOW(),
  NOW()
), (
  'CHANNEL_B',
  'CHAT_B',
  'VIDEO_A',
  'super_chat',
  '43:21',
  '1234567890123456',
  '[{"type":"text","text":"hello"},{"type":"emoji","image_url":"https://placehold.jp/64x64.png?text=dummy","label":":dummy:"}]',
  1000,
  '¥',
  '{"header_background_color":"ffffffff","header_text_color":"ff000000","body_background_color":"ffeeeeee","body_text_color":"ff333333","author_text_color":"ff000000"}',
  NOW(),
  NOW()
);

INSERT INTO badges (
	chat_id,
	badge_type,
	image_url,
	label,
  created_at,
  updated_at
) VALUES (
  'CHAT_A',
  'member',
  'https://placehold.jp/64x64.png?text=dummy',
  'dummy member badge',
  NOW(),
  NOW()
), (
  'CHAT_B',
  'moderator',
  '',
  '',
  NOW(),
  NOW()
);
