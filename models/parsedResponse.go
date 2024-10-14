package models

type ParsedResponse struct {
	Word     string
	Pos      string
	Meaning  string
	Example  interface{}
	Synonyms interface{}
	Antonyms interface{}
}
