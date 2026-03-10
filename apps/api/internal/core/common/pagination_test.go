package common_test

import (
	"testing"

	"github.com/Sanaruca/condominio/internal/core/common"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ID   int
	Name string
}

func TestPaginator_Sanitize(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected common.Paginator
	}{
		{
			name:  "default values",
			page:  0,
			limit: 0,
			expected: common.Paginator{
				Page:  1,
				Limit: common.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name:  "negative values",
			page:  -1,
			limit: -5,
			expected: common.Paginator{
				Page:  1,
				Limit: common.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name:  "valid values",
			page:  2,
			limit: 10,
			expected: common.Paginator{
				Page:  2,
				Limit: 10,
			},
		},
		{
			name:  "limit exceeds max",
			page:  1,
			limit: 150,
			expected: common.Paginator{
				Page:  1,
				Limit: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := common.Paginator{
				Page:  tt.page,
				Limit: tt.limit,
			}
			p.Sanitize()
			assert.Equal(t, tt.expected.Page, p.Page)
			assert.Equal(t, tt.expected.Limit, p.Limit)
		})
	}
}

func TestNewPaginated(t *testing.T) {
	testData := []testItem{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
	}

	tests := []struct {
		name      string
		data      []testItem
		total     int
		paginator common.Paginator
		expected  common.Paginated[testItem]
	}{
		{
			name:  "first page with default pagination",
			data:  testData,
			total: 50,
			paginator: common.Paginator{
				Page:  1,
				Limit: 20,
			},
			expected: common.Paginated[testItem]{
				Data:  testData,
				Total: 50,
				Page:  1,
				Pages: 3,
				Limit: 20,
			},
		},
		{
			name:  "second page with custom limit",
			data:  testData,
			total: 30,
			paginator: common.Paginator{
				Page:  2,
				Limit: 10,
			},
			expected: common.Paginated[testItem]{
				Data:  testData,
				Total: 30,
				Page:  2,
				Pages: 3,
				Limit: 10,
			},
		},
		{
			name:  "empty result set",
			data:  []testItem{},
			total: 0,
			paginator: common.Paginator{
				Page:  1,
				Limit: 20,
			},
			expected: common.Paginated[testItem]{
				Data:  []testItem{},
				Total: 0,
				Page:  1,
				Pages: 1,
				Limit: 20,
			},
		},
		{
			name:  "exact page division",
			data:  testData,
			total: 40,
			paginator: common.Paginator{
				Page:  2,
				Limit: 20,
			},
			expected: common.Paginated[testItem]{
				Data:  testData,
				Total: 40,
				Page:  2,
				Pages: 2,
				Limit: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.NewPaginated(tt.data, tt.total, tt.paginator)
			assert.Equal(t, tt.expected.Data, result.Data)
			assert.Equal(t, tt.expected.Total, result.Total)
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.Pages, result.Pages)
			assert.Equal(t, tt.expected.Limit, result.Limit)
		})
	}
}

func TestPaginated_WithSanitizedPaginator(t *testing.T) {
	data := []testItem{{ID: 1, Name: "Test"}}

	tests := []struct {
		name      string
		paginator common.Paginator
		expected  common.Paginated[testItem]
	}{
		{
			name: "sanitize page",
			paginator: common.Paginator{
				Page:  0,
				Limit: 10,
			},
			expected: common.Paginated[testItem]{
				Page:  1,
				Limit: 10,
			},
		},
		{
			name: "sanitize limit",
			paginator: common.Paginator{
				Page:  1,
				Limit: 0,
			},
			expected: common.Paginated[testItem]{
				Page:  1,
				Limit: common.DEFAULT_PAGE_SIZE,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := common.NewPaginated(data, 1, tt.paginator)
			assert.Equal(t, tt.expected.Page, result.Page)
			assert.Equal(t, tt.expected.Limit, result.Limit)
		})
	}
}
