package server

import (
	"fmt"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
)

const (
	dbFilename = "../snomed.db" // real, live database
	port       = 8080
)

func TestMain(m *testing.M) {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		log.Printf("Skipping live tests in the absence of a live database %s", dbFilename)
		os.Exit(0)
	}
	svc, err := terminology.NewService(dbFilename, false)
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()
	go RunRPCServer(svc, port)
	os.Exit(m.Run())
}

func TestRpcClient(t *testing.T) {

	// Set up a connection to the Server.
	address := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := snomed.NewSnomedCTClient(conn)
	// Test GetConcept as a subtest
	t.Run("GetConcept", func(t *testing.T) {
		c, err := c.GetConcept(context.Background(), &snomed.SctID{Identifier: 24700007})
		if err != nil {
			t.Fatal(err)
		}
		if c.GetId() != 24700007 {
			t.Error("Expected '24700007', got ", c.GetId())
		}

	})
	t.Run("Translate", func(t *testing.T) {
		// test translating MS into emergency care reference set - should give MS
		t1, err := c.Translate(context.Background(), &snomed.TranslateRequest{ConceptId: 24700007, TargetId: 991411000000109})
		if err != nil {
			t.Fatal(err)
		}
		if t1.GetConcept().Id != 24700007 {
			t.Fatalf("failed to find multiple sclerosis in the emergency care reference set. found: %v", t1)
		}
		// test translating MS into ICD-10
		t2, err := c.Translate(context.Background(), &snomed.TranslateRequest{ConceptId: 24700007, TargetId: 999002261000000108})
		if err != nil {
			t.Fatal(err)
		}
		if t2.GetReferenceSet().GetComplexMap().GetMapTarget() != "G35X" {
			t.Fatalf("didn't match multiple sclerosis to ICD-10. expected: G35X, got: %v", t2)
		}
		// test translating ADEM into emergency care reference set - should get encephalitis (45170000)
		t3, err := c.Translate(context.Background(), &snomed.TranslateRequest{ConceptId: 83942000, TargetId: 991411000000109})
		if err != nil {
			t.Fatal(err)
		}
		if t3.GetConcept().Id != 45170000 {
			t.Fatalf("did not translate ADEM into encephalitis via emergency unit reference set. got: %v", t3)
		}
		// test subsumption - could use a static table here...
		s1, err := c.Subsumes(context.Background(), &snomed.SubsumptionRequest{CodeA: 45170000, CodeB: 83942000})
		if err != nil {
			t.Fatal(err)
		}
		if s1.GetResult() != snomed.SubsumptionResponse_SUBSUMES {
			t.Fatalf("Encephalitis does not subsume ADEM, and it should. response:%v", s1.GetResult())
		}
		s2, err := c.Subsumes(context.Background(), &snomed.SubsumptionRequest{CodeA: 83942000, CodeB: 45170000})
		if err != nil {
			t.Fatal(err)
		}
		if s2.GetResult() != snomed.SubsumptionResponse_SUBSUMED_BY {
			t.Fatalf("Encephalitis does not subsume ADEM, and it should. response:%v", s2.GetResult())
		}
		s3, err := c.Subsumes(context.Background(), &snomed.SubsumptionRequest{CodeA: 83942000, CodeB: 24700007})
		if err != nil {
			t.Fatal(err)
		}
		if s3.GetResult() != snomed.SubsumptionResponse_NOT_SUBSUMED {
			t.Fatalf("Encephalitis subsumes multiple sclerosis, and it should not. response:%v", s3.GetResult())
		}

	})
}
