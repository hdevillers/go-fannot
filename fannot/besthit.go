package fannot

import (
	"fmt"

	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
)

type BestHit struct {
	QuerySeq    seq.Seq
	HitSeq      seq.Seq
	NumHits     int
	Similarity  float64
	LengthRatio float64
	NeedleParam needle.Param
	IdRule      int
}

func NewBestHit(q seq.Seq, np needle.Param) *BestHit {
	var bh BestHit
	bh.QuerySeq = q
	bh.NumHits = 0
	bh.Similarity = 0.0
	bh.LengthRatio = 0.0
	bh.NeedleParam = np
	bh.IdRule = -1

	return &bh
}

// Return the minimal length ratio
func computeMinLengthRatio(l1, l2 int) float64 {
	if l1 < l2 {
		return float64(l1) / float64(l2)
	} else {
		return float64(l2) / float64(l1)
	}
}

/*
   Consider a new hit and store it if it is the first on or
   if it is the best one.
*/
func (bh *BestHit) CheckHit(h seq.Seq) {
	// Increment hit counter
	bh.NumHits++

	// Use needle alignment to compare hit with the query
	ndl := needle.NewNeedle(bh.QuerySeq, h)
	ndl.Par = &bh.NeedleParam
	err := ndl.Align()
	if err != nil {
		panic(fmt.Sprintf("Failed to align query %s againt ref %s, error: %s.", bh.QuerySeq.Id, h.Id, err.Error()))
	}

	// Check if the hit is better than the one stored
	if ndl.Rst.GetSimilarityPct() > bh.Similarity {
		bh.HitSeq = h
		bh.Similarity = ndl.Rst.GetSimilarityPct()
		bh.LengthRatio = computeMinLengthRatio(bh.QuerySeq.Length(), h.Length())
	}
}
