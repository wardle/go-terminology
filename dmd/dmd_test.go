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
