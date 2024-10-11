package domain

// an abstraction of the main application logic, which is driven by external agents such as the REST APIs that we build later in the article.
type Service interface {
	Find(code string) (*Product, error)
	Store(product *Product) error
	Update(product *Product) error
	FindAll() ([]*Product, error)
	Delete(code string) error
}

// enables the product catalogue service to connect to data store adapters to persist and query the product catalogue data
type Repository interface {
	Find(code string) (*Product, error)
	Store(product *Product) error
	Update(product *Product) error
	FindAll() ([]*Product, error)
	Delete(code string) error
}
