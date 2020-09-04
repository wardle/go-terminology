// Package dmd provides functionality to process and understand data from the
// UK Dictionary of Medicines and Devices (dm+d).
// This is, by definition, a UK-only module and will not give expected results for
// drugs outside of the UK Product root.
//
//
// See https://www.nhsbsa.nhs.uk/sites/default/files/2018-10/doc_SnomedCTUKDrugExtensionModel%20-%20v1.0.pdf
//
// The dm+d model consists of the following components:
// VTM
// VMP
// VMPP
// TF
// AMP
// AMPP
//
// The relationships between these components are:
//
// VMP <<- IS_A -> VTM
// VMP <<- HAS_SPECIFIC_ACTIVE_INGREDIENT ->> SUBSTANCE
// VMP <<- HAS_DISPENSED_DOSE_FORM ->> QUALIFIER
// VMPP <<- HAS_VMP -> VMP
// AMPP <<- IS_A -> VMPP
// AMPP <<- HAS_AMP -> AMP
// AMP <<- IS_A -> VMP
// AMP <<- IS_A -> TF
// AMP <<- HAS_EXCIPIENT ->> QUALIFIER
// TF <<- HAS_TRADE_FAMILY_GROUP ->> QUALIFIER
//
// Cardinality rules are: (see https://www.nhsbsa.nhs.uk/sites/default/files/2017-02/Technical_Specification_of_data_files_R2_v3.1_May_2015.pdf)
// The SNOMED dm+d data file documents the cardinality rules for AMP<->TF (https://www.nhsbsa.nhs.uk/sites/default/files/2017-04/doc_UKTCSnomedCTUKDrugExtensionEditorialPolicy_Current-en-GB_GB1000001_v7_0.pdf)
//
package dmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
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

	// Language reference sets
	NhsDmdRealmLanguageReferenceSet                = 999000671000001103 // language reference set - dm+d
	NhsRealmPharmacyLanguageReferenceSet           = 999000691000001104
	NhsRealmClinicalLanguageReferenceSet           = 999001261000000100
	NhsEPrescribingRouteAdministrationReferenceSet = 999000051000001100
	// Other reference sets
	DoseFormReferenceSet               = 999000781000001107
	SugarFreeReferenceSet              = 999000601000001109
	GlutenFreeReferenceSet             = 999000611000001106
	PreservativeFreeReferenceSet       = 999000621000001102
	CombinationDrugVtm                 = 999000771000001105
	ChlorofluorocarbonFreeReferenceSet = 999000651000001105
	BlackTriangleReferenceSet          = 999000661000001108

	// Important relevant relationship types for dm+d concepts
	IsA                                  = 116680003
	PendingMove                          = 900000000000492006
	HasActiveIngredient                  = 127489000
	HasVmp                               = 10362601000001103
	HasAmp                               = 10362701000001108
	HasTradeFamilyGroup                  = 9191701000001107
	HasSpecificActiveIngredient          = 10362801000001104
	HasDispensedDoseForm                 = 10362901000001105 // UK dm+d version of "HasDoseForm"
	HasDoseForm                          = 411116001         // Do not use - from International release - use dm+d relationship instead
	HasExcipient                         = 8653101000001104
	PrescribingStatus                    = 8940001000001105
	NonAvailabilityIndicator             = 8940601000001102
	LegalCategory                        = 8941301000001102
	DiscontinuedIndicator                = 8941901000001101
	HasBasisOfStrength                   = 10363001000001101
	HasUnitOfAdministration              = 13085501000001109
	HasUnitOfPresentation                = 763032000
	HasNHSdmdBasisOfStrength             = 10363001000001101
	HasNHSControlledDrugCategory         = 13089101000001102
	HasVMPNonAvailabilityIndicator       = 8940601000001102
	VMPPrescribingStatus                 = 8940001000001105
	HasNHSdmdVmpRouteOfAdministration    = 13088401000001104
	HasNHSdmdVmpOntologyFormAndRoute     = 13088501000001100
	HasPresentationStrengthNumerator     = 732944001
	HasPresentationStrengthDenominator   = 732946004
	HasPresentationStrengthNumeratorUnit = 732945000
)

