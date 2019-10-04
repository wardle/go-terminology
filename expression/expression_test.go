package expression

import (
	"fmt"
	"os"
	"testing"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
)

const (
	dbFilename = "../snomed.db" // real, live database
)

var (
	languageTags = []language.Tag{terminology.BritishEnglish.Tag()}
)

var etests = []struct {
	name                string
	expression          string
	numFocusConcepts    int
	numRefinements      int
	numRefinementGroups int
	totalRefinements    int
	f                   func(e *snomed.Expression) error
}{
	{
		"Simple",
		"73211009 |Diabetes mellitus|",
		1, 0, 0, 0,
		func(e *snomed.Expression) error {
			if e.GetClause().GetFocusConcepts()[0].ConceptId != 73211009 {
				return fmt.Errorf("did not identify diabetes mellitus. got %v instead", e)
			}
			return nil
		},
	},
	{
		"Simple refinement",
		"83152002 |oophorectomy|: 405815000|procedure device| = 122456005 |laser device|",
		1, 1, 0, 1,
		func(e *snomed.Expression) error {
			devices := e.GetClause().GetRefinementOfType(405815000)
			if len(devices) != 1 {
				return fmt.Errorf("incorrect number of devices. expected 1. got:%d (%v)", len(devices), e.GetClause())
			}
			if devices[0].GetConceptValue().GetConceptId() != 122456005 {
				return fmt.Errorf("did not get laser device. expected:122456005 got: %v", devices[0])
			}
			return nil
		},
	},
	{
		"Multiple attributes",
		"71388002 |procedure|:	405815000|procedure device| = 122456005 |laser device|, 260686004 |method| = 129304002 |excision - action|,405813007 |procedure site - direct| = 15497006 |ovarian structure|",
		1, 3, 0, 3, nil,
	},
	{
		"Conjoined expression",
		"119189000 |ulna part| + 312845000 |epiphysis of upper limb|:272741003 |laterality| = 7771000 |left|",
		2, 1, 0, 1, nil,
	},
	{
		"Complex expression",
		"3415004 |cyanosis| + 363696006 |neonatal cardiovascular disorder|:246454002 |occurrence| = 255407002 |neonatal|,	363698007 |finding site| = 113257007 |structure of cardiovascular system|",
		2, 2, 0, 2, nil,
	},
	{
		"Attribute group 1",
		"71388002 |procedure|:{ 260686004 |method| = 129304002 |excision - action|,405813007 |procedure site - direct| = 15497006 |ovarian structure|} { 260686004 |method| = 129304002 |excision - action|,405813007 |procedure site - direct| = 31435000 |fallopian tube structure|}",
		1, 0, 2, 4, nil,
	},
	{
		"Attribute group 2",
		"71388002 |procedure|:{ 260686004 |method| = 129304002 |excision - action|,405813007 |procedure site - direct| = 20837000 |structure of right ovary|,424226004 |using device| = 122456005 |laser device|} {260686004 |method| = 261519002 |diathermy excision - action|,405813007 |procedure site - direct| = 113293009 |structure of left fallopian tube|}",
		1, 0, 2, 5, nil,
	},
	{
		"Nested expression",
		"397956004 |prosthetic arthroplasty of the hip|:363704007 |procedure site| = (24136001 |hip joint structure|:272741003 |laterality| = 7771000 |left|)",
		1, 1, 0, 1, nil,
	},
	{
		"Concrete value",
		"27658006 |amoxicillin |:411116001 |has dose form| = 385049006 |capsule|,{ 127489000 |has active ingredient| = 372687004 |amoxicillin|,111115|has basis of strength| = (111115 |amoxicillin only|:111115|strength magnitude| = #500, 111115|strength unit| = 258684004 |mg|)}",
		1, 1, 1, 3,
		func(e *snomed.Expression) error {
			if e.GetDefinitionStatus() != snomed.Expression_EQUIVALENT_TO {
				return fmt.Errorf("Failed to determine appropriate definition status. Got %v", e.GetDefinitionStatus())
			}
			return nil
		},
	},
	{
		"Test Equivalent To",
		"=== 46866001 |fracture of lower limb| + 428881005 |injury of tibia|: 116676008 |associated morphology| = 72704001 |fracture|, 363698007 |finding site| = 12611008 |bone structure of tibia|",
		2, 2, 0, 2,
		func(e *snomed.Expression) error {
			if e.GetDefinitionStatus() != snomed.Expression_EQUIVALENT_TO {
				return fmt.Errorf("Failed to determine appropriate definition status. Got %v", e.GetDefinitionStatus())
			}
			return nil
		},
	},
	{
		"Test Subtype of",
		"<<< 73211009 |diabetes mellitus|: 363698007 |finding site| = 113331007 |endocrine system|",
		1, 1, 0, 1,
		func(e *snomed.Expression) error {
			if e.GetDefinitionStatus() != snomed.Expression_SUBTYPE_OF {
				return fmt.Errorf("Failed to determine appropriate definition status. Got %v", e.GetDefinitionStatus())
			}
			return nil
		},
	},
}

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

