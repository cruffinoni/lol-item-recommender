package api

import (
	"fmt"
)

const (
	// Summoner endpoints
	summonerBySummonerNameEndpoint = "/summoner/v4/summoners/by-name/%s"
	summonersByLeague              = "/league-exp/v4/entries/%s/%s/%s"

	// Match endpoints
	playerMatchListEndpoint = "/match/v5/matches/by-puuid/%s/ids?start=0&count=100" // TODO: Change count to be a parameter
	playerMatchInfoEndpoint = "/match/v5/matches/%s"

	// Static data endpoints
	staticDataChampionsEndpoint = "/cdn/%s/data/en_US/champion.json"
	staticDataItemEndpoint      = "/cdn/%s/data/en_US/item.json"
	staticDDragonBaseURL        = "https://ddragon.leagueoflegends.com"

	DDragonStaticVersionsURL = staticDDragonBaseURL + "/api/versions.json"
)

type EndpointsManager struct {
	Apikey           string
	Region           string
	apiBaseURL       string
	regionApiBaseURL string
}

func NewEndpointsManager(apikey, region string) *EndpointsManager {
	return &EndpointsManager{
		Apikey:           apikey,
		Region:           region,
		apiBaseURL:       "https://" + region + ".api.riotgames.com/lol",
		regionApiBaseURL: "https://europe.api.riotgames.com/lol",
	}
}

// GetSummonerByName returns the URL for the summoner by name endpoint.
func (em *EndpointsManager) GetSummonerByName(summonerName string) string {
	return fmt.Sprintf(em.apiBaseURL+summonerBySummonerNameEndpoint+"?api_key="+em.Apikey, summonerName)
}

func (em *EndpointsManager) GetSummonersByLeague(queue, division, tier string, page int) string {
	return fmt.Sprintf(em.apiBaseURL+summonersByLeague+"?page=%d&api_key=%s", queue, division, tier, page, em.Apikey)
}

// GetMatchInfoURL returns the URL for the match info endpoint.
func (em *EndpointsManager) GetMatchInfoURL(matchID string) string {
	return fmt.Sprintf(em.regionApiBaseURL+playerMatchInfoEndpoint+"?api_key=%s", matchID, em.Apikey)
}

func (em *EndpointsManager) GetMatchListURL(pUUID string) string {
	return fmt.Sprintf(em.regionApiBaseURL+playerMatchListEndpoint+"&api_key=%s", pUUID, em.Apikey)
}

// GetStaticDataChampionsURL returns the URL for the static data champions endpoint.
func (em *EndpointsManager) GetStaticDataChampionsURL(version string) string {
	return fmt.Sprintf(staticDDragonBaseURL+staticDataChampionsEndpoint, version)
}

func (em *EndpointsManager) GetStaticDataItemsURL(version string) string {
	return fmt.Sprintf(staticDDragonBaseURL+staticDataItemEndpoint, version)
}
