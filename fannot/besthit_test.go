package fannot

import (
	"testing"

	"github.com/hdevillers/go-needle"
	"github.com/hdevillers/go-seq/seq"
)

// Testing the besthit selection
func TestBestHitCheck(t *testing.T) {
	// Create a collection of sequences
	seq1 := seq.Seq{
		Id:       "seq1",
		Sequence: []byte("MEQNRFKKETKTCSASWPRAPQSTLCATDRLELTYDVYTSAERQRRSRTATRLNLVFLHG"),
	}
	// 53 identical bases w/r to seq1
	// 1 insertion
	// Expected similarity: 88.5%
	// Expected identity: 86.9%
	// Expected LRA: 0.9836
	seq2 := seq.Seq{
		Id:       "seq2",
		Sequence: []byte("MEQNRRFKKERKTCSASWARAPQSTLCATDQLERTYDVYTSAERQRRSRAATRLVLVFSHG"),
	}
	// 12 identical bases w/r to seq1
	// 57 gaps
	// Expected similarity: 19.4%
	// Expected identity: 12.2%
	// Expected LRA: 0.7595
	seq3 := seq.Seq{
		Id:       "seq3",
		Sequence: []byte("MDFKRTPSVDSLGTAHSVYSSKTPTRSRSNIGSRTSISVLPSLPTNVSNESIEEDLHNHSNIDEKLESEATPETSYVEF"),
	}

	// Set the seq1 as the query
	bh := NewBestHit(seq1, *needle.NewParam())

	if bh.QuerySeq.Id != "seq1" {
		t.Fatalf(`Wrong query sequence Id, expected 'seq1', obtained %s.`, bh.QuerySeq.Id)
	}

	// Test seq3 as first hit
	bh.CheckHit(seq3)

	if bh.NumHits != 1 {
		t.Fatalf(`Wrong number of checked hits, expected 1 hit, found %d hit(s).`, bh.NumHits)
	}

	if bh.HitSeq.Id != "seq3" {
		t.Fatalf(`Wrong hit sequence Id, expected 'seq3', obtained %s.`, bh.HitSeq.Id)
	}

	// Similarity is between 19% and 20%
	if bh.Similarity < 19.0 || bh.Similarity > 20.0 {
		t.Fatalf(`Hit similarity with 'seq3' should be between 19 and 20, found %.02f.`, bh.Similarity)
	}

	// The length ratio is between 0.75 and 0.76
	if bh.LengthRatio < 0.75 || bh.LengthRatio > 0.76 {
		t.Fatalf(`Hit length ratio with 'seq3' should be between 0.75 and 0.76, found %.02f.`, bh.LengthRatio)
	}

	// Test seq2 as second hit
	bh.CheckHit(seq2)

	if bh.NumHits != 2 {
		t.Fatalf(`Wrong number of checked hits, expected 2 hit, found %d hit(s).`, bh.NumHits)
	}

	// Seq2 should replace seq3 as a better hit
	if bh.HitSeq.Id != "seq2" {
		t.Fatalf(`Wrong hit sequence Id, expected 'seq2', obtained %s.`, bh.HitSeq.Id)
	}

	// Similarity is between 19% and 20%
	if bh.Similarity < 88.0 || bh.Similarity > 89.0 {
		t.Fatalf(`Hit similarity with 'seq2' should be between 88 and 89, found %.02f.`, bh.Similarity)
	}

	// The length ratio is between 0.98 and 0.98
	if bh.LengthRatio < 0.98 || bh.LengthRatio > 0.99 {
		t.Fatalf(`Hit length ratio with 'seq2' should be between 0.98 and 0.99, found %.02f.`, bh.LengthRatio)
	}

	// Test again seq3, it should not replace seq2
	bh.CheckHit(seq3)

	if bh.NumHits != 3 {
		t.Fatalf(`Wrong number of checked hits, expected 3 hits, found %d hit(s).`, bh.NumHits)
	}

	if bh.HitSeq.Id != "seq2" {
		t.Fatalf(`Wrong hit sequence Id, expected 'seq2', obtained %s.`, bh.HitSeq.Id)
	}

	// Similarity is still between 19% and 20%
	if bh.Similarity < 88.0 || bh.Similarity > 89.0 {
		t.Fatalf(`Hit similarity with 'seq2' should be still between 88 and 89, found %.02f.`, bh.Similarity)
	}

	// Now test seq1, a perfect hit
	bh.CheckHit(seq1)

	if bh.NumHits != 4 {
		t.Fatalf(`Wrong number of checked hits, expected 4 hits, found %d hit(s).`, bh.NumHits)
	}

	if bh.HitSeq.Id != "seq1" {
		t.Fatalf(`Wrong hit sequence Id, expected 'seq1', obtained %s.`, bh.HitSeq.Id)
	}

	if bh.Similarity != 100.0 {
		t.Fatalf(`Hit identical to queyr, expected a similarity of 100, obtained %.02f.`, bh.Similarity)
	}

	if bh.LengthRatio != 1.0 {
		t.Fatalf(`Hit identical to queyr, expected a length ratio of 1, obtained %.02f.`, bh.LengthRatio)
	}
}