func TestBadParse1(t *testing.T) {
	exp, err := Parse("wibble")
	if exp != nil && err == nil {
		t.Fatalf("parsed a bad request and did not flag an error")
	}
	if perr, ok := err.(*ParseError); ok {
		if perr.Line != 1 || perr.Column != 0 || perr.OffendingToken != "w" {
			t.Fatalf("incorrect parse error. got: %v", perr)
		}
	} else {
		t.Fatalf("syntax error in parsing did not result in a structured parse error. got: %v", err)
	}
}

func TestBadParse2(t *testing.T) {
	_, err := Parse(`404684003|Clinical finding|:
					116676008|Associated morphology|=72704001|Fracture|,
					363698007|Finding site|=62413002|Bone structure of radius`)
	if err == nil {
		t.Error("did not correctly identify malformed expression")
	}
	if _, ok := err.(*ParseError); !ok {
		t.Fatalf("missing structured parse error. got %v", err)
	}
}

func TestPostcoordinationTests(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tags, _, _ := language.ParseAcceptLanguage("en-GB") // TODO(mw): better language support

	ms, err := svc.ExtendedConcept(24700007, tags)
	if err != nil {
		t.Error(err)
	}
	e := CreateSimpleExpression(ms.Concept)
	if e.IsPrecoordinated() == false {
		t.Errorf("Multiple sclerosis not correctly identified as a pre-coordinated expression")
	}
	if e.IsPostcoordinated() == true {
		t.Errorf("Multiple sclerosis incorrectly identified as a post-coordinated expression")
	}
}

func TestExpressions(t *testing.T) {
	for _, test := range etests {
		e, err := Parse(test.expression)
		if err != nil {
			t.Errorf("%s: %s", test.name, err)
		}
		if test.numFocusConcepts != len(e.GetClause().GetFocusConcepts()) {
			t.Errorf("%s: Invalid number of focus concepts. Expected %d, got %v\n", test.name, test.numFocusConcepts, e.GetClause().GetFocusConcepts())
		}
		if test.numRefinementGroups != len(e.GetClause().GetRefinementGroups()) {
			t.Errorf("%s: Invalid number of refinement groups. Expected %d, got %v\n", test.name, test.numRefinementGroups, e.GetClause().GetRefinementGroups())
		}
		if test.numRefinements != len(e.GetClause().GetRefinements()) {
			t.Errorf("%s: Invalid number of refinements. Expected %d, got %v\n", test.name, test.numRefinements, e.GetClause().GetRefinements())
		}
		total := len(e.GetClause().GetRefinements())
		for _, g := range e.GetClause().GetRefinementGroups() {
			total += len(g.GetRefinements())
		}
		if test.totalRefinements != total {
			t.Errorf("%s: Invalid total number of refinements. Expected %d, got %d\n", test.name, test.totalRefinements, total)
		}
		if test.f != nil {
			if err := test.f(e); err != nil {
				t.Errorf("%s: %s", test.name, err)
			}
		}
	}
}

// TestPrimitive1 tests deriving a primitive as per information on:
// https://confluence.ihtsdotools.org/display/DOCTSG/12.3.8+Normal+Form+of+a+Fully-Defined+Concept+with+No+Intermediate+Primitives
func TestPrimitive1(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	fractureFemur, err := svc.Concept(71620000)
	if err != nil {
		t.Error(err)
	}
	if fractureFemur.IsPrimitive() == true {
		t.Fatalf("fracture of femur incorrectly identified as a primitive type")
	}
	if fractureFemur.IsSufficientlyDefined() == false {
		t.Errorf("fracture of femur not correctly identified as a sufficiently defined type")
	}
	primitive, err := svc.Primitive(fractureFemur)
	if primitive.Id != 64572001 {
		t.Fatalf("primitive supertype of fracture of femur not correctly identified as 'disease'")
	}
	normalized, err := NormalizeConcept(svc, fractureFemur)
	if err != nil {
		t.Fatal(err)
	}
	if normalized.Clause.FocusConcepts[0].ConceptId != 64572001 {
		t.Fatalf("Did not correctly normalize fracture of femur to disease")
	}
	if len(normalized.Clause.Refinements) != 2 {
		t.Fatalf("Fracture of femur should have two refinements. Found: %v", normalized.Clause.Refinements)
	}
	refinements := make(map[int64]int64)
	for _, r := range normalized.Clause.Refinements {
		refinements[r.GetRefinementConcept().ConceptId] = r.GetConceptValue().ConceptId
	}
	if r, ok := refinements[363698007]; !ok || r != 71341001 {
		t.Fatalf("fracture of femur not correctly normalised to include finding site attribute 'bone structure of femur'")
	}
	if r, ok := refinements[116676008]; !ok || r != 72704001 {
		t.Fatalf("fracture of femur not correctly normalised to include morphology attribute 'fracture'")
	}
}

