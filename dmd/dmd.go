// Package dmd provides functionality to process and understand data from the
// UK Dictionary of Medicines and Devices (dm+d).
// This is, by definition, a UK-only module and will not give expected results for
// drugs outside of the UK Product root.
//
//
// See https://www.nhsbsa.nhs.uk/sites/default/files/2018-10/doc_SnomedCTUKDrugExtensionModel%20-%20v1.0.pdf
package dmd

import (
	"fmt"
	"strings"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
)

// dm+d related constants - definitions of relationship types
const (
	// Core concepts - types of dm+d product
	UKProduct                   = 10363601000001109
	ActualMedicinalProduct      = 10363901000001102
	ActualMedicinalProductPack  = 10364001000001104
	VirtualMedicinalProduct     = 10363801000001108
	VirtualMedicinalProductPack = 8653601000001108
	VirtuaTherapeuticMoiety     = 10363701000001104
	TradeFamily                 = 9191801000001103

	// dm+d reference sets - membership of a reference set tells us which of six types of product
	VtmReferenceSet  = 999000581000001102
	TfReferenceSet   = 999000631000001100
	AmpReferenceSet  = 999000541000001108
	AmppReferenceSet = 999000551000001106
	VmpReferenceSet  = 999000561000001109
	VmppReferenceSet = 999000571000001104

	EnglishLanguageReferenceSet = 999000671000001103 // language reference set

	DoseFormReferenceSet               = 999000781000001107
	SugarFreeReferenceSet              = 999000601000001109
	GlutenFreeReferenceSet             = 999000611000001106
	PreservativeFreeReferenceSet       = 999000621000001102
	CombinationDrugVtm                 = 999000771000001105
	ChlorofluorocarbonFreeReferenceSet = 999000651000001105
	BlackTriangleReferenceSet          = 999000661000001108

	// Important relevant relationship types for dm+d concepts
	IsA                         = 116680003
	PendingMove                 = 900000000000492006
	HasActiveIngredient         = 127489000
	HasVmp                      = 10362601000001103
	HasAmp                      = 10362701000001108
	HasTradeFamilyGroup         = 9191701000001107
	HasSpecificActiveIngredient = 10362801000001104
	HasDispensedDoseForm        = 10362901000001105 // UK dm+d version of "HasDoseForm"
	HasDoseForm                 = 411116001         // Do not use - from International release - use dm+d relationship instead
	HasExcipient                = 8653101000001104
	PrescribingStatus           = 8940001000001105
	NonAvailabilityIndicator    = 8940601000001102
	LegalCategory               = 8941301000001102
	DiscontinuedIndicator       = 8941901000001101
	HasBasisOfStrength          = 10363001000001101
	HasUnitOfAdministration     = 13085501000001109
	HasNHSdmdBasisOfStrength    = 10363001000001101
)

// Prescribing status - descendants of 8940101000001106 -  https://termbrowser.nhs.uk/?perspective=full&conceptId1=8940101000001106&edition=uk-edition&release=v20181001&server=https://termbrowser.nhs.uk/sct-browser-api/snomed&langRefset=999001261000000100,999000691000001104
const (
	NeverValidToPrescribeAsVrp           = 12459601000001102
	NeverValidToPrescribeAsVmp           = 8940401000001100
	NotRecommendedToPrescribeAsVmp       = 8940501000001101
	InvalidAsPrescribableProduct         = 8940301000001108
	NotRecommendedBrandsNotBioequivalent = 9900001000001104
	NotRecommendedNoProductSpecification = 12468201000001102
	NotRecommendedPatientTraining        = 9900101000001103
	VmpValidAsPrescribableProduct        = 8940201000001104
	VrpValidAsPrescribableProduct        = 12223601000001104
)

func prescribingStatus(statusConceptID int64) (valid bool, recommended bool) {
	switch statusConceptID {
	case NeverValidToPrescribeAsVrp:
		return false, false
	case NeverValidToPrescribeAsVmp:
		return false, false
	case NotRecommendedToPrescribeAsVmp:
		return true, false
	case InvalidAsPrescribableProduct:
		return false, false
	case NotRecommendedBrandsNotBioequivalent:
		return true, false
	case NotRecommendedNoProductSpecification:
		return true, false
	case NotRecommendedPatientTraining:
		return true, false
	case VmpValidAsPrescribableProduct:
		return true, true
	case VrpValidAsPrescribableProduct:
		return true, true
	}
	return false, false
}

// Product is any dm+d product
type Product struct {
	svc *terminology.Svc
	*snomed.ExtendedConcept
}

// NewProduct creates a new UK dm+d product
func NewProduct(svc *terminology.Svc, ec *snomed.ExtendedConcept) Product {
	return Product{svc: svc, ExtendedConcept: ec}
}

func (p Product) String() string {
	if p.IsProduct() == false {
		return p.ExtendedConcept.String()
	}
	var sb strings.Builder
	sb.WriteString(p.GetPreferredDescription().GetTerm())
	sb.WriteString(" (")
	switch {
	case p.IsVMP():
		sb.WriteString("VMP")
	case p.IsVMPP():
		sb.WriteString("VMPP")
	case p.IsAMP():
		sb.WriteString("AMP")
	case p.IsAMPP():
		sb.WriteString("AMPP")
	case p.IsVTM():
		sb.WriteString("VTM")
	case p.IsTradeFamily():
		sb.WriteString("TF")
	default:
		sb.WriteString("Unknown")
	}
	sb.WriteString(")")
	return sb.String()
}

