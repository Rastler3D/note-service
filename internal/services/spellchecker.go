package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SpellcheckService interface {
	CorrectText(text string) (string, error)
}

type YandexSpellcheckService struct {
	apiURL string
}

type SpellcheckResult struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func NewSpellcheckService(apiURL string) SpellcheckService {
	return &YandexSpellcheckService{apiURL: apiURL}
}

func (s *YandexSpellcheckService) CheckText(text string) ([]SpellcheckResult, error) {
	params := url.Values{}
	params.Add("text", text)

	resp, err := http.Get(fmt.Sprintf("%s?%s", s.apiURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []SpellcheckResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *YandexSpellcheckService) CorrectText(text string) (string, error) {
	results, err := s.CheckText(text)
	if err != nil {
		return "", err
	}

	correctedText := text
	for i := len(results) - 1; i >= 0; i-- {
		result := results[i]
		if len(result.S) > 0 {
			correctedText = strings.Replace(correctedText, result.Word, result.S[0], 1)
		}
	}

	return correctedText, nil
}
