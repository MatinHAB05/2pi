package translation

type Sentence string

// TODO : add func Translate(enum string language,...args) string
const (
	HelpMe          Sentence = "/help"
	StartUnknown    Sentence = "/start-unknown"
	StartRegistered Sentence = "/start-registered"
)
