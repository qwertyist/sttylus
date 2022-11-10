package abbreviation

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
)

//AbbService is
type AbbService interface {
	Abbreviate(userID, abbQuery string) (string, bool)
	Lookup(userID, abb string) ([]*Abbreviation, bool)

	CreateAbb(listID string, abb *Abbreviation) error
	CreateManyAbbs(listID string, abbs []*Abbreviation) error
	GetAbb(listID, abb string) (*Abbreviation, error)
	GetAbbs(listID string) ([]*Abbreviation, error)
	UpdateAbb(listID string, abb *Abbreviation) error
	DeleteAbb(listID string, abb string) (*Abbreviation, error)
	DeleteAbbByID(listID string, id string) (*Abbreviation, error)
	CopyStandardList(userID string) (string, error)
	CreateList(list *List) (string, error)
	CreateManyLists(lists []*List) error
	GetList(id string) (*List, error)
	GetLists(listIDs []string) ([]*List, error)
	GetUserLists(userID string) ([]*List, error)
	UpdateList(list *List) error
	TouchList(listID string) error
	DeleteList(userID string, id string) (*List, error)

	InitCache(users []string)
	Cache(q query) error
	UpdateCache(userID string) error

	InitSharedList(userID string, baseListID string) string
	JoinSharedList(userid, listid string) error
	CreateSharedAbb(listid string, abb Abbreviation) error
	RemoveSharedAbb(listid, abb string) error
	GetSharedAbbs(listid string) ([]Abbreviation, error)

	GetSuggestions(userID string) ([]string, error)
	IgnoreSuggestion(userID, word string) bool
	IgnoreAllSuggestions(userID string) error
	DontRemindAbb(userID, abb string) error

	Import(userid string, abbs []*Abbreviation) []*Abbreviation
	ImportAbbsToList(listID string, abbs []*Abbreviation) error
	ImportProtype(userID string, dat []byte) ([]*Abbreviation, error)
	ImportTextOnTop(userID string, dat []byte) (map[string][]*Abbreviation, error)
	ImportIllumiType(userID string, dat []byte) (map[string][]*Abbreviation, error)
	CheckForConflicts(listID string, abbs []*Abbreviation) ([]conflict, error)
	CreateToTExport(standard string, addon []string) (textOnTopExport, error)
	CreateProtypeExport(user, standard string, addon []string) error
	GetPublicList(short_id string) (*List, error)
	CreatePublicList(id string, list *List) error

	TouchAllAbbs(userIDs []string)
}

type abbService struct {
	repo           AbbRepository
	cache          abbCache
	suggestions    map[string]int
	serverURL      string
	globalStandard string
}

//NewAbbService returns an abbService connected to the repo db and initializes the suggestion datatype
func NewAbbService(repo AbbRepository, globalStandard string) AbbService {
	suggs := make(map[string]int)
	a := abbService{
		repo:           repo,
		suggestions:    suggs,
		globalStandard: globalStandard,
	}
	return &a
}

func capitalize(abb, word string) string {
	abbLength := len(abb)
	n := len(abb)
	if n == 1 {
		r, _ := utf8.DecodeRuneInString(abb)
		if !unicode.IsUpper(r) {
			return word
		}
		return strings.Title(word)
	}
	for len(abb) > 0 {
		r, size := utf8.DecodeRuneInString(abb)
		if n == abbLength && !unicode.IsUpper(r) {
			return word
		}
		if unicode.IsUpper(r) {
			n--
		}
		abb = abb[size:]
	}
	if n != 0 {
		capitalized := strings.Split(word, " ")
		if len(capitalized) > 1 {
			return strings.Title(capitalized[0]) + " " + strings.Join(capitalized[1:], " ")
		}

		return strings.Title(capitalized[0])
	}
	return strings.ToUpper(word)
}

func (s *abbService) Abbreviate(userID, abbQuery string) (string, bool) {
	//log.Printf("service|Abbreviate got abb %s from user %s\n", q.Abb, q.UserID)
	sharedListID := s.cache.UserSharedList[userID]
	if sharedListID != "" {
		w, f := s.cache.SharedAbbs[sharedListID][strings.ToLower(abbQuery)]
		if f {
			w = capitalize(abbQuery, w)
			return w, false
		}
	}
	word, found := s.cache.UserAbbs[userID][strings.ToLower(abbQuery)]
	/*	for _, a := range s.cache.UserAbbs[q.UserID] {
		fmt.Println("in this cache:", a)
	}*/
	if found {
		abb := capitalize(abbQuery, word.Word)
		return abb, false
	}
	missed, found := s.cache.UserLookup[userID][strings.ToLower(abbQuery)]

	if found {
		for _, a := range missed {
			if a.Remind {
				return missed[len(missed)-1].Abb, true
			}
			return abbQuery, false
		}
	}

	if len(abbQuery) > 5 {
		sugg := strings.TrimRight(abbQuery, "\\/^'\"#¤%&¨*<>|=+`´.,!?()[]{}:;-_")
		sugg = strings.ToLower(sugg)
		//log.Printf("%s is a long word without an abbreviation\n", q.Abb)

		if s.cache.UserSuggestions[userID][sugg] == 0 {
			s.cache.UserSuggestions[userID][sugg] = 1
			//		log.Printf("storing suggestion\n")
		}
	}

	return abbQuery, false
}

func (s *abbService) Lookup(userID, word string) ([]*Abbreviation, bool) {
	abbs, found := s.cache.UserLookup[userID][strings.ToLower(word)]
	return abbs, found
}

