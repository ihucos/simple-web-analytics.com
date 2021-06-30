package endpoints

import (
	"github.com/ihucos/counter.dev/models"
	"github.com/ihucos/counter.dev/lib"
	"net/http"
)

type UserDump struct {
	Id    string            `json:"id"`
	Token string            `json:"token"`
	Prefs map[string]string `json:"prefs"`
}

type SitesDumpVal struct {
	Count  int                `json:"count"`
	Logs   models.LogData     `json:"logs"`
	Visits models.TimedVisits `json:"visits"`
}

type SitesDump map[string]SitesDumpVal
type Meta map[string]string

type Dump struct {
	Sites SitesDump         `json:"sites"`
	User  UserDump          `json:"user"`
	Meta  map[string]string `json:"meta"`
}

func Origin2SiteId(origin string) string {
	// this function returns
	var re = regexp.MustCompile(`^.*?:\/\/(?:www.)?(.*)$`)
	var match = re.FindStringSubmatch(origin)
	if len(match) < 1 {
		return origin
	}
	return match[1]
}

func LoadSitesDump(user models.User, utcOffset int) (SitesDump, error) {
	sitesDump := make(SitesDump)

	sitesLink, err := user.GetPreferredSiteLinks()
	if err != nil {
		return SitesDump{}, err
	}

	for siteId, count := range sitesLink {
		site := user.NewSite(siteId)
		logs, err := site.GetLogs()
		if err != nil {
			return SitesDump{}, err
		}
		visits, err := site.GetVisits(utcOffset)
		if err != nil {
			return SitesDump{}, err
		}
		sitesDump[siteId] = SitesDumpVal{
			Logs:   logs,
			Visits: visits,
			Count:  count,
		}
	}
	return sitesDump, nil
}

func LoadDump(user models.User, utcOffset int) (Dump, error) {
	prefsData, err := user.GetPrefs()
	if err != nil {
		return Dump{}, err
	}

	token, err := user.ReadToken()
	if err != nil {
		return Dump{}, err
	}

	sitesDump, err := LoadSitesDump(user, utcOffset)
	if err != nil {
		return Dump{}, err
	}

	userDump := UserDump{Id: user.Id, Token: token, Prefs: prefsData}
	return Dump{User: userDump, Sites: sitesDump, Meta: Meta{}}, nil
}
