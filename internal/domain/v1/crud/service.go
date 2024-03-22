package crud

import "github.com/Filipza/excel-mapping-tool/internal/settings"

type CRUDService[T, L any] interface {
	List(...settings.Option) ([]*L, error)
	Create(*T, ...settings.Option) (*T, error)
	Read(string, ...settings.Option) (*T, error)
	Update(string, *T, ...settings.Option) (*T, error)
	Delete(string, ...settings.Option) (*T, error)
}
