package mangoplus

type ChapterListGroup struct {
	ChapterNumbers   string    `json:"chapterNumbers"`
	FirstChapterList []Chapter `json:"firstChapterList"`
	MidChapterList   []Chapter `json:"midChapterList"`
	LastChapterList  []Chapter `json:"lastChapterList"`
}

type Chapter struct {
	TitleId        int     `json:"titleId"`
	ChapterId      int     `json:"chapterId"`
	Name           string  `json:"name"`
	SubTitle       *string `json:"subTitle"`
	ThumbnailUrl   string  `json:"thumbnailUrl"`
	StartTimeStamp int     `json:"startTimeStamp"`
	EndTimeStamp   int     `json:"endTimeStamp"`
	AlreadyViewed  bool    `json:"alreadyViewed"`
	ViewCount      int     `json:"viewCount"`
	CommentCount   int     `json:"commentCount"`
	IsVerticalOnly bool    `json:"isVerticalOnly"`
}
