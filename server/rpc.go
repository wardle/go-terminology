package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/wardle/go-terminology/expression"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"strings"
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
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
	if err := snomed.RegisterSnomedCTHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
	return
}

// ensures GRPC gateway passes through the standard HTTP header Accept-Language as "accept-language"
// rather than munging the name prefixed with grpcgateway.
// delegates to default implementation for other headers.
func headerMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Accept-Language" {
		return "accept-language", true
	}
	return runtime.DefaultHeaderMatcher(headerName)
}

func (ss *myServer) languageTags(ctx context.Context) ([]language.Tag, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		preferred := md.Get("accept-language")
		if len(preferred) > 0 {
			pref := strings.Join(preferred, ";")
			tags, _, err := language.ParseAcceptLanguage(pref)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid accept-language header: %s", err)
			}
			return tags, nil
		}
	}
	tags, _, _ := language.ParseAcceptLanguage("en-GB")
	return tags, nil
}
func (ss *myServer) GetConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.Concept, error) {
	c, err := ss.svc.Concept(conceptID.Identifier)
	if err == terminology.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "Concept not found with identifier %d", conceptID.Identifier)
	}
	return c, err
}

func (ss *myServer) GetExtendedConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.ExtendedConcept, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	return ss.svc.ExtendedConcept(conceptID.Identifier, tags)
}

func (ss *myServer) GetDescriptions(conceptID *snomed.SctID, stream snomed.SnomedCT_GetDescriptionsServer) error {
	descs, err := ss.svc.Descriptions(conceptID.Identifier)
	if err != nil {
		return err
	}
	for _, d := range descs {
		stream.Send(d)
	}
	return nil
}

func (ss *myServer) GetDescription(ctx context.Context, id *snomed.SctID) (*snomed.Description, error) {
	return ss.svc.Description(id.Identifier)
}

// CrossMap translates a SNOMED CT concept into an external code system, as defined by the map reference
// set specified in this request.
func (ss *myServer) CrossMap(tr *snomed.TranslateToRequest, stream snomed.SnomedCT_CrossMapServer) error {
	targets, err := ss.svc.ComponentFromReferenceSet(tr.TargetId, tr.ConceptId)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return status.Errorf(codes.NotFound, "Unable to map %d to reference set %d", tr.ConceptId, tr.TargetId)
	}
	for _, target := range targets {
		stream.Send(target)
	}
	return nil
}

// Map translates a SNOMED CT concept into the best match in a destination simple reference set
func (ss *myServer) Map(ctx context.Context, tr *snomed.TranslateToRequest) (*snomed.Concept, error) {
	c, err := ss.svc.Concept(tr.ConceptId)
	if err != nil {
		return nil, err
	}
	members, err := ss.svc.ReferenceSetComponents(tr.TargetId) // get all reference set members
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, status.Errorf(codes.NotFound, "Reference set %d not installed or has no members", tr.TargetId)
	}
	generic, found := ss.svc.GenericiseTo(c, members)
	if found {
		return generic, nil
	}
	return nil, status.Errorf(codes.NotFound, "Unable to translate %d to %d", tr.ConceptId, tr.TargetId)
}

// FromCrossMap translates an external code into SNOMED CT, if possible.
func (ss *myServer) FromCrossMap(ctx context.Context, r *snomed.TranslateFromRequest) (*snomed.TranslateFromResponse, error) {
	items, err := ss.svc.MapTarget(r.RefsetId, r.S)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, status.Errorf(codes.NotFound, "target '%s' not found in refset %d", r.S, r.RefsetId)
	}
	response := new(snomed.TranslateFromResponse)
	rr := make([]*snomed.TranslateFromResponse_Item, len(items))
	response.Translations = rr
	for i, item := range items {
		rr[i] = new(snomed.TranslateFromResponse_Item)
		rr[i].ReferenceSetItem = item
		c, err := ss.svc.Concept(item.ReferencedComponentId)
		if err != nil {
			return nil, err
		}
		rr[i].Concept = c
	}
	return response, nil
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
	c, err := ss.svc.Concept(r.CodeB)
	if err != nil {
		return nil, err
	}
	if ss.svc.IsA(c, r.CodeA) {
		res.Result = snomed.SubsumptionResponse_SUBSUMES
		return &res, nil
	}
	c, err = ss.svc.Concept(r.CodeA)
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

func (ss *myServer) Parse(ctx context.Context, r *snomed.ParseRequest) (*snomed.Expression, error) {
	return expression.ParseExpression(r.S)
}

// Refinements determines the appropriate refinements for an arbitrary concept
// It is quite easy to do, we find the relationships and additionally determine
// whether the concept's attributes exist in the lateralisable reference set.
// TODO: this would be better deprecated in favour of using only expressions
// that would mean normalising any concept into an expression and *then* deriving
// possible refinements for that expression, instead.
func (ss *myServer) Refinements(ctx context.Context, r *snomed.RefinementRequest) (*snomed.RefinementResponse, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	c, err := ss.svc.Concept(r.ConceptId)
	if err != nil {
			return nil, status.Errorf(codes.NotFound, "Concept %d not found", r.ConceptId)
		}
		return nil, err
	}
	rels, err := ss.svc.ParentRelationships(c.Id)
	if err != nil {
		return nil, err
	}
	attrs := make([]*snomed.RefinementResponse_Refinement, 0)
	for _, rel := range rels {
		if rel.Active {
			cc, err := ss.svc.Concepts(rel.TypeId, rel.DestinationId)
			if err != nil {
				return nil, err
			}
			attr := new(snomed.RefinementResponse_Refinement)
			attr.Attribute = makeConceptReference(ss.svc, cc[0], tags)
			attr.Parent = makeConceptReference(ss.svc, cc[1], tags)
			attrs = append(attrs, attr)
		}
	}
	response := new(snomed.RefinementResponse)
	response.Concept = c
	response.Refinements = attrs
	return response, nil
}

// isLateralisable finds out whether the specific concept is lateralisable
func isLateralisable(svc *terminology.Svc, id int64) (bool, error) {
	rsis, err := svc.ComponentFromReferenceSet(snomed.LateralisableReferenceSet, id)
	if err != nil {
		return false, err
	}
	for _, rsi := range rsis {
		if rsi.Active {
			return true, nil
		}
	}
	return false, nil
}

func makeConceptReference(svc *terminology.Svc, c *snomed.Concept, tags []language.Tag) *snomed.ConceptReference {
	r := new(snomed.ConceptReference)
	r.ConceptId = c.Id
	d := svc.MustGetPreferredSynonym(c, tags)
	r.Term = d.Term
	return r
}

var _ snomed.SnomedCTServer = (*myServer)(nil)
