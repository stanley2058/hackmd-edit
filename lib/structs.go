package lib

import "net/http"

type Context struct {
	Note              string
	TeamPath          string
	BaseUrl           string
	ApiToken          string
	Editor            string
	Client            *http.Client
	LastStoredContent string
	TokenUsage        int
}

type Note struct {
	Id              string   `json:"id"`
	Title           string   `json:"title"`
	Tags            []string `json:"tags"`
	CreatedAt       int64    `json:"createdAt"`
	TitleUpdatedAt  int64    `json:"titleUpdatedAt"`
	TagsUpdatedAt   int64    `json:"tagsUpdatedAt"`
	PublishType     string   `json:"publishType"`
	PublishedAt     *int64   `json:"publishedAt"`
	Permalink       *string  `json:"permalink"`
	PublishLink     string   `json:"publishLink"`
	ShortId         string   `json:"shortId"`
	Content         string   `json:"content"`
	LastChangedAt   int64    `json:"lastChangedAt"`
	UserPath        *string  `json:"userPath"`
	TeamPath        *string  `json:"teamPath"`
	ReadPermission  string   `json:"readPermission"`
	WritePermission string   `json:"writePermission"`
}
