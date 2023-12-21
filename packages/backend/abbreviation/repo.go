package abbreviation

import "log"

// AbbRepository is
type AbbRepository interface {
	GetAbb(listID string, abb string) (*Abbreviation, error)
	GetAbbs(listID string) ([]*Abbreviation, error)
	GetList(id string) (*List, error)
	GetLists(IDs []string) ([]*List, error)
	GetUserLists(userID string) ([]*List, error)
	CreateAbb(listID string, abb *Abbreviation) error
	CreateManyAbbs(listID string, abbs []*Abbreviation) error
	UpdateAbb(listID string, abb *Abbreviation) error
	DeleteAbb(listID, abb string) (*Abbreviation, error)
	DeleteAbbByID(listID, ID string) (*Abbreviation, error)

	CreateList(list *List) error
	CreateManyLists(lists []*List) error
	UpdateList(list *List) error
	DeleteList(id string) (*List, error)

	ImportAbbsToList(listID string, abbs []*Abbreviation) error
}

func (s *abbService) TouchAllAbbs(userIDs []string) {
	var ll []*List
	//counter := 0
	for _, u := range userIDs {
		log.Println("User ID: ", u)
		lists, err := s.GetUserLists(u)
		if err != nil {
			log.Println(err)
			return
		}
		for _, l := range lists {
			//log.Println("List name:", l.Name)
			ll = append(ll, l)
			abbs, _ := s.GetAbbs(l.ID)
			for _, abb := range abbs {
				s.UpdateAbb(l.ID, abb)
			}
		}
	}
	//log.Println("Total number of lists:", len(ll))
	//log.Println("Total number of abbs:", counter)
}
