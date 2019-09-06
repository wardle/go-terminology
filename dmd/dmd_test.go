package dmd

import (
	"context"
	"fmt"
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
	dispensedDoseForm := amlodipineVmp.DisposedDoseForm()
	if dispensedDoseForm != 385023001 { // oral solution
		t.Errorf("amlodipine oral solution not categorised as oral solution. expected: 385023001. got %d", dispensedDoseForm)
	}
	vtms := amlodipineVmp.VTMs()
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
	vmpps, err := amlodipineVmp.VMPPs()
	if err != nil {
		t.Fatal(err)
	}
	for _, vmpp := range vmpps {
		ec, err := svc.ExtendedConcept(vmpp, tags)
		if err != nil {
			t.Fatal(err)
		}
		p := NewProduct(svc, ec)
		if p.IsVMPP() == false {
			t.Fatalf("Get VMPPs did not return a VMPP. Got: %v", p)
		}
	}

	amps, err := amlodipineVmp.AllAMPs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v has %d AMPs\n", amlodipineVmp, len(amps))
	for _, amp := range amps {
		ec, err := svc.ExtendedConcept(amp, tags)
		if err != nil {
			t.Fatal(err)
		}
		p, err := NewAMP(svc, ec)
		if err != nil {
			t.Fatal(err)
		}
		if p.IsAMP() == false {
			t.Fatalf("%v is not an AMP", p)
		}
		vmp, err := p.VMP()
		if err != nil {
			t.Fatal(err)
		}
		if vmp != amlodipineVmp.Concept.Id {
			t.Fatalf("incorrect VMP for AMP: %v. expected: %v, got:%v", p, amlodipineVmp, vmp)
		}
	}
	if len(amps) == 0 {
		t.Fatalf("got 0 AMPs for VMP %v", amlodipineVmp)
	}
	vtmIngredients, err := amlodipineVtm.SpecificActiveIngredients(ctx, tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(vtmIngredients) != 1 {
		t.Fatalf("incorrect ingredients for amlodipine VTM. expected 1. got:%v", vtmIngredients)
	}
	if vtmIngredients[0] != 386864001 {
		t.Fatalf("incorrect ingredients for amlodipine VTM. expected 386864001. got:%d", vtmIngredients[0])
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
	vtms := vmp.VTMs()
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

// VMPs have children that are VMPs
func TestVmpsHaveChildVmps(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	allVmps, err := svc.ReferenceSetComponents(VmpReferenceSet)
	if err != nil {
		t.Fatal(err)
	}
	vmpsHaveChildVmps := 0
	for vmpID := range allVmps {
		children, err := svc.Children(vmpID)
		if err != nil {
			t.Fatal(err)
		}
		for _, child := range children {
			if _, exists := allVmps[child]; exists {
				vmpsHaveChildVmps++
			}
		}
	}
	if vmpsHaveChildVmps == 0 {
		t.Fatal("VMP structures have changed. No VMPs have children that are VMPs")
	}
}

//AMPs do not have children that are AMPs
func TestAmpsDoNotHaveChildAmps(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	allAmps, err := svc.ReferenceSetComponents(AmpReferenceSet)
	if err != nil {
		t.Fatal(err)
	}
	ampsHaveChildAmps := 0
	for ampID := range allAmps {
		children, err := svc.Children(ampID)
		if err != nil {
			t.Fatal(err)
		}
		for _, child := range children {
			if _, exists := allAmps[child]; exists {
				ampsHaveChildAmps++
			}
		}
	}
	if ampsHaveChildAmps != 0 {
		t.Fatalf("AMP structures have changed. There are %d AMPs that have children that are AMPs", ampsHaveChildAmps)
	}
}

func TestVTMsHaveChildVTMs(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	allVtms, err := svc.ReferenceSetComponents(VtmReferenceSet)
	if err != nil {
		t.Fatal(err)
	}
	vtmsHaveChildVtms := 0
	for vtmID := range allVtms {
		children, err := svc.Children(vtmID)
		if err != nil {
			t.Fatal(err)
		}
		for _, child := range children {
			if _, exists := allVtms[child]; exists {
				vtmsHaveChildVtms++
			}
		}
	}
	if vtmsHaveChildVtms == 0 {
		t.Fatal("VTM structures have changed. There are 0 VTMs that have children that are VTMs")
	}
}

// Language handling, particularly derivation of preferred term, for dm+d, is not quite the same
// as normal SNOMED CT rules for language. Instead, preferred terms are derived by membership of
// the dm+d realm description reference set (999000671000001103).
// See https://www.nhsbsa.nhs.uk/sites/default/files/2017-02/Secondary_Care_Electronic_Prescribing_Implementation_Guidance_5_0.pdf
// although this documentation is out of date as dm+d preferred synonym for otic is now "Ear".

func TestLanguageDmd(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()

	tests := []struct {
		conceptID int64
		usual     string // what we should get from standard SNOMED CT
		expected  string // what we expect for dm+d / for e-prescribing
	}{
		{10547007, "Otic route", "Ear"},
		{6064005, "Topical route", "Topical"},
		{420254004, "Body cavity route", "Body cavity"},
		{26643006, "Oral route", "Oral"},
	}
	for _, test := range tests {
		d1, err := svc.PreferredSynonym(test.conceptID, tags)
		if err != nil {
			t.Fatal(err)
		}
		d2, err := svc.PreferredSynonymByReferenceSet(test.conceptID, NhsDmdRealmLanguageReferenceSet, tags)
		if err != nil {
			t.Fatal(err)
		}
		if d1.Id == d2.Id {
			t.Fatalf("Standard preferred term and dm+d term should be different for %d", test.conceptID)
		}
		if d1.Term != test.usual {
			t.Fatalf("incorrect standard preferred term, expected:'%s', got:'%v'", test.usual, d1)
		}
		if d2.Term != test.expected {
			t.Fatalf("incorrect dm+d term, expected:'%s', got:'%v'", test.expected, d2)
		}
	}

}

// This tests the structures of dm+d rather than the code...
// ensuring that all TFs in dm+d have only ONE active trade family group relationship
func TestTradeFamilyGroupOrdinality(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tfs, err := svc.ReferenceSetComponents(TfReferenceSet)
	if err != nil {
		t.Fatal(err)
	}
	for tfID := range tfs {
		ec, err := svc.ExtendedConcept(tfID, tags)
		if err != nil {
			t.Fatal(err)
		}
		tf := NewProduct(svc, ec)
		rels := tf.GetRelationships()
		count := 0
		for _, rel := range rels {
			if rel.Active && rel.TypeId == HasTradeFamilyGroup {
				count++
			}
		}
		if count > 1 {
			t.Errorf("TF '%s'(%d) has more than more TF group. expected:1, got:%d\n", tf.PreferredDescription.Term, tf.Concept.Id, count)
		}
	}
}

func TestTradeFamily(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tests := []struct {
		ampConceptID int64
		expectedTF   int64
		expectedTFG  int64
	}{
		{172711000001100, 9496401000001106, 0},                   // istin - no TFG
		{11497911000001105, 9642801000001104, 12809001000001109}, //  XPHEN TYR Maxamum powder (SHS International Ltd) (product) - XPHEN TYR Maxamum (product) - "maxamum"
	}
	for _, test := range tests {
		ampEC, err := svc.ExtendedConcept(test.ampConceptID, tags) // this is the AMP
		if err != nil {
			t.Fatal(err)
		}
		amp, err := NewAMP(svc, ampEC)
		if err != nil {
			t.Fatal(err)
		}
		tfID, err := amp.TF()
		if err != nil {
			t.Fatal(err)
		}
		if tfID != test.expectedTF {
			t.Fatalf("incorrect TF for %v, expected:%d, got:%d", amp, test.expectedTF, tfID)
		}
		tfEC, err := svc.ExtendedConcept(tfID, tags)
		if err != nil {
			t.Fatal(err)
		}
		tf, err := NewTF(svc, tfEC)
		if err != nil {
			t.Fatal(err)
		}
		if tfg := tf.TradeFamilyGroup(); tfg != test.expectedTFG {
			t.Fatalf("incorrect trade family group for %v. expected:%d, got:%d", tf, test.expectedTFG, tfg)
		}
	}

}
