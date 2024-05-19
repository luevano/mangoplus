package mangoplus

type ImageQuality string

const (
	ImageQualityLow       ImageQuality = "low"
	ImageQualityHigh      ImageQuality = "high"
	ImageQualitySuperHigh ImageQuality = "super_high"
)

type Rating string

const (
	RatingAllAges  Rating = "ALLAGES"
	RatingTeen     Rating = "TEEN"
	RatingTeenPlus Rating = "TEENPLUS"
	RatingMature   Rating = "MATURE"
)

type Language string

const (
	LanguageEnglish      Language = "ENGLISH"
	LanguageSpanish      Language = "SPANISH"
	LanguageFrench       Language = "FRENCH"
	LanguageIndonesian   Language = "INDONESIAN"
	LanguagePortugueseBR Language = "PORTUGUESE_BR"
	LanguageRussian      Language = "RUSSIAN"
	LanguageThai         Language = "THAI"
	LanguageVietnamese   Language = "VIETNAMESE"
	LanguageGerman       Language = "GERMAN"
)

// ToCode: Get the corresponding Language code.
//
// Using MangaDex as reference: https://api.mangadex.org/docs/3-enumerations/#language-codes--localization
func (l Language) ToCode() string {
	switch l {
	case LanguageEnglish:
		return "en"
	case LanguageSpanish:
		// Could also return "es-la", but "es" only should be more common
		return "es"
	case LanguageFrench:
		return "fr"
	case LanguageIndonesian:
		return "id"
	case LanguagePortugueseBR:
		// Same, could be "pt-br"
		return "pt"
	case LanguageRussian:
		return "ru"
	case LanguageThai:
		return "th"
	case LanguageVietnamese:
		return "vi"
	case LanguageGerman:
		return "de"
	default:
		// TODO: add warning
		return "en"
	}
}

// FromLanguageCode: Get the corresponding Language type from a Language code.
//
// Using MangaDex as reference: https://api.mangadex.org/docs/3-enumerations/#language-codes--localization
func FromLanguageCode(code string) Language {
	switch code {
	case "en":
		return LanguageEnglish
	case "es", "es-la":
		return LanguageSpanish
	case "fr":
		return LanguageFrench
	case "id":
		return LanguageIndonesian
	case "pt", "pt-br":
		return LanguagePortugueseBR
	case "ru":
		return LanguageRussian
	case "th":
		return LanguageThai
	case "vi":
		return LanguageVietnamese
	case "de":
		return LanguageGerman
	default:
		// TODO: add warning
		return LanguageEnglish
	}
}

type ReleaseSchedule string

const (
	ReleaseScheduleDisabled   ReleaseSchedule = "DISABLED"
	ReleaseScheduleEveryday   ReleaseSchedule = "EVERYDAY"
	ReleaseScheduleWeekly     ReleaseSchedule = "WEEKLY"
	ReleaseScheduleBiweekly   ReleaseSchedule = "BIWEEKLY"
	ReleaseScheduleMonthly    ReleaseSchedule = "MONTHLY"
	ReleaseScheduleBiMonthly  ReleaseSchedule = "BIMONTHLY"
	ReleaseScheduleTriMonthly ReleaseSchedule = "TRIMONTHLY"
	ReleaseScheduleOther      ReleaseSchedule = "OTHER"
	ReleaseScheduleCompleted  ReleaseSchedule = "COMPLETED"
)

type LabelCode string

const (
	LabelCodeCreators LabelCode = "CREATORS"
	LabelCodeGiga     LabelCode = "GIGA"
	LabelCodeJPlus    LabelCode = "J_PLUS"
	LabelCodeOthers   LabelCode = "OTHERS"
	LabelCodeRevival  LabelCode = "REVIVAL"
	LabelCodeSKJ      LabelCode = "SKJ"
	LabelCodeSQ       LabelCode = "SQ"
	LabelCodeTYJ      LabelCode = "TYJ"
	LabelCodeVJ       LabelCode = "VJ"
	LabelCodeYJ       LabelCode = "YJ"
	LabelCodeWSJ      LabelCode = "WSJ"
)

func (l Label) Long() string {
	switch l.Label {
	case LabelCodeCreators:
		return "MANGA Plus Creators"
	case LabelCodeGiga:
		return "Shounen Jump Giga"
	case LabelCodeJPlus:
		return "Shounen Jump+"
	case LabelCodeOthers:
		return "Others"
	case LabelCodeRevival:
		return "Revival"
	case LabelCodeSKJ:
		return "Saikyou Jump"
	case LabelCodeSQ:
		return "Jump SQ."
	case LabelCodeTYJ:
		return "Tonari no Young Jump"
	case LabelCodeVJ:
		return "V Jump"
	case LabelCodeYJ:
		return "Weekly Young Jump"
	case LabelCodeWSJ:
		return "Weekly Shounen Jump"
	default:
		return ""
	}
}
