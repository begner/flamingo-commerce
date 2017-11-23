package domain

import (
	"context"
	"errors"
)

type (
	// Filter interface for search queries
	Filter interface {
		Value() (string, []string)
	}

	// KeyValueFilter allows simple k -> []values filtering
	KeyValueFilter struct {
		k string
		v []string
	}

	// SearchMeta data
	SearchMeta struct {
		Query      string
		Page       int
		NumPages   int
		NumResults int
	}

	FacetType string

	FacetItem struct {
		Label    string
		Value    string
		Active   bool
		Selected bool

		// Tree Facet
		Items []*FacetItem

		// Range Facet
		Min, Max                 float64
		SelectedMin, SelectedMax float64
	}

	Facet struct {
		Type  string
		Name  string
		Label string
		Items []*FacetItem
	}

	// Result defines a search result for one type
	Result struct {
		SearchMeta SearchMeta
		Hits       []Document
		Facets     map[string]Facet
	}

	// Document holds a search result document
	Document interface{}

	// SearchService defines how to access search
	SearchService interface {
		//Types() []string
		Search(ctx context.Context, filter ...Filter) (results map[string]Result, err error)
		SearchFor(ctx context.Context, typ string, filter ...Filter) (result Result, err error)
	}
)

const (
	ListFacet  FacetType = "ListFacet"
	TreeFacet            = "TreeFacet"
	RangeFacet           = "RangeFacet"
)

var (
	_ Filter = NewKeyValueFilter("a", []string{"b", "c"})

	// ErrNotFound error
	ErrNotFound = errors.New("search not found")
)

// NewKeyValueFilter factory
func NewKeyValueFilter(k string, v []string) *KeyValueFilter {
	return &KeyValueFilter{
		k: k,
		v: v,
	}
}

// Value of the current filter
func (f *KeyValueFilter) Value() (string, []string) {
	return f.k, f.v
}