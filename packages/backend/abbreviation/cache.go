package abbreviation

import (
	"fmt"
	"sort"
	"strings"
)

type abbCache struct {
	UserAbbLists    map[string][]string
	UserAbbs        map[string]map[string]*Abbreviation
	UserLookup      map[string]map[string][]*Abbreviation
	UserSuggestions map[string]map[string]int
	UserSharedList  map[string]string
	SharedAbbs      map[string]map[string]string

	GlobalAbbs   map[string][]*Abbreviation
	GlobalLookup map[string][]*Abbreviation
}

func (s *abbService) InitCache(userIDs []string) {
	s.cache.UserAbbs = make(map[string]map[string]*Abbreviation)
	s.cache.UserLookup = make(map[string]map[string][]*Abbreviation)
	s.cache.UserAbbLists = make(map[string][]string)
	s.cache.UserSuggestions = make(map[string]map[string]int)
	s.cache.UserSharedList = make(map[string]string)
	s.cache.SharedAbbs = make(map[string]map[string]string)

	for _, u := range userIDs {
		s.cache.UserSuggestions[u] = make(map[string]int)
		s.cache.UserAbbs[u] = make(map[string]*Abbreviation)
		s.cache.UserLookup[u] = make(map[string][]*Abbreviation)
	}
}

func (s *abbService) Cache(q query) error {
	if s.cache.UserSuggestions[q.UserID] == nil {
		s.cache.UserSuggestions[q.UserID] = make(map[string]int)
	}
	cache := make(map[string]*Abbreviation)
	lookup := make(map[string][]*Abbreviation)

	listIDs := []string{q.Standard}
	listIDs = append(listIDs, q.Addon...)
	lists, err := s.repo.GetLists(listIDs)
	if err != nil {
		return err
	}

	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Type < lists[j].Type
	})

	for _, l := range lists {
		abbs, err := s.repo.GetAbbs(l.ID)

		if err != nil {
			return fmt.Errorf("service|Cache failed getting abbs: %q", err)
		}

		//	log.Printf("service|Cache - caching %s with %d abbs\n", l.ID, len(abbs))
		for _, a := range abbs {
			cache[a.Abb] = a
			lowercaseWord := strings.ToLower(a.Word)
			lookup[lowercaseWord] = append(lookup[lowercaseWord], a)
		}
	}

	s.cache.UserAbbs[q.UserID] = cache
	s.cache.UserLookup[q.UserID] = lookup
	//	log.Printf("Cached %d abbs for user %s, queried with: %q\n", len(s.cache.UserAbbs[q.UserID]), q.UserID, q.ListIDs)
	s.cache.UserAbbLists[q.UserID] = listIDs
	return nil
}

func (s *abbService) UpdateCache(userID string) error {
	cache := make(map[string]*Abbreviation)
	lookup := make(map[string][]*Abbreviation)

	listIDs := s.cache.UserAbbLists[userID]
	lists, err := s.repo.GetLists(listIDs)
	if err != nil {
		return err
	}

	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Type < lists[j].Type
	})

	for _, l := range lists {
		abbs, err := s.repo.GetAbbs(l.ID)

		if err != nil {
			return err
		}

		//log.Printf("service|UpdateCache - caching %s with %d abbs\n", l.ID, len(abbs))
		for _, a := range abbs {
			cache[a.Abb] = a
			lowercaseWord := strings.ToLower(a.Word)
			lookup[lowercaseWord] = append(lookup[lowercaseWord], a)
		}
	}

	s.cache.UserAbbs[userID] = cache
	s.cache.UserLookup[userID] = lookup
	//log.Printf("Updated cached with %d abbs for user %s\n", len(s.cache.UserAbbs[userID]), userID)
	return nil
}
