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
	h, err := blake2xb.New(16) // XXX: how many calls does rand make?
	if err != nil {
		panic(err)
	}
	h.Write([]byte(key))
	h.Write([]byte(seed))
	rnd := rand.NewWithReader(h)
	return rnd.Float64()
}

// Table represents WRH hash table
type Table map[string]float64

// WeightScore calculated weighted score of bucket named name with
// weight for given key.
func WeightedScore(key, name string, weight float64) float64 {
	return -weight / math.Log(hashFloat64(key, name))
}

// Calc calculates new table containing all the scores for the
// given key.
func (t Table) Calc(key string) Table {
	wt := make(Table)
	for name, weight := range t {
		wt[name] = WeightedScore(key, name, weight)
	}
	return wt
}

type scoredValue struct {
	Name  string
	Score float64
}

// Sort sorts score table by score for given key in descending order.
// It returns slice of sorted keys.
func (t Table) Sort() []string {
	bs := make([]scoredValue, 0, len(t))
	for k, v := range t {
		bs = append(bs, scoredValue{Name: k, Score: v})
	}
	sort.Slice(bs, func(i, j int) bool {
		return bs[i].Score >= bs[j].Score
	})
	sm := make([]string, 0)
	for _, v := range bs {
		sm = append(sm, v.Name)
	}
	return sm
}

// List is a convenience shourtcut function that returns list
// of keys sorted by their score.
func (t Table) List(key string) []string {
	return t.Calc(key).Sort()
}
