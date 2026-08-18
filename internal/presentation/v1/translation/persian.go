package translation

var persianTranslations = map[Sentence]string{
	HelpMe:       "چطوری ایرانی ؟ خوش میگذره زیر بمب و موشک و جنگنده و گرونی و فساد و بیکاری و بی برقی و محیط زیست  به تاراج رفتو و مردمان روحی ضعیفو و فقر و کوفت و زهرمار و...",
	StartUnknown: "سلام و صد تا سلام به کاربر جدید",
}

func (s Sentence) ToPersian() string {
	if translation, ok := persianTranslations[s]; ok {
		return translation
	}
	return string(s)
}
