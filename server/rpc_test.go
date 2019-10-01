package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"

	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	context "golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	dbFilename = "../snomed.db" // real, live database
	port       = ":50080"
	lang       = "en-GB"
)

func Server(svc *terminology.Svc) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	tags, _, _ := language.ParseAcceptLanguage(lang)
	snomed.RegisterSnomedCTServer(s, &coreServer{svc: svc, lang: tags})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
func TestMain(m *testing.M) {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		log.Printf("Skipping live tests in the absence of a live database %s", dbFilename)
		os.Exit(0)
	}
	svc, err := terminology.NewService(dbFilename, true)
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	go Server(svc)
	os.Exit(m.Run())
}

func TestRpcClient(t *testing.T) {
	address := fmt.Sprintf("localhost%s", port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := snomed.NewSnomedCTClient(conn)
	header := metadata.New(map[string]string{"accept-language": "en-GB"})
	ctx := metadata.NewOutgoingContext(context.Background(), header)

	// Test GetConcept as a subtest
	t.Run("GetConcept", func(t *testing.T) {
		c, err := c.GetConcept(ctx, &snomed.SctID{Identifier: 24700007})
		if err != nil {
			t.Fatal(err)
		}
		if c.GetId() != 24700007 {
			t.Error("Expected '24700007', got ", c.GetId())
		}
	})
	t.Run("GetExtendedConcept", func(t *testing.T) {
		header := metadata.New(map[string]string{"accept-language": "en-US"})
		usCtx := metadata.NewOutgoingContext(context.Background(), header)

		c, err := c.GetExtendedConcept(usCtx, &snomed.SctID{Identifier: 80146002})
		if err != nil {
			t.Fatal(err)
		}
		if c.GetPreferredDescription().Term != "Appendectomy" {
			t.Errorf("Appendicectomy not correctly localised for language preferences. Got %s", c.GetPreferredDescription().Term)
		}
	})
	t.Run("Map", func(t *testing.T) {
		// test translating MS into emergency care reference set - should give MS
		t1, err := c.Map(ctx, &snomed.TranslateToRequest{ConceptId: 24700007, TargetId: 991411000000109})
		if err != nil {
			t.Fatal(err)
		}
		if t1.Translations[0].GetConcept().GetConceptId() != 24700007 {
			t.Errorf("failed to find multiple sclerosis in the emergency care reference set. found: %v", t1)
		}
		// test translating ADEM into emergency care reference set - should get encephalitis (45170000)
		t2, err := c.Map(ctx, &snomed.TranslateToRequest{ConceptId: 83942000, TargetId: 991411000000109})
		if err != nil {
			t.Fatal(err)
		}

		if t2.Translations[0].GetConcept().GetConceptId() != 45170000 {
			t.Fatalf("did not translate ADEM into encephalitis via emergency unit reference set. got: %v", t2)
		}
	})
	t.Run("FromCrossMap", func(t *testing.T) {
		response, err := c.FromCrossMap(ctx, &snomed.TranslateFromRequest{RefsetId: 999002271000000101, S: "G35X"})
		if err != nil {
			t.Fatal(err)
		}
		result := make(map[int64]struct{})
		for _, item := range response.Translations {
			result[item.ReferenceSetItem.ReferencedComponentId] = struct{}{}
		}
		if _, ok := result[24700007]; !ok {
			t.Fatalf("Did not correctly Reverse map G35X from ICD-10 into 'multiple sclerosis'. got:%v\n", result)
		}
	})
	t.Run("CrossMap", func(t *testing.T) {
		// test translating MS into ICD-10
		t2, err := c.CrossMap(ctx, &snomed.TranslateToRequest{ConceptId: 24700007, TargetId: 999002271000000101})
		if err != nil {
			t.Fatal(err)
		}
		result := make([]*snomed.ReferenceSetItem, 0)
		icd10codes := make(map[string]struct{})
		for {
			crossmap, err := t2.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			result = append(result, crossmap)
			if simpleMap := crossmap.GetSimpleMap(); simpleMap != nil {
				icd10codes[simpleMap.MapTarget] = struct{}{}
			}
			if complexMap := crossmap.GetComplexMap(); complexMap != nil {
				icd10codes[complexMap.MapTarget] = struct{}{}
			}
		}
		if _, ok := icd10codes["G35X"]; !ok {
			t.Fatalf("didn't match multiple sclerosis to ICD-10. expected: G35X, got: %v", icd10codes)
		}
	})
	t.Run("Subsumption", func(t *testing.T) {
		// test subsumption - could use a static table here...
		s1, err := c.Subsumes(ctx, &snomed.SubsumptionRequest{CodeA: 45170000, CodeB: 83942000})
		if err != nil {
			t.Fatal(err)
		}
		if s1.GetResult() != snomed.SubsumptionResponse_SUBSUMES {
			t.Fatalf("Encephalitis does not subsume ADEM, and it should. response:%v", s1.GetResult())
		}
		s2, err := c.Subsumes(ctx, &snomed.SubsumptionRequest{CodeA: 83942000, CodeB: 45170000})
		if err != nil {
			t.Fatal(err)
		}
		if s2.GetResult() != snomed.SubsumptionResponse_SUBSUMED_BY {
			t.Fatalf("Encephalitis does not subsume ADEM, and it should. response:%v", s2.GetResult())
		}
		s3, err := c.Subsumes(ctx, &snomed.SubsumptionRequest{CodeA: 83942000, CodeB: 24700007})
		if err != nil {
			t.Fatal(err)
		}
		if s3.GetResult() != snomed.SubsumptionResponse_NOT_SUBSUMED {
			t.Fatalf("Encephalitis subsumes multiple sclerosis, and it should not. response:%v", s3.GetResult())
		}

	})

	t.Run("Parse", func(t *testing.T) {
		e := "64572001 |disease|: 246454002 |occurrence| = 255407002 |neonatal|,  363698007 |finding site| = 113257007 |structure of cardiovascular system|"
		exp, err := c.Parse(context.Background(), &snomed.ParseRequest{S: e})
		if err != nil {
			t.Fatal(err)
		}
		if exp.GetClause().GetFocusConcepts()[0].ConceptId != 64572001 {
			t.Errorf("expression not parsed correctly, got %v\n", exp)
		}
	})
}
