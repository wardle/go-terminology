package expression

import (
	"testing"
)

// Test normalizing a simple expression
// e.g. See https://confluence.ihtsdotools.org/display/DOCTSG/12.3.13+Normal+Form+of+a+Simple+Expression?src=sidebarhttps://confluence.ihtsdotools.org/display/DOCTSG/12.3.13+Normal+Form+of+a+Simple+Expression?src=sidebar
func TestNormalize1(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	s := `12676007 |fracture of radius| + 397181002 |open fracture| :
	272741003 |laterality| = 7771000 |left|,
	42752001 |due to| = 297186008 |motorcycle accident|`
	e1, err := ParseExpression(s)
	if err != nil {
		t.Error(err)
	}
	normalizer := NewNormalizer(svc)
	focusConcepts, err := normalizer.normalizedFocusConcepts(e1)
	if err != nil {
		t.Error(err)
	}
	if len(focusConcepts) != 2 {
		t.Errorf("wrong number of focus concepts. expected 2, got: %d (%v)", len(focusConcepts), focusConcepts)
	}
	// we expect the two focus concepts to be normalised like this:
	expected1, err := ParseExpression(`404684003|Clinical finding|:
					116676008|Associated morphology|=72704001|Fracture|,
					363698007|Finding site|=62413002|Bone structure of radius|`)
	if err != nil {
		t.Error(err)
	}
	if !Equal(expected1, focusConcepts[0]) && !Equal(expected1, focusConcepts[1]) {
		t.Errorf("did not correctly normalise 'fracture of radius'. expected:\n%sgot:\n", Render(expected1))
		for _, fc := range focusConcepts {
			t.Error(Render(fc))
		}
	}
}

// Test merging two groups that cannot be merged because there are no shared attributes
// see https://confluence.ihtsdotools.org/display/DOCTSG/12.4.10+Merging+Groups
func TestMergeNoMatchNoMerge(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	n := NewNormalizer(svc)
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
	n := NewNormalizer(svc)
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
	n := NewNormalizer(svc)
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
