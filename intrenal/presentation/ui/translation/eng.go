package translation

var engTranslations = map[Sentence]string{
	HelpMe: "Help me text......",
}

func (s Sentence) ToEng() string {
	if translation, ok := engTranslations[s]; ok {
		return translation
	}
	return string(s)
}

//TODO : use user-lang to translate to user'langs === maybe need redis user cache for it
