package app

import (
	"homework10/internal/ads"
	"testing"
)

func FuzzAuthorPredicate(f *testing.F) {
	testcases := []int64{-1, 0, 1}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, s int64) {
		got := authorPredicate(s, ads.Ad{
			AuthorID: 0,
			Text:     "1",
		})
		expect := false
		if got != expect && s != 0 && s != -1 {
			t.Errorf("For (%d) Expect: %t, but got: %t", s, expect, got)
		}
	})
}
