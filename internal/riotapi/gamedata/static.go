package gamedata

import (
	"encoding/json"
	"os"

	"LoLItemRecommender/internal/levenshtein"
	"LoLItemRecommender/internal/printer"
	"LoLItemRecommender/internal/riotapi/api"
)

type StaticData struct {
	APIVersion     string
	ChampionsStats map[string]*ChampionStats
	ItemsData      map[string]*ItemData
	client         *api.Client
	em             *api.EndpointsManager
	apiKey         string
}

func NewStaticData(em *api.EndpointsManager, client *api.Client) *StaticData {
	return &StaticData{
		apiKey:         os.Getenv("RIOT_API_KEY"),
		em:             em,
		client:         client,
		ChampionsStats: make(map[string]*ChampionStats),
		ItemsData:      make(map[string]*ItemData),
	}
}

const (
	StrongProbability = 70.0
	WeakProbability   = 50.0
)

func (sd *StaticData) GetChampionStats(name string) *ChampionStats {
	for c, s := range sd.ChampionsStats {
		//printer.Debug("Similarity: %s & %s = %f", c, name, levenshtein.StringSimilarity(c, name))
		if levenshtein.StringSimilarity(c, name) >= StrongProbability {
			return s
		}
	}
	return nil
}

func (sd *StaticData) GetChampionsStatsWithCloseName(name string) []*ChampionStats {
	cs := make([]*ChampionStats, 0)
	for c, s := range sd.ChampionsStats {
		//printer.Debug("Similarity: %s & %s = %f", c, name, levenshtein.StringSimilarity(c, name))
		if levenshtein.StringSimilarity(c, name) >= WeakProbability {
			cs = append(cs, s)
		}
	}
	return cs
}

func (sd *StaticData) RetrieveAPIVersion() error {
	body, err := sd.client.Get(api.DDragonStaticVersionsURL)
	if err != nil {
		return err
	}
	var v []string
	if err = json.Unmarshal(body, &v); err != nil {
		return err
	}
	sd.APIVersion = v[0]
	return nil
}

func (sd *StaticData) RetrieveChampionsStats() error {
	body, err := sd.client.Get(sd.em.GetStaticDataChampionsURL(sd.APIVersion))
	if err != nil {
		return err
	}
	var v ChampionsDataResponse
	if err = json.Unmarshal(body, &v); err != nil {
		return err
	}
	for n, c := range v.Data {
		pc := c
		sd.ChampionsStats[n] = &pc
	}
	printer.Printf("{-F_CYAN,BOLD}%d {-RESET}champions found", len(sd.ChampionsStats))
	return nil
}

func (sd *StaticData) RetrieveItems() error {
	body, err := sd.client.Get(sd.em.GetStaticDataItemsURL(sd.APIVersion))
	if err != nil {
		return err
	}
	var i ItemResponse
	if err = json.Unmarshal(body, &i); err != nil {
		return err
	}
	for n, item := range i.Data {
		pc := item
		sd.ItemsData[n] = &pc
	}
	printer.Printf("{-F_CYAN,BOLD}%d {-RESET}items found", len(sd.ChampionsStats))
	return nil
}
