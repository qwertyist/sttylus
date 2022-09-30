package document

type DocRepository interface {
	Init() error

	CreateDoc(doc *Document) error
	CreateManyDocs(docs []*Document) error
	GetDoc(id string) (*Document, error)
	GetDocs(userID string) ([]*Document, error)
	UpdateDoc(doc *Document) error
	DeleteDoc(id string) error
}
