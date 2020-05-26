package models

type Chat struct {
	ChannelID       string
	VideoID         string
	Timestamp       string
	TimestampUsec   string
	MessageElements []MessageElement
	PurchaseAmount  float64
	CurrencyUnit    string
	IsModerator     bool
	Badge           Badge
}

type MessageElement struct {
	Type  string
	Text  string
	Label string
	URL   string
}

type Badge struct {
	Label string
	URL   string
}
