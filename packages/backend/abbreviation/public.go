package abbreviation

import (
	"fmt"
	"strings"
)

func (s *abbService) GetPublicList(short_id string) (*List, error) {
	list, ok := s.cache.PublicLists[short_id]
	//log.Println("GetPublicList[list, ok]:", list, ok)
	if ok {
		return list, nil
	}
	return nil, fmt.Errorf("no public list corresponding to short_id")

}

func (s *abbService) CreatePublicList(id string, list *List) error {
	short_id := strings.Split(id, "-")
	if list == nil {
		list.ID = id
	}
	//log.Println("short_id:", short_id)
	s.cache.PublicLists[short_id[0]] = list
	return nil
}
