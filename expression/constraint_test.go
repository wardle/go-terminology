package expression

import (
	"testing"
)

var textExpressions = [...]string{
	"< 19829001 |Disorder of lung| : 116676008 |Associated morphology|= 79654002 |Edema|",
	"< 373873005 |Pharmaceutical / biologic product| : [1..3] 127489000 |Has active ingredient| = < 105590001 |Substance|",
	"< 105590001 |Substance| :R 127489000 |Has active ingredient| =249999999101 |TRIPHASIL tablet|",
	"< 19829001 |Disorder of lung|  AND <  301867009 |Edema of trunk|",
	"<  27658006 |Amoxicillin| : 411116001 |Has dose form|  =  <<  385055001 |Tablet dose form| , { 179999999100 |Has basis of strength|  = ( 219999999102 |Amoxicillin only| : 189999999103 |Has strength magnitude| >= #200,199999999101 |Has strength unit|  =  258684004 |mg| )}",
}

func TestExpressionConstraint(t *testing.T) {
}

func TestConstraintSelf(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()

	clinicalFinding, err := svc.Concept(404684003)
	if err != nil {
		t.Fatal(err)
	}
	clinicalFindingExp, err := NormalizeConcept(svc, clinicalFinding)
	if err != nil {
		t.Fatal(err)
	}

	// this constraint is satisfied only by the concept specified
	s1 := "404684003 |Clinical finding|"
	t1, err := ApplyConstraint(svc, clinicalFindingExp, s1)
	if err != nil {
		t.Fatal(err)
	}
	if t1 != true {
		t.Fatalf("clinical finding is a clinical finding. got: %v", t1)
	}
}

func TestConstraintSyntaxError(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()

	s1 := "wibble"
	t1, err := ApplyConstraint(svc, nil, s1)
	if err == nil {
		t.Fatalf("failed to identify syntax error. got %v", t1)
	}
}