// TestPrimitive2 tests normal forms of a primitive concept
// See https://confluence.ihtsdotools.org/display/DOCTSG/12.3.9+Normal+Forms+of+a+Primitive+Concept
// Although this page is out of date and has the wrong defined relationships!
func TestPrimitive2(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	asthma, err := svc.Concept(195967001)
	if err != nil {
		t.Fatal(err)
	}
	if asthma.IsPrimitive() == false {
		t.Fatal("asthmia incorrectly defined as a fully defined concept, rather than primitive")
	}
	normalized, err := NormalizeConcept(svc, asthma)
	if normalized.Clause.FocusConcepts[0].ConceptId != 195967001 {
		t.Fatalf("focus concept for normal form of asthma incorrect. expected 195967001, was: %v", normalized.Clause.FocusConcepts)
	}
	if len(normalized.Clause.Refinements) != 1 {
		t.Fatalf("incorrect number of refinements for normal form of asthma. was %v", normalized.Clause.Refinements)
	}
	r := normalized.Clause.Refinements[0]
	if r.RefinementConcept.ConceptId != 363698007 || r.GetConceptValue().ConceptId != 89187006 {
		t.Fatalf("Asthma not correctly identified as a disease of the airways. was : %v", r)
	}
}

// TestPrimitive3 tests normalization of a concept with an intermediate primitive
// Unfortunately, the documentation here (https://confluence.ihtsdotools.org/display/DOCTSG/12.3.10+Normal+Form+of+a+Fully-Defined+Concept+with+an+Intermediate+Primitive)
// is out-of-date. TODO: find concept with intermediate primitives
func TestPrimitive3(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	appendicitis, err := svc.Concept(74400008)
	if err != nil {
		t.Fatal(err)
	}
	if appendicitis.IsPrimitive() {
		t.Fatal("appendicitis incorrectly flagged as primitive")
	}
	normalized, err := NormalizeConcept(svc, appendicitis)
	if err != nil {
		t.Error(err)
	}
	if normalized.IsPostcoordinated() == false {
		t.Errorf("normalization failed to return post coordinated expression")
	}
}

func printExpression(exp *snomed.Expression) {
	for _, c := range exp.GetClause().GetFocusConcepts() {
		fmt.Printf("focus concept:%v\n", c)
	}
	for _, r := range exp.GetClause().GetRefinements() {
		fmt.Printf("refinement: %v = %v\n", r.GetRefinementConcept(), r.GetValue())
	}
	for i, g := range exp.GetClause().GetRefinementGroups() {
		for _, r := range g.GetRefinements() {
			fmt.Printf("group %d: refinement: %v = %v\n", i, r.GetRefinementConcept(), r.GetValue())
		}
	}
}

var roundtrips = []struct {
	expression string
	canonical  string
}{
	{
		`===  46866001 |Fracture of lower limb| +  428881005 |Injury of tibia|:
        116676008 |Associated morphology| =  72704001 |Fracture|,
		363698007 |Finding site| =  12611008 |Bone structure of tibia|`,
		"428881005+46866001:116676008=72704001,363698007=12611008",
	},
	{
		"<<<  73211009 |Diabetes mellitus|:  363698007 |Finding site| =  113331007 |Endocrine system|",
		"<<<73211009:363698007=113331007",
	},
	{
		`373873005 |Pharmaceutical / biologic product|:
		411116001 |Has dose form| = ( 421720008 |Spray dose form| +  7946007 |Drug suspension|)`,
		"373873005:411116001=(421720008+7946007)",
	},
	{
		`  27658006 |Amoxicillin|:
		411116001 |Has dose form| =  385049006 |Capsule|,
	   {  127489000 |Has active ingredient| =  372687004 |Amoxicillin|,
		  179999999100 |Has basis of strength| = ( 219999999102 |Amoxicillin only|:
		  189999999103 |Has strength magnitude| = #500,  199999999101 |Has strength unit| =  258684004 |mg|)}`,
		"27658006:411116001=385049006,{127489000=372687004,179999999100=(219999999102:189999999103=#500,199999999101=258684004)}",
	},
	{
		" <<<  73211009 |Diabetes mellitus|:  363698007 |Finding site| =  113331007 |Endocrine system|",
		"<<<73211009:363698007=113331007",
	},
}

func TestExpressionConstraint(t *testing.T) {
	s1 := "< 19829001 |Disorder of lung| : 116676008 |Associated morphology|= 79654002 |Edema|"
	ParseExpressionConstraint(s1)
	s2 := "< 373873005 |Pharmaceutical / biologic product| : [1..3] 127489000 |Has active ingredient| = < 105590001 |Substance|"
	ParseExpressionConstraint(s2)
	s3 := "< 105590001 |Substance| :R 127489000 |Has active ingredient| =249999999101 |TRIPHASIL tablet|"
	ParseExpressionConstraint(s3)
}