// Prescribing status - descendants of 8940101000001106 -  https://termbrowser.nhs.uk/?perspective=full&conceptId1=8940101000001106&edition=uk-edition&release=v20181001&server=https://termbrowser.nhs.uk/sct-browser-api/snomed&langRefset=999001261000000100,999000691000001104
const (
	CautionAMPLevelPrescribingAdvised    = 13291401000001100
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
	case CautionAMPLevelPrescribingAdvised:
		return true, false
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
	if !p.IsProduct() {
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

// PrettyPrint provides a pretty formatted version of this product.
func (p Product) PrettyPrint() string {
	s, _ := json.MarshalIndent(p, "", "\t")
	return string(s)
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

// Relationships returns the active (parent) relationships of the specified type.
func (p Product) Relationships(relationshipType int64) (result []int64) {
	if p.ExtendedConcept == nil {
		return
	}
	for _, rel := range p.ExtendedConcept.GetRelationships() {
		if rel.Active && rel.GetTypeId() == relationshipType {
			result = append(result, rel.DestinationId)
		}
	}
	return
}

// Relationship returns the single active relationship of the specified type.
func (p Product) Relationship(relationshipType int64) int64 {
	for _, rel := range p.ExtendedConcept.GetRelationships() {
		if rel.Active && rel.TypeId == relationshipType {
			return rel.DestinationId
		}
	}
	return 0
}

// IsProduct confirms that this is a UK dm+d product.
func (p Product) IsProduct() bool {
	if p.ExtendedConcept == nil {
		return false
	}
	if p.Concept.Id == UKProduct {
		return true
	}
	for _, rs := range p.ExtendedConcept.GetAllParentIds() {
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

// NewVMP creates a new VMP from the specified concept.
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
	return prescribingStatus(vmp.Relationship(PrescribingStatus))
}

// HasNHSControlledDrugCategory returns the controlled drug category for this VMP
func (vmp VMP) HasNHSControlledDrugCategory() int64 {
	return vmp.Relationship(HasNHSControlledDrugCategory)
}

// HasNHSdmdVmpRouteAdministration returns the route of administration for this VMP
func (vmp VMP) HasNHSdmdVmpRouteAdministration() int64 {
	return vmp.Relationship(HasNHSdmdVmpRouteOfAdministration)
}

// HasNHSdmdVmpOntologyFormAndRoute returns the ontology form and route for this VMP
func (vmp VMP) HasNHSdmdVmpOntologyFormAndRoute() int64 {
	return vmp.Relationship(HasNHSdmdVmpOntologyFormAndRoute)
}

// HasPresentationStrengthDenominator returns the presentation strength denominator
func (vmp VMP) HasPresentationStrengthDenominator() []int64 {
	return vmp.Relationships(HasPresentationStrengthDenominator)
}

// HasPresentationStrengthNumerator returns the presentation strength numerator
func (vmp VMP) HasPresentationStrengthNumerator() []int64 {
	return vmp.Relationships(HasPresentationStrengthNumerator)
}

// HasPresentationStrengthNumeratorUnit returns the presentation strength numerator unit
func (vmp VMP) HasPresentationStrengthNumeratorUnit() []int64 {
	return vmp.Relationships(HasPresentationStrengthNumeratorUnit)
}

// HasBasisOfStrength  returns NHS dm+d determined basis of strength
func (vmp VMP) HasBasisOfStrength() []int64 {
	return vmp.Relationships(HasNHSdmdBasisOfStrength) // TODO: see 377269004 - only non-nhs has basis of strength correct!!
}

// VTMs returns the VTM(s) for the given VMP
// 	VMP -> IS-A -> VTM
func (vmp VMP) VTMs() (result []int64) {
	for _, parent := range vmp.GetAllParentIds() {
		items, err := vmp.svc.ComponentFromReferenceSet(VtmReferenceSet, parent)
		if err == nil && len(items) > 0 {
			result = append(result, parent)
		}
	}
	return
}

// AMPs returns the AMP(s) for the given VMP
// AMP -> IS-A -> VMP
// Note, this method does not include AMPs for child VMPs of this VMP.
func (vmp VMP) AMPs(ctx context.Context) (result []int64, err error) {
	children, err := vmp.svc.Children(vmp.GetConcept().GetId())
	if err != nil {
		return
	}
	for _, child := range children {
		isAmp, err := vmp.svc.IsInReferenceSet(child, AmpReferenceSet)
		if err != nil {
			return nil, err
		}
		if isAmp {
			result = append(result, child)
		}
	}
	return
}

// AllAMPs return all (recursive) AMPs for this VMP
// Because a VMP can have children that are VMPs, this returns AMPs for this VMP
// *and* child VMPs.
// In most use-cases, use AMPs() instead.
func (vmp VMP) AllAMPs(ctx context.Context) (result []int64, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for child := range vmp.svc.StreamAllChildrenIDs(ctx, vmp.GetConcept().GetId(), 50000) {
		if child.Err != nil {
			err = child.Err
		}
		isAMP, err := vmp.svc.IsInReferenceSet(child.ID, AmpReferenceSet)
		if err != nil {
			return nil, err
		}
		if isAMP {
			result = append(result, child.ID)
		}
	}
	return
}

// VMPPs returns the VMPPs for this VMP
func (vmp VMP) VMPPs() (result []int64, err error) {
	rels, err := vmp.svc.ChildRelationships(vmp.ExtendedConcept.Concept.Id)
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		if rel.Active && rel.TypeId == HasVmp {
			result = append(result, rel.SourceId)
		}
	}
	return
}

// SpecificActiveIngredients returns the ingredients for this VMP
// 	VMP -> HAS_SPECIFIC_ACTIVE_INGREDIENT -> [...]
func (vmp VMP) SpecificActiveIngredients() (result []int64) {
	return vmp.Relationships(HasSpecificActiveIngredient)
}

// DisposedDoseForm returns the disposed dose form of this VMP
// 	VMP -> HAS_DISPENSED_DOSE_FORM -> [...]
func (vmp VMP) DisposedDoseForm() int64 {
	return vmp.Relationship(HasDispensedDoseForm)
}

// VTM is a Virtual Therapeutic Moiety
// It will no have HAS_INGREDIENT relationship itself as part of dm+d, but
// the International release may include some with that type of relationship. That
// means derivation of ingredients, for dm+d products, is via VMP.
// See https://www.nhsbsa.nhs.uk/sites/default/files/2017-12/doc_UKTCSnomedCTUKDrugExtensionEditorialPolicy_Current-en-GB_GB1000001_20171227.pdf
//
// It is related to other components in dm+d thusly:
// VMP - IS_A - VTM
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

// AllVMPs returns the VMPs for this VTM
// Recursive children will include VMPS and AMPs. Direct children will include other VTMs.
// ie. you can have VMP -IS-A- VMP -IS-A- VMP -IS-A- VTM
// We therefore need to recursively search children, and filter both.
// // VMP - IS_A - VTM
// // VMP - IS_A - VMP - IS_A - VTM
// // VMP - IS_A - VTM - IS_A - VTM
func (vtm VTM) AllVMPs(ctx context.Context) (result []int64, err error) {
	children := vtm.svc.StreamAllChildrenIDs(ctx, vtm.Concept.Id, 50000)
	for child := range children {
		if child.Err != nil {
			err = child.Err
			return
		}
		isVmp, err := vtm.svc.IsInReferenceSet(child.ID, VmpReferenceSet)
		if err != nil {
			return nil, err
		}
		if isVmp {
			result = append(result, child.ID)
		}
	}
	return
}

// VMPs returns only direct IS-A VMPs for this VTM.
func (vtm VTM) VMPs(ctx context.Context) (result []int64, err error) {
	children, err := vtm.svc.Children(vtm.Concept.Id)
	if err != nil {
		return
	}
	for _, child := range children {
		isVmp, err := vtm.svc.IsInReferenceSet(child, VmpReferenceSet)
		if err != nil {
			return nil, err
		}
		if isVmp {
			result = append(result, child)
		}
	}
	return
}

// SpecificActiveIngredients returns the specific active ingredients for this VTM
// This simply returns the active ingredients via its (direct) VMPs.
func (vtm VTM) SpecificActiveIngredients(ctx context.Context, tags []language.Tag) ([]int64, error) {
	ingredients := make(map[int64]struct{})
	vmps, err := vtm.VMPs(ctx)
	if err != nil {
		return nil, err
	}
	for _, vmpID := range vmps {
		ec, err := vtm.svc.ExtendedConcept(vmpID, tags)
		if err != nil {
			return nil, err
		}
		vmp, err := NewVMP(vtm.svc, ec)
		if err != nil {
			return nil, err
		}
		for _, ingredient := range vmp.SpecificActiveIngredients() {
			ingredients[ingredient] = struct{}{}
		}
	}
	result := make([]int64, 0, len(ingredients))
	for ingredient := range ingredients {
		result = append(result, ingredient)
	}
	return result, nil
}

// AMP is an actual medicinal product
// It is related to other components thusly:
// AMPP - HAS_AMP - AMP
// AMP - IS_A - VMP
// AMP - IS_A - TF
// AMP - HAS_EXCIPIENT - QUALIFIER
type AMP struct {
	Product
}

// NewAMP creates a new AMP from the specified concept.
// It is an error to use a concept that is not a AMP
func NewAMP(svc *terminology.Svc, ec *snomed.ExtendedConcept) (*AMP, error) {
	product := NewProduct(svc, ec)
	if product.IsAMP() {
		return &AMP{Product: product}, nil
	}
	return nil, fmt.Errorf("%s is not an AMP", product)
}

// AMPPs returns the AMPPs for the given AMP
// AMPP - HAS_AMP - AMP
func (amp AMP) AMPPs() (result []int64, err error) {
	rels, err := amp.svc.ChildRelationships(amp.ExtendedConcept.Concept.Id)
	if err != nil {
		return nil, err
	}
	for _, rel := range rels {
		if rel.Active && rel.TypeId == HasAmp {
			result = append(result, rel.SourceId)
		}
	}
	return
}

// VMP returns the VMP for the given AMP
// AMP - IS_A - VMP
func (amp AMP) VMP() (vmp int64, err error) {
	var isVmp bool
	for _, parent := range amp.DirectParentIds {
		if isVmp, err = amp.svc.IsInReferenceSet(parent, VmpReferenceSet); err != nil {
			return
		}
		if isVmp {
			return parent, nil
		}
	}
	return
}

// TF returns the trade family for this AMP
// AMP - IS_A - TF
func (amp AMP) TF() (tf int64, err error) {
	var isTf bool
	for _, parent := range amp.DirectParentIds {
		if isTf, err = amp.svc.IsInReferenceSet(parent, TfReferenceSet); err != nil {
			return
		}
		if isTf {
			return parent, nil
		}
	}
	return
}

// Excipients returns the excipients for this AMP
/// AMP - HAS_EXCIPIENT - QUALIFIER
// If there are no excipients, dm+d explicitly returns "8653301000001102 = Excipient Not Declared"
func (amp AMP) Excipients() (excipients []int64) {
	return amp.Relationships(HasExcipient)
}

// ControlledDrugStatus returns the NHS controlled drug category
// which is a child of "13089201000001109".
func (amp AMP) ControlledDrugStatus() int64 {
	return amp.Relationship(HasNHSControlledDrugCategory)
}

// VMPNonAvailabilityIndicator returns whether the AMP's VMP is available, or not
func (amp AMP) VMPNonAvailabilityIndicator() int64 {
	return amp.Relationship(HasVMPNonAvailabilityIndicator)
}

// VMPPrescribingStatus returns whether the AMP's VMP is prescribable.
func (amp AMP) VMPPrescribingStatus() int64 {
	return amp.Relationship(VMPPrescribingStatus)
}

// TF is a trade family and is related to other components in dm+d thusly:
// AMP <<- IS_A -> TF
// TF <<- HAS_TRADE_FAMILY_GROUP ->> QUALIFIER
type TF struct {
	Product
}

// NewTF creates a new TF from the specified concept.
// It is an error to use a concept that is not a TF
func NewTF(svc *terminology.Svc, ec *snomed.ExtendedConcept) (*TF, error) {
	product := NewProduct(svc, ec)
	if product.IsTradeFamily() {
		return &TF{Product: product}, nil
	}
	return nil, fmt.Errorf("%s is not a TF", product)
}

// AMPs returns the AMPs for this TF
// AMP <<- IS_A -> TF
func (tf TF) AMPs(ctx context.Context) (result []int64, err error) {
	children := tf.svc.StreamAllChildrenIDs(ctx, tf.Concept.Id, 50000)
	for child := range children {
		if child.Err != nil {
			err = child.Err
			return
		}
		isAmp, err := tf.svc.IsInReferenceSet(child.ID, AmpReferenceSet)
		if err != nil {
			return nil, err
		}
		if isAmp {
			result = append(result, child.ID)
		}
	}
	return
}

// TradeFamilyGroup returns the trade family group for this trade family.
func (tf TF) TradeFamilyGroup() int64 {
	rels := tf.GetRelationships()
	for _, rel := range rels {
		if rel.Active && rel.TypeId == HasTradeFamilyGroup {
			return rel.DestinationId
		}
	}
	return 0
}
