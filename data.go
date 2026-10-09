package mangoplus

type ImageQuality string

const (
	ImageQualityLow       ImageQuality = "low"
	ImageQualityHigh      ImageQuality = "high"
	ImageQualitySuperHigh ImageQuality = "super_high"
)

// StringToImageQuality: Get the corresponding ImageQuality given a string.
//
// Defaults to ImageQualitySuperHigh.
func StringToImageQuality(str string) ImageQuality {
	switch str {
	case "low":
		return ImageQualityLow
	case "high":
		return ImageQualityHigh
	case "super_high":
		return ImageQualitySuperHigh
	default:
		return ImageQualitySuperHigh
	}
}

// ToCode: Get the corresponding Language code.
//
// Defaults to "en".
//
// Using MangaDex as reference: https://api.mangadex.org/docs/3-enumerations/#language-codes--localization
func (l Language) ToCode() string {
	switch l {
	case Language_ENGLISH:
		return "en"
	case Language_SPANISH:
		// Could also return "es-la", but "es" only should be more common
		return "es"
	case Language_FRENCH:
		return "fr"
	case Language_INDONESIAN:
		return "id"
	case Language_PORTUGUESE_BR:
		// Same, could be "pt-br"
		return "pt"
	case Language_RUSSIAN:
		return "ru"
	case Language_THAI:
		return "th"
	case Language_VIETNAMESE:
		return "vi"
	case Language_GERMAN:
		return "de"
	case Language_ITALIAN:
		return "it"
	default:
		return "en"
	}
}

// APICode: Get the language code used in the API params.
//
// Defaults to "eng".
func (l Language) APICode() string {
	switch l {
	case Language_SPANISH:
		return "esp"
	case Language_FRENCH:
		return "fra"
	case Language_INDONESIAN:
		return "ind"
	case Language_PORTUGUESE_BR:
		return "ptb"
	case Language_RUSSIAN:
		return "rus"
	case Language_THAI:
		return "tha"
	case Language_VIETNAMESE:
		return "vie"
	case Language_GERMAN:
		return "deu"
	default:
		return "eng"
	}
}

// StringToLanguage: Get the corresponding Language type from a Language code.
//
// Defaults to Language_ENGLISH.
//
// Using MangaDex as reference: https://api.mangadex.org/docs/3-enumerations/#language-codes--localization
func StringToLanguage(languageCode string) Language {
	switch languageCode {
	case "en":
		return Language_ENGLISH
	case "es", "es-la":
		return Language_SPANISH
	case "fr":
		return Language_FRENCH
	case "id":
		return Language_INDONESIAN
	case "pt", "pt-br":
		return Language_PORTUGUESE_BR
	case "ru":
		return Language_RUSSIAN
	case "th":
		return Language_THAI
	case "vi":
		return Language_VIETNAMESE
	case "de":
		return Language_GERMAN
	case "it":
		return Language_ITALIAN
	default:
		return Language_ENGLISH
	}
}

// Long: Get the full name of the magazine.
func (l LabelCode) Long() string {
	switch l {
	case LabelCode_CREATORS:
		return "MANGA Plus Creators"
	case LabelCode_GIGA:
		return "Shounen Jump Giga"
	case LabelCode_J_PLUS:
		return "Shounen Jump+"
	case LabelCode_OTHERS:
		return "Others"
	case LabelCode_REVIVAL:
		return "Revival"
	case LabelCode_SKJ:
		return "Saikyou Jump"
	case LabelCode_SQ:
		return "Jump SQ."
	case LabelCode_TYJ:
		return "Tonari no Young Jump"
	case LabelCode_VJ:
		return "V Jump"
	case LabelCode_YJ:
		return "Weekly Young Jump"
	case LabelCode_WJ:
		return "Weekly Shounen Jump"
	default:
		return l.String()
	}
}
