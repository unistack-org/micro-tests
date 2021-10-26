package id

import (
	"testing"

	"github.com/stretchr/testify/require"
	mid "go.unistack.org/micro/v3/util/id"
)

func TestHasNoCollisions(t *testing.T) {
	tries := 100_000
	used := make(map[string]bool, tries)
	for i := 0; i < tries; i++ {
		id := mid.Must()
		require.False(t, used[id], "shouldn't return colliding IDs")
		used[id] = true
	}
}

func TestFlatDistribution(t *testing.T) {
	tries := 100_000
	alphabet := "abcdefghij"
	size := 10
	chars := make(map[rune]int)
	for i := 0; i < tries; i++ {
		id := mid.Must(mid.Alphabet(alphabet), mid.Size(size))
		for _, r := range id {
			chars[r]++
		}
	}

	for _, count := range chars {
		require.InEpsilon(t, size*tries/len(alphabet), count, .01, "should have flat distribution")
	}
}

// Benchmark id generator
func BenchmarkNanoid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = mid.New()
	}
}
