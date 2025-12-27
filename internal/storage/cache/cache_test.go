package cache

import (
	"testing"
	"time"
)

func addTestDataToCache(cache *Cache) {
	cache.cache["id_1"] = time.Now()
	cache.cache["id_2"] = time.Now()
	cache.cache["id_3"] = time.Now()
	cache.cache["id_4"] = time.Now()
}

func TestCache_CheckBySubID(t *testing.T) {
	testCases := []struct {
		name         string
		cache        *Cache
		checkIDS     []string
		checkResults []bool
	}{
		{
			name:         "test_1_low_ttl_expect_deleted",
			cache:        NewCache(time.Microsecond),
			checkIDS:     []string{"id_1", "id_2", "id_3", "id_4"},
			checkResults: []bool{false, false, false, false},
		},
		{
			name:         "test_2_long_ttl_expect_true",
			cache:        NewCache(time.Second * 5),
			checkIDS:     []string{"id_1", "id_2", "id_3", "id_4"},
			checkResults: []bool{true, true, true, true},
		},
		{
			name:         "non_existent_id",
			cache:        NewCache(time.Second),
			checkIDS:     []string{"non_existent", "another_missing"},
			checkResults: []bool{false, false},
		},
		{
			name:         "empty_string_id",
			cache:        NewCache(time.Second),
			checkIDS:     []string{"", "  ", "\t", "\n"},
			checkResults: []bool{false, false, false, false},
		},
		{
			name:         "case_sensitive_ids",
			cache:        NewCache(time.Second),
			checkIDS:     []string{"ID", "id", "Id", "iD"},
			checkResults: []bool{false, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			addTestDataToCache(tc.cache)
			for i, id := range tc.checkIDS {
				res := tc.cache.CheckBySubID(id)
				if res != tc.checkResults[i] {
					t.Errorf("EXPECTED: got: %v, expect: %v\n", res, tc.checkResults[i])
				}

				if _, ok := tc.cache.cache[id]; ok != tc.checkResults[i] {
					t.Errorf("cache by id: %s not deleted!\n", id)
				}
			}
		})
	}
}
