package fannot

import (
	"fmt"
	"strings"
)

type Result struct {
	Note        string
	Product     string
	GeneName    string
	Function    string
	Status      int
	Similarity  float64
	LengthRatio float64
	HitNum      int
	HitOvrWrt   bool
	HitId       string
	HitLocus    string
	HitSpecies  string
	IpsId       []string
	IpsAnnot    []string
	DbId        string
}

// Generate result header
func Header() string {
	return "GeneID\tProduct\tNote\tFunction\tOrganism\tRefID\tRefLocus\tRefName\tIPSID\tIPSAnnot\tStatus\tSimilarity\tLengthRatio\tDBID\tHitNum\tOverWritten\n"
}

// Convert a Result object into string
func (r *Result) ToString(gid string) string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%.03f\t%.03f\t%s\t%d\t%t\n",
		gid, r.Product, r.Note, r.Function, r.HitSpecies,
		r.HitId, r.HitLocus, r.GeneName, strings.Join(r.IpsId, ","),
		strings.Join(r.IpsAnnot, "; "), r.Status, r.Similarity,
		r.LengthRatio, r.DbId, r.HitNum, r.HitOvrWrt,
	)
}
