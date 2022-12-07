package document

import (
	"os"
	"time"

	"github.com/jaevor/go-nanoid"
)

type DocService interface {
	CreateDoc(doc *Document) (*Document, error)
	CreateManyDocs(docs []*Document) error
	ImportDoc(file *os.File) (string, error)
	GetDoc(id string) (*Document, error)
	GetDocs(userID string) ([]*Document, error)
	UpdateDoc(doc *Document) (*Document, error)
	DeleteDoc(id string) error
}

type docService struct {
	repo DocRepository
}

func NewDocService(repo DocRepository) DocService {
	d := docService{
		repo: repo,
	}
	return &d
}

func (s *docService) GetDoc(id string) (*Document, error) {
	doc, err := s.repo.GetDoc(id)
	return doc, err
}

func (s *docService) GetDocs(userID string) ([]*Document, error) {
	docs, err := s.repo.GetDocs(userID)
	return docs, err
}

func (s *docService) CreateDoc(doc *Document) (*Document, error) {
	id, err := nanoid.CustomASCII("abcdef0123456789", 8)
	if err != nil {
		panic(err)
	}
	doc.ID = id()
	doc.Created = time.Now()
	doc.Updated = time.Now()
	err = s.repo.CreateDoc(doc)
	return doc, err
}

func (s *docService) CreateManyDocs(docs []*Document) error {
	return s.repo.CreateManyDocs(docs)
}

func (s *docService) UpdateDoc(doc *Document) (*Document, error) {
	doc.Updated = time.Now()
	err := s.repo.UpdateDoc(doc)
	return doc, err
}

func (s *docService) DeleteDoc(id string) error {
	err := s.repo.DeleteDoc(id)
	return err
}
