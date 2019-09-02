package dmd

import (
	"context"
	"os"
	"testing"

	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
)

const (
	dbFilename = "../snomed.db" // real, live database
)

func setUp(tb testing.TB) *terminology.Svc {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		tb.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
	svc, err := terminology.NewService(dbFilename, true)
	if err != nil {
		tb.Fatal(err)
	}
	return svc

}

const (
	amlodipineVmp = 29826211000001109 // Amlodipine 10mg/5ml oral solution sugar free - VMP (sugar free)
	lithiumVmp    = 4559411000001104  // Lithium modified release  - VMP
)

var (
	tags = []language.Tag{terminology.BritishEnglish.Tag()}
)

func TestPrescribableVmp(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	ctx := context.Background()
	amlodipine, err := svc.ExtendedConcept(amlodipineVmp, tags)
	if err != nil {
		t.Fatal(err)
	}
	amlodipineVmp, err := NewVMP(svc, amlodipine)
	if err != nil {
		t.Errorf("amlodipine 10mg/5ml sugar free, not appropriately categorised as a VMP: %s", err)
	}
	if amlodipineVmp.IsSugarFree() == false {
		t.Errorf("amlodipine sugar free is misclassified as containing sugar")
	}
	valid, recommended := amlodipineVmp.PrescribingStatus()
	if !valid || !recommended {
		t.Errorf("amlodipine incorrectly recorded as not being prescribable as a VMP")
	}
	ingredients := amlodipineVmp.SpecificActiveIngredients()
	if len(ingredients) != 1 {
		t.Errorf("incorrect ingredients for amlodipine. expected 1, got: %d(%v)", len(ingredients), ingredients)
	}
	if ingredients[0] != 386864001 { // should be amlodipine substance
		t.Errorf("incorrect ingredient. expected amlodipine got: %d", ingredients[0])
	}
	vtms := amlodipineVmp.GetVTMs()
	if len(vtms) != 1 {
		t.Fatalf("did not return correct VTM for this VMP. expected: 108537001 got: %v", vtms)
	}
	amlodipineVtmID := vtms[0]
	amlodipineVtmEc, err := svc.ExtendedConcept(amlodipineVtmID, tags)
	if err != nil {
		t.Fatal(err)
	}
	amlodipineVtm, err := NewVTM(svc, amlodipineVtmEc)
	if err != nil {
		t.Fatal(err)
	}
	if amlodipineVtm.IsVTM() == false {
		t.Errorf("amlodipine VTM not appropriately classified as a VTM")
	}
	amps, err := amlodipineVmp.GetAMPs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, amp := range amps {
		ec, err := svc.ExtendedConcept(amp, tags)
		if err != nil {
			t.Fatal(err)
		}
		p := NewProduct(svc, ec)
		if p.IsAMP() == false {
			t.Fatalf("%v is not an AMP", p)
		}
	}
	if len(amps) == 0 {
		t.Fatalf("got 0 AMPs for VMP %v", amlodipineVmp)
	}
}
func TestNonPrescribableVmp(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	lithium, err := svc.ExtendedConcept(lithiumVmp, tags)
	if err != nil {
		t.Fatal(err)
	}
	vmp, err := NewVMP(svc, lithium)
	if err != nil {
		t.Errorf("lithium modified release tablets not categorised at a VMP")
	}
	if vmp.IsSugarFree() {
		t.Errorf("lithium m/r misclassified as sugar free")
	}
	valid, recommended := vmp.PrescribingStatus()
	if valid || recommended {
		t.Errorf("lithium modified release incorrectly recorded as valid or recommended for prescribing as VMP")
	}
	vtms := vmp.GetVTMs()
	if len(vtms) == 0 {
		t.Fatal("did not return any VTMs for this VMP")
	}
	for _, vtm := range vtms {
		ec, err := svc.ExtendedConcept(vtm, tags)
		if err != nil {
			t.Fatal(err)
		}
		vtmProduct := NewProduct(svc, ec)
		if vtmProduct.IsVTM() == false {
			t.Fatal("Returned a non-VTM dm+d product during fetch of VTMs for a VMP")
		}
	}
}
