package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
	_ "github.com/LassiHeikkila/flmnchll/content/contentdb"
)

// supported search parameters:
// - name (would be nice to support fuzzy search but out of scope now)
// - attribution (would be nice to support fuzzy search but out of scope now)
// - duration (longer or shorter than)
// - category (exact match)
// - upload date (before or after)

const (
	searchKeyName         = "name"
	searchKeyAttribution  = "attr"
	searchKeyDuration     = "dur"
	searchKeyCategory     = "cat"
	searchKeyUploadedDate = "uploaded"

	searchParamGreaterThan = ">"
	searchParamLessThan    = "<"
)

func VideoSearchHandler(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()

	var searchOpts []contentdb.SearchOption

	if values.Has(searchKeyName) {
		param := values.Get(searchKeyName)
		searchOpts = append(searchOpts, contentdb.SearchVideoByName(param))
	}
	if values.Has(searchKeyAttribution) {
		param := values.Get(searchKeyAttribution)
		searchOpts = append(searchOpts, contentdb.SearchVideoByName(param))
	}
	if values.Has(searchKeyDuration) {
		param := values.Get(searchKeyDuration)
		if strings.HasPrefix(param, searchParamGreaterThan) {
			param = strings.TrimPrefix(param, searchParamGreaterThan)
			d, err := time.ParseDuration(param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			searchOpts = append(searchOpts, contentdb.SearchVideoByDuration(d, false))
		} else if strings.HasPrefix(param, searchParamLessThan) {
			param = strings.TrimPrefix(param, searchParamLessThan)
			d, err := time.ParseDuration(param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			searchOpts = append(searchOpts, contentdb.SearchVideoByDuration(d, true))
		}
		// TODO: default case?
	}
	if values.Has(searchKeyCategory) {
		param := values.Get(searchKeyCategory)
		searchOpts = append(searchOpts, contentdb.SearchVideoByCategory(param))
	}
	if values.Has(searchKeyUploadedDate) {
		param := values.Get(searchKeyUploadedDate)
		if strings.HasPrefix(param, searchParamGreaterThan) {
			param = strings.TrimPrefix(param, searchParamGreaterThan)
			t, err := time.Parse(time.RFC3339, param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			searchOpts = append(searchOpts, contentdb.SearchVideoByUploadedBeforeOrAfterDate(t, false))
		} else if strings.HasPrefix(param, searchParamLessThan) {
			param = strings.TrimPrefix(param, searchParamLessThan)
			t, err := time.Parse(time.RFC3339, param)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			searchOpts = append(searchOpts, contentdb.SearchVideoByUploadedBeforeOrAfterDate(t, true))
		}
		// TODO: default case?
	}

	videos, err := contentdb.SearchVideos(searchOpts...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	_ = enc.Encode(videos)
	// TODO: can we do something about possible error?
	// possible actions:
	// - logging
	// - ???
}
