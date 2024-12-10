package abbreviation

import (
	"encoding/json"
	//"log"
	"net/http"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type abbProviderCtx struct {
	ListID      string `json:"listId"`
	CurrentPage int    `json:"currentPage"`
	PerPage     int    `json:"perPage"`
	Filter      string `json:"filter"`
	SortBy      string `json:"sortBy"`
	SortDesc    bool   `json:"sortDesc"`
}

func fuzzyFind(filter string, userList []*Abbreviation) []*Abbreviation {
	var result []*Abbreviation

	for _, abb := range userList {
		if fuzzy.Match(filter, abb.Abb) {
			result = append(result, abb)
			continue
		}
		if fuzzy.Match(filter, abb.Word) {
			result = append(result, abb)
		}
	}
	return result

}

func paginate(abbs []*Abbreviation, currPage, perPage int) []*Abbreviation {
	skip := currPage * perPage
	if skip > len(abbs) {
		skip = len(abbs)
	}

	end := skip + perPage
	if end > len(abbs) {
		end = len(abbs)
	}

	return abbs[skip:end]
}

func (h *abbHandler) FilterAbbs(w http.ResponseWriter, r *http.Request) {
	var ctx abbProviderCtx

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	abbs := []*Abbreviation{}
	listAbbs, err := h.abbService.GetAbbs(ctx.ListID)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("couldn't get abbs for list" + ctx.ListID))
		return
	}
	//log.Println(ctx)

	if ctx.Filter != "" {
		if strings.HasSuffix(ctx.Filter, "**") {
		  abbs = fuzzyFind(ctx.Filter[:len(ctx.Filter)-2], listAbbs)
		} else {
		  for _, abb := range listAbbs {
			if strings.HasPrefix(abb.Word, ctx.Filter) {
			  abbs = append(abbs, abb)
			}
		  }
		}
	} else {
		abbs = listAbbs
	}

	switch ctx.SortBy {
	case "abb":
		sort.Slice(abbs, func(i, j int) bool {
			if ctx.SortDesc {
				return abbs[i].Abb < abbs[j].Abb
			}
			return abbs[i].Abb > abbs[j].Abb
		})
		break
	case "word":
		sort.Slice(abbs, func(i, j int) bool {
			if ctx.SortDesc {
				return abbs[i].Word < abbs[j].Word
			}
			return abbs[i].Word > abbs[j].Word
		})
		break
	case "updated":
		sort.Slice(abbs, func(i, j int) bool {
			if ctx.SortDesc {
				return abbs[i].Updated.Before(abbs[j].Updated)
			}
			return abbs[i].Updated.After(abbs[j].Updated)
		})
		break
	}

	//log.Println(ctx.Filter)


	result := struct {
		Rows int             `json:"rows"`
		Abbs []*Abbreviation `json:"abbs"`
	}{
		Rows: len(abbs),
		Abbs: abbs,
	}

	result.Abbs = paginate(abbs, ctx.CurrentPage-1, ctx.PerPage)

	if abbs == nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("[]"))
		return

	}
	resp, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't marshal abbs for list " + ctx.ListID))
		return
	}

	w.Write(resp)
}