// IsInReferenceSet returns whether the product is in the specified reference set or not
func (p Product) IsInReferenceSet(refset int64) bool {
	if p.ExtendedConcept == nil {
		return false
	}
	for _, rs := range p.ExtendedConcept.GetConceptRefsets() {
		if rs == refset {
			return true
		}
	}
	return false
}

// GetRelationship returns the relationship of the specified type, indicated success or failure
func (p Product) GetRelationship(relationshipType int64) (*snomed.Relationship, bool) {
	if p.ExtendedConcept == nil {
		return nil, false
	}
	for _, rel := range p.ExtendedConcept.GetRelationships() {
		if rel.GetTypeId() == relationshipType {
			return rel, true
		}
	}
	return nil, false
}

// IsProduct confirms that this is a UK dm+d product.
func (p Product) IsProduct() bool {
	if p.ExtendedConcept == nil {
		return false
	}
	if p.Concept.Id == UKProduct {
		return true
	}
	for _, rs := range p.ExtendedConcept.GetRecursiveParentIds() {
		if rs == UKProduct {
			return true
		}
	}
	return false
}

// IsAMP returns whether this product is an Actual Medicinal Product (AMP)
func (p Product) IsAMP() bool {
	return p.IsInReferenceSet(AmpReferenceSet)
}

// IsAMPP returns whether this product is an Actual Medicinal Product Pack (AMPP)
func (p Product) IsAMPP() bool {
	return p.IsInReferenceSet(AmppReferenceSet)
}

// IsVTM returns whether this product is a virtual therapeutic moiety (VTM)
func (p Product) IsVTM() bool {
	return p.IsInReferenceSet(VtmReferenceSet)
}

// IsVMP returns whether this product is a virtual medicinal product (VMP)
func (p Product) IsVMP() bool {
	return p.IsInReferenceSet(VmpReferenceSet)
}

// IsVMPP returns whether this product is a virtual medicinal product pack (VMPP)
func (p Product) IsVMPP() bool {
	return p.IsInReferenceSet(VmppReferenceSet)
}

// IsTradeFamily returns whether this product is a trade family (TF)
func (p Product) IsTradeFamily() bool {
	return p.IsInReferenceSet(TfReferenceSet)
}

// IsSugarFree returns whether this product is flagged as sugar-free
func (p Product) IsSugarFree() bool {
	return p.IsInReferenceSet(SugarFreeReferenceSet)
}

// VMP is a virtual medicinal product.
// It is related to other dm+d models thusly:
// 	VMP -> IS-A -> VTM
// 	VMPP -> HAS VMP -> VMP
// 	AMP -> IS-A -> VMP
// 	VMP -> HAS_DISPENSED_DOSE_FORM -> [...]
// 	VMP -> HAS_SPECIFIC_ACTIVE_INGREDIENT -> [...]
type VMP struct {
	Product
}

// NewVmp creates a new VMP from the specified concept.
// It is an error to use a concept that is not a VMP
func NewVMP(svc *terminology.Svc, ec *snomed.ExtendedConcept) (*VMP, error) {
	product := NewProduct(svc, ec)
	if product.IsVMP() {
		return &VMP{Product: product}, nil
	}
	return nil, fmt.Errorf("%s is not a VMP", product)
}

// PrescribingStatus returns whether this VMP can be prescribed
func (vmp VMP) PrescribingStatus() (valid, recommended bool) {
	if vmp.IsVMP() == false {
		return
	}
	if rel, ok := vmp.GetRelationship(PrescribingStatus); ok {
		return prescribingStatus(rel.GetDestinationId())
	}
	return
}

// GetVTMs returns the VTM(s) for the given VMP
func (vmp VMP) GetVTMs() (result []int64) {
	for _, parent := range vmp.GetRecursiveParentIds() {
		items, err := vmp.svc.ComponentFromReferenceSet(VtmReferenceSet, parent)
		if err == nil && len(items) > 0 {
			result = append(result, parent)
		}
	}
	return
}

// VTM is a Virtual Therapeutic Moiety
// It will no have HAS_INGREDIENT relationship itself as part of dm+d, but
// the International release may include some with that type of relationship. That
// means derivation of ingredients, for dm+d products, is via VMP.
// See https://www.nhsbsa.nhs.uk/sites/default/files/2017-12/doc_UKTCSnomedCTUKDrugExtensionEditorialPolicy_Current-en-GB_GB1000001_20171227.pdf
type VTM struct {
	Product
}

// NewVTM creates a new VTM from the specified concept.
// It is an error to use a concept that is not a VTM
func NewVTM(svc *terminology.Svc, ec *snomed.ExtendedConcept) (*VTM, error) {
	product := NewProduct(svc, ec)
	if product.IsVTM() {
		return &VTM{Product: product}, nil
	}
	return nil, fmt.Errorf("%s is not a VMP", product)
}
