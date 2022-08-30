package url

import "testing"

func TestUrlRotator_Next(t *testing.T) {
	rotator := NewUrlRotator([]string{"1", "2", "3"})
	url1, url2, url3 := rotator.Next(), rotator.Next(), rotator.Next()

	next := rotator.Next()
	if url1 != next {
		t.Errorf("expected %s, got %s", url1, next)
	}
	next = rotator.Next()
	if url2 != next {
		t.Errorf("expected %s, got %s", url2, next)
	}
	next = rotator.Next()
	if url3 != next {
		t.Errorf("expected %s, got %s", url3, next)
	}
}
