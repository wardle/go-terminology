package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type myServer struct {
	svc *terminology.Svc
}

// RunServer runs a GRPC and a gateway REST server concurrently
func RunServer(svc *terminology.Svc, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	go func() {
		server := grpc.NewServer()
		snomed.RegisterSnomedCTServer(server, &myServer{svc: svc})
		log.Printf("gRPC Listening on %s\n", lis.Addr().String())
		server.Serve(lis)
	}()
	clientAddr := fmt.Sprintf("localhost:%d", port)
	addr := fmt.Sprintf(":%d", port+1)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()
	if err := snomed.RegisterSnomedCTHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
	return
}

func (ss *myServer) GetConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.Concept, error) {
	return ss.svc.GetConcept(conceptID.Identifier)
}

func (ss *myServer) GetExtendedConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.ExtendedConcept, error) {
	tags, _, _ := language.ParseAcceptLanguage("en-GB") // TODO(mw): better language support
	return ss.svc.GetExtendedConcept(conceptID.Identifier, tags)
}

func (ss *myServer) GetDescriptions(conceptID *snomed.SctID, server snomed.SnomedCT_GetDescriptionsServer) error {
	c, err := ss.svc.GetConcept(conceptID.Identifier)
	if err != nil {
		return err
	}
	descs, err := ss.svc.GetDescriptions(c)
	if err != nil {
		return err
	}
	for _, d := range descs {
		server.Send(d)
	}
	return nil
}

func (ss *myServer) GetDescription(ctx context.Context, id *snomed.SctID) (*snomed.Description, error) {
	return ss.svc.GetDescription(id.Identifier)
}

// Translate translates a SNOMED CT concept into the best match in a destination simple reference set
func (ss *myServer) Translate(ctx context.Context, tr *snomed.TranslateRequest) (*snomed.TranslateResponse, error) {
	response := new(snomed.TranslateResponse)
	target, err := ss.svc.GetFromReferenceSet(tr.TargetId, tr.ConceptId)
	if err != nil {
		return nil, err
	}
	if target != nil { // we have found our concept in the reference set, so return that entry
		simple := target.GetSimple() // if we have a simple refset, then return it
		if simple != nil {
			c, err := ss.svc.GetConcept(tr.ConceptId)
			if err != nil {
				return nil, err
			}
			rc := new(snomed.TranslateResponse_Concept)
			rc.Concept = c
			response.Result = rc
			return response, nil
		}
		// we must have a map, so let's return
		mr := new(snomed.TranslateResponse_MappedResponse)
		mr.Items = make([]*snomed.ReferenceSetItem, 1)
		mr.Items[0] = target
		mapped := new(snomed.TranslateResponse_Mapped)
		mapped.Mapped = mr
		response.Result = mapped
		return response, nil
	}
	// we have not found the source concept, so try to genericise
	c, err := ss.svc.GetConcept(tr.ConceptId)
	if err != nil {
		return nil, err
	}
	members, err := ss.svc.GetReferenceSetItems(tr.TargetId) // get all reference set members
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("Reference set %d not found or has no members", tr.TargetId)
	}
	generic, found := ss.svc.GenericiseTo(c, members)
	if found {
		rc := new(snomed.TranslateResponse_Concept)
		rc.Concept = generic
		response.Result = rc
		return response, nil
	}
	return nil, fmt.Errorf("Unable to translate %d to %d", tr.ConceptId, tr.TargetId)
}

// Subsumes determines whether code A subsumes code B, according to the definition
// in the HL7 FHIR terminology service specification.
// See https://www.hl7.org/fhir/terminology-service.html
func (ss *myServer) Subsumes(ctx context.Context, r *snomed.SubsumptionRequest) (*snomed.SubsumptionResponse, error) {
	res := snomed.SubsumptionResponse{}
	if r.CodeA == r.CodeB {
		res.Result = snomed.SubsumptionResponse_EQUIVALENT
		return &res, nil
	}
	c, err := ss.svc.GetConcept(r.CodeB)
	if err != nil {
		return nil, err
	}
	if ss.svc.IsA(c, r.CodeA) {
		res.Result = snomed.SubsumptionResponse_SUBSUMES
		return &res, nil
	}
	c, err = ss.svc.GetConcept(r.CodeA)
	if err != nil {
		return nil, err
	}
	if ss.svc.IsA(c, r.CodeB) {
		res.Result = snomed.SubsumptionResponse_SUBSUMED_BY
		return &res, nil
	}
	res.Result = snomed.SubsumptionResponse_NOT_SUBSUMED
	return &res, nil
}

var _ snomed.SnomedCTServer = (*myServer)(nil)
