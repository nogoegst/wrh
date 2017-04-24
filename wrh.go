// wrh.go - Weighted Rendezvous Hash implementation.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to wrh, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package wrh

import (
	"math"
	"sort"

	"github.com/nogoegst/blake2xb"
	"github.com/nogoegst/rand"
)

// hashFloat64 calculates hash of key||seed and converts
// result into float64 number in [0:1).
func hashFloat64(key string, seed string) float64 {
	cfg := blake2xb.NewXConfig(0)
	h, err := blake2xb.NewX(cfg)
	if err != nil {
		panic(err)
	}
	h.Write([]byte(key))
	h.Write([]byte(seed))
	rnd := rand.NewWithReader(h)
	return rnd.Float64()
}

// WRH bucket
type Bucket interface {
	Name() string
	Weight() float64
}

// WeightScore calculated weighted score of bucket b for
// given key.
func WeightedScore(b Bucket, key string) float64 {
	return -b.Weight() / math.Log(hashFloat64(key, b.Name()))
}

// Sort sorts buckets by weighted score for given key in descending order.
// It returns resulting slice and does not modify buckets.
func Sort(buckets interface{}, key string) []Bucket {
	bckts := buckets.(map[string]Bucket)
	bs := make([]Bucket, 0, len(bckts))
	for _, b := range bckts {
		bs = append(bs, b.(Bucket))
	}
	sort.Slice(bs, func(i, j int) bool {
		return WeightedScore(bs[i], key) >= WeightedScore(bs[j], key)
	})
	return bs
}
