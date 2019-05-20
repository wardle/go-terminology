package expression

import (
	"testing"
)

// Test merging two groups that cannot be merged because there are no shared attributes
// see https://confluence.ihtsdotools.org/display/DOCTSG/12.4.10+Merging+Groups
func TestMergeNoMatchNoMerge(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	n := NewNormalizer(svc, nil)
	e1, err := ParseExpression("71388002 |Procedure| :363704007 |procedure site| = 421235005 |structure of femur|,363700003 |direct morphology| = 72704001 |fracture|")
	if err != nil {
		t.Error(err)
	}
	e2, err := ParseExpression("71388002 |Procedure| :260686004 |method| = 129371009 |fixation - action| ,424226004 |using device| = 31031000 |Orthopedic internal fixation system, device|")
	if err != nil {
		t.Error(err)
	}

	m1, merged1, err := n.mergeRefinements(e1.GetClause().GetRefinements(), e2.GetClause().GetRefinements())
	if err != nil {
		t.Error(err)
	}
	if m1 != false || merged1 != nil {
		t.Errorf("should not be able to merge:\n1:%v\n2:%v", Render(e1), Render(e2))
	}
}

// Test merging two groups with same attribute but values don't match and one isn't subsumed by other
// see https://confluence.ihtsdotools.org/display/DOCTSG/12.4.10+Merging+Groups
func TestMergeMatchNoMerge(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	n := NewNormalizer(svc, nil)
	e1, err := ParseExpression("64572001 |disease| :363698007 |finding site| = 62413002 |radius| ,	116676008 |associated morphology| = 72704001 |fracture|")
	if err != nil {
		t.Error(err)
	}

	e2, err := ParseExpression("64572001 |disease| :363698007 |finding site| = 87342007 |fibula|, 116676008 |associated morphology| = 72704001 |fracture|")
	if err != nil {
		t.Error(err)
	}

	m1, merged1, err := n.mergeRefinements(e1.GetClause().GetRefinements(), e2.GetClause().GetRefinements())
	if err != nil {
		t.Error(err)
	}
	if m1 != false || merged1 != nil {
		t.Errorf("should not be able to merge:\n1:%v\n2:%v", Render(e1), Render(e2))
	}
}

// Test merging two groups with same attribute and values mergable - ie one subsumed by other
// see https://confluence.ihtsdotools.org/display/DOCTSG/12.4.10+Merging+Groups
func TestMergeMatchAndMerge(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	n := NewNormalizer(svc, nil)
	e1, err := ParseExpression("64572001 |disease| :363698007 |finding site| = 62413002 |radius| ,116676008 |associated morphology| = 72704001 |fracture|")
	if err != nil {
		t.Error(err)
	}

	e2, err := ParseExpression("64572001 |disease| :363698007 |finding site| = 75129005 |distal radius|,116676008 |associated morphology| = 72704001 |fracture|")
	if err != nil {
		t.Error(err)
	}

	m1, merged1, err := n.mergeRefinements(e1.GetClause().GetRefinements(), e2.GetClause().GetRefinements())
	if err != nil {
		t.Error(err)
	}
	if m1 == false || merged1 == nil {
		t.Errorf("should be able to merge:\n1:%v\n2:%v", Render(e1), Render(e2))
	}
}
