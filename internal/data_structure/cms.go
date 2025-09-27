package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// Log10PointFive is a precomputed value for log10(0.5).
const Log10PointFive = -0.30102999566

// CMS is the Count-Min Sketch data structure.
// The counter field has been changed to a 2D slice for better clarity and indexing.
type CMS struct {
	width uint64
	depth uint64
	// counter is now a 2D slice of uint32. The outer slice represents the rows (depth),
	// and the inner slice represents the columns (width).
	counters [][]uint64
}

func CreateCMS(width uint64, depth uint64) *CMS {
	cms := &CMS{
		width: width,
		depth: depth,
	}

	cms.counters = make([][]uint64, depth)
	for i := uint64(0); i < depth; i++ {
		cms.counters[i] = make([]uint64, width)
	}
	return cms
}

// CalcCMSDim calculates the dimensions (width and depth) of the CMS
// based on the desired error rate and probability.
func CalcCMSDim(errRate float64, errProb float64) (int64, int64) {
	w := int64(math.Ceil(2.0 / errRate))
	d := int64(math.Ceil(math.Log10(errProb) / Log10PointFive))
	return w, d
}

// calcHash calculates a 32-bit hash for the given item and seed.
func (c *CMS) calcHash(item string, seed int32) uint64 {
	hasher := murmur3.New64WithSeed(uint32(seed))
	hasher.Write([]byte(item))
	return hasher.Sum64()
}

// IncrBy increments the count for an item by a specific value.
// It returns the estimated count for the item after the increment.
func (c *CMS) IncrBy(item string, value int64) uint64 {
	var minCount uint64 = math.MaxUint64

	// loop through each row
	for i := uint64(0); i < c.depth; i++ {
		// Calculate a new hash for each row using the row index as the seed.
		hash := c.calcHash(item, int32(i))
		// Use the hash to get the column index within the row.
		j := hash % c.width

		// Safely add the value to prevent overflow.
		if math.MaxUint64-c.counters[i][j] < uint64(value) {
			c.counters[i][j] = math.MaxUint64
		} else {
			c.counters[i][j] += uint64(value)
		}
		// Update the minimum count if the current count is smaller.
		if c.counters[i][j] < minCount {
			minCount = c.counters[i][j]
		}
	}
	return minCount
}

// Count returns the estimated count for an item.
// It retrieves the minimum count across all hash functions to provide the most accurate estimate.
func (c *CMS) Count(item string) uint64 {
	var minCount uint64 = math.MaxUint64

	// loop through each row
	for i := uint64(0); i < c.depth; i++ {
		// Calculate a new hash for each row using the row index as the seed.
		hash := c.calcHash(item, int32(i))
		// Use the hash to get the column index within the row.
		j := hash % c.width
		// Update the minimum count if the current count is smaller.
		if c.counters[i][j] < minCount {
			minCount = c.counters[i][j]
		}
	}
	return minCount
}