func (s *abbService) GetAbb(listID, ID string) (*Abbreviation, error) {
	abb, err := s.repo.GetAbb(listID, ID)
	return abb, err
}
func (s *abbService) GetAbbs(listID string) ([]*Abbreviation, error) {
	abbs, err := s.repo.GetAbbs(listID)
	return abbs, err
}

func (s *abbService) CreateAbb(listID string, abb *Abbreviation) error {
	abb.ID = uuid.New().String()
	abb.Updated = time.Now()
	abb.Remind = true
	abb.Abb = strings.ToLower(abb.Abb)
	err := s.repo.CreateAbb(listID, abb)

	err = s.TouchList(listID)
	if err != nil {
		return fmt.Errorf("createAbb|couldn't touch list:\n%s", err.Error())
	}

	return err
}

func (s *abbService) CreateManyAbbs(listID string, abbs []*Abbreviation) error {
	return s.repo.CreateManyAbbs(listID, abbs)
}

func (s *abbService) UpdateAbb(listID string, abb *Abbreviation) error {
	tmp, err := s.repo.GetAbb(listID, abb.Abb)
	if tmp == nil {
		return fmt.Errorf("no abb to update")
	}
	abb.Updated = time.Now()
	err = s.repo.UpdateAbb(listID, abb)
	if err != nil {
		return fmt.Errorf("service UpdateAbb failed: %s", err)
	}
	err = s.TouchList(listID)
	return err
}

func (s *abbService) DeleteAbb(listID, abb string) (*Abbreviation, error) {
	deleted, err := s.repo.DeleteAbb(listID, abb)
	if err != nil || deleted == nil {
		return nil, fmt.Errorf("service DeleteAbb failed: %s", err)
	}
	s.TouchList(listID)

	return deleted, err
}

func (s *abbService) DeleteAbbByID(listID, ID string) (*Abbreviation, error) {
	deleted, err := s.repo.DeleteAbbByID(listID, ID)
	if err != nil || deleted == nil {
		return nil, fmt.Errorf("service DeleteAbbByID failed: %s", err)
	}
	s.TouchList(listID)

	return deleted, err
}

func (s *abbService) GetList(listID string) (*List, error) {
	list, err := s.repo.GetList(listID)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (s *abbService) GetLists(listIDs []string) ([]*List, error) {
	lists, err := s.repo.GetLists(listIDs)
	return lists, err
}

func (s *abbService) GetUserLists(userID string) ([]*List, error) {
	return s.repo.GetUserLists(userID)
}

func (s *abbService) CopyStandardList(userID string) (string, error) {
	g, err := s.GetList(s.globalStandard)
	if err != nil {
		return "", err
	}
	abbs, err := s.GetAbbs(g.ID)
	g.Creator = userID
	newID, err := s.CreateList(g)
	if err != nil {
		return "", err
	}
	for _, a := range abbs {
		a.Creator = userID
		s.CreateAbb(newID, a)
	}

	log.Printf("global: %s, new: %s\n", g.ID, newID)
	return newID, nil
}

func (s *abbService) CreateList(list *List) (string, error) {
	list.ID = uuid.New().String()
	list.Created = time.Now()
	list.Updated = time.Now()

	err := s.repo.CreateList(list)
	if err != nil {
		return "", fmt.Errorf("service|CreateList failed: %q", err)
	}
	return list.ID, nil
}

func (s *abbService) CreateManyLists(lists []*List) error {
	return s.repo.CreateManyLists(lists)
}

func (s *abbService) UpdateList(list *List) error {
	list.Updated = time.Now()
	err := s.repo.UpdateList(list)
	return err
}

func (s *abbService) TouchList(listID string) error {
	l, err := s.repo.GetList(listID)
	if err != nil {
		return fmt.Errorf("service TouchList failed: %s", err)
	}
	l.Updated = time.Now()
	abbs, err := s.repo.GetAbbs(listID)
	if err != nil {
		return fmt.Errorf("service TouchList couldn't count abbs: %s", err)
	}
	l.Counter = len(abbs)
	err = s.repo.UpdateList(l)
	return err
}

func (s *abbService) DeleteList(userID, listID string) (*List, error) {
	list, err := s.repo.DeleteList(listID)
	return list, err
}

func (s *abbService) GetSuggestions(userID string) ([]string, error) {
	var suggs []string
	//log.Printf("userID %s get suggestions:\n", userID)
	for word, keep := range s.cache.UserSuggestions[userID] {
		if keep == 1 {
			suggs = append(suggs, word)
		}
	}
	return suggs, nil
}

func (s *abbService) IgnoreSuggestion(userID, word string) bool {
	suggested := s.cache.UserSuggestions[userID][word]
	if suggested <= 0 {
		return false
	}
	s.cache.UserSuggestions[userID][word] = -1
	return true
}

func (s *abbService) IgnoreAllSuggestions(userID string) error {
	for word := range s.cache.UserSuggestions[userID] {
		s.cache.UserSuggestions[userID][word] = -1
	}
	return nil
}

func (s *abbService) DontRemindAbb(userID, abb string) error {
	a, found := s.cache.UserAbbs[userID][abb]
	if !found {
		log.Println("nothing stored for abb:", abb)
		return nil
	}
	a.Remind = false

	ll := s.cache.UserAbbLists[userID]
	for _, l := range ll {
		s.UpdateAbb(l, a)
	}
	err := s.UpdateCache(userID)

	return err

}
