package server

import (
	"fmt"
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
	c := NewSnomedCTClient(conn)
	// Test GetConcept as a subtest
	t.Run("GetConcept", func(t *testing.T) {
		c, err := c.GetConcept(context.Background(), &SctID{Identifier: 24700007})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Concept: %d", c.GetId())
		if c.GetId() != 24700007 {
			t.Error("Expected '24700007', got ", c.GetId())
		}

	})
}
