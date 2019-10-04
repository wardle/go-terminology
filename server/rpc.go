package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/wardle/go-terminology/expression"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type coreServer struct {
	svc  *terminology.Svc
	lang []language.Tag // default language to use, if not explitly requested
}

// Options defines the options for a server.
type Options struct {
	RPCPort         int
	RESTPort        int
	DefaultLanguage string
}

// DefaultOptions provides some default options
var DefaultOptions = &Options{
	RPCPort:         8081,
	RESTPort:        8080,
	DefaultLanguage: "en-GB",
}

// RunServer runs a GRPC and a gateway REST server concurrently
func RunServer(svc *terminology.Svc, opts Options) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", opts.RPCPort))
	if err != nil {
		log.Printf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	tags, _, err := language.ParseAcceptLanguage(opts.DefaultLanguage)
	if err != nil {
		return err
	}
	go func() {
		impl := &coreServer{svc: svc, lang: tags}
		server := grpc.NewServer()
		health.RegisterHealthServer(server, impl)
		snomed.RegisterSnomedCTServer(server, impl)
		snomed.RegisterSearchServer(server, impl)
		log.Printf("gRPC Listening on %s\n", lis.Addr().String())
		server.Serve(lis)
	}()
	clientAddr := fmt.Sprintf("localhost:%d", opts.RPCPort)
	addr := fmt.Sprintf(":%d", opts.RESTPort)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()} // TODO:use better options
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(headerMatcher))
	if err := snomed.RegisterSnomedCTHandlerFromEndpoint(ctx, mux, clientAddr, dialOpts); err != nil {
		log.Fatalf("failed to create HTTP reverse proxy: %v", err)
	}
	if err := snomed.RegisterSearchHandlerFromEndpoint(ctx, mux, clientAddr, dialOpts); err != nil {
		log.Fatalf("failed to create reverse proxy for search service: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
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

// determine preferred language tags from the context, or fallback to a reasonable default
func (ss *coreServer) languageTags(ctx context.Context) ([]language.Tag, error) {
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
	return ss.lang, nil
}
func (ss *coreServer) GetConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.Concept, error) {
	c, err := ss.svc.Concept(conceptID.Identifier)
	if err == terminology.ErrNotFound {
		return nil, status.Errorf(codes.NotFound, "Concept not found with identifier %d", conceptID.Identifier)
	}
	return c, err
}

func (ss *coreServer) GetExtendedConcept(ctx context.Context, conceptID *snomed.SctID) (*snomed.ExtendedConcept, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	return ss.svc.ExtendedConcept(conceptID.Identifier, tags)
}

func (ss *coreServer) GetReferenceSets(conceptID *snomed.SctID, server snomed.SnomedCT_GetReferenceSetsServer) error {
	refsets, err := ss.svc.ComponentReferenceSets(conceptID.Identifier)
	if err != nil {
		return err
	}
	for _, refsetID := range refsets {
		items, err := ss.svc.ComponentFromReferenceSet(refsetID, conceptID.Identifier)
		if err != nil {
			return err
		}
		for _, item := range items {
			server.Send(item)
		}
	}
	return nil
}

func (ss *coreServer) GetReferenceSetItem(ctx context.Context, itemID *snomed.ReferenceSetItemID) (*snomed.ReferenceSetItem, error) {
	return ss.svc.ReferenceSetItem(itemID.Identifier)
}

func (ss *coreServer) GetDescriptions(ctx context.Context, conceptID *snomed.SctID) (*snomed.ConceptDescriptions, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	result := new(snomed.ConceptDescriptions)
	if result.Concept, err = ss.svc.Concept(conceptID.Identifier); err != nil {
		if err == terminology.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "Concept not found with identifier %d", conceptID.Identifier)
		}
		return nil, err
	}
	result.PreferredDescription, err = ss.svc.PreferredSynonym(conceptID.Identifier, tags)
	if err != nil {
		return nil, err
	}
	synonyms := make([]*snomed.Description, 0)
	definitions := make([]*snomed.Description, 0)
	descs, err := ss.svc.Descriptions(conceptID.Identifier)
	if err != nil {
		return nil, err
	}
	for _, d := range descs {
		if d.IsFullySpecifiedName() {
			result.FullySpecifiedName = d
		} else if d.IsSynonym() {
			synonyms = append(synonyms, d)
		} else if d.IsDefinition() {
			definitions = append(definitions, d)
		}

	}
	result.Synonyms = synonyms
	result.Definitions = definitions
	return result, nil
}

func (ss *coreServer) GetAllChildren(conceptID *snomed.SctID, stream snomed.SnomedCT_GetAllChildrenServer) error {
	tags, err := ss.languageTags(stream.Context())
	if err != nil {
		return err
	}
	children := ss.svc.StreamAllChildrenIDs(stream.Context(), conceptID.Identifier, 1000000)
	crch := ss.svc.StreamConceptReferences(stream.Context(), children, 4, tags)
	for cr := range crch {
		if cr.Err != nil {
			return cr.Err
		}
		stream.Send(cr.ConceptReference)
	}
	return nil
}

func (ss *coreServer) GetDescription(ctx context.Context, id *snomed.SctID) (*snomed.Description, error) {
	return ss.svc.Description(id.Identifier)
}

// CrossMap translates a SNOMED CT concept into an external code system, as defined by the map reference
// set specified in this request.
func (ss *coreServer) CrossMap(tr *snomed.CrossMapRequest, stream snomed.SnomedCT_CrossMapServer) error {
	targets, err := ss.svc.ComponentFromReferenceSet(tr.RefsetId, tr.ConceptId)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return status.Errorf(codes.NotFound, "Unable to map %d to reference set %d", tr.ConceptId, tr.RefsetId)
	}
	for _, target := range targets {
		stream.Send(target)
	}
	return nil
}

// Map translates a SNOMED CT concept into the best match in a destination simple reference set
func (ss *coreServer) Map(ctx context.Context, tr *snomed.MapRequest) (*snomed.MapResponse, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	members := make(map[int64]struct{})
	if tr.RefsetId != 0 {
		members, err = ss.svc.ReferenceSetComponents(tr.RefsetId) // get all reference set members
		if err != nil {
			return nil, err
		}
		if len(members) == 0 {
			return nil, status.Errorf(codes.NotFound, "Reference set %d not installed or has no members", tr.RefsetId)
		}
	}
	if len(tr.TargetId) > 0 {
		for _, c := range tr.TargetId {
			members[c] = struct{}{}
		}
	}
	if len(members) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to map: no valid target")
	}

	includeParents := false
	if tr.Parents == snomed.MapRequest_ALWAYS {
		includeParents = true
	}
DoMap:
	mapped, err := ss.svc.GenericiseTo(tr.ConceptId, includeParents, members)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to map %d: %s", tr.ConceptId, err)
	}
	count := len(mapped)
	if count == 0 {
		if includeParents == false && tr.Parents == snomed.MapRequest_FALLBACK {
			includeParents = true
			goto DoMap
		}
		return nil, status.Errorf(codes.NotFound, "Unable to map %d", tr.ConceptId)
	}
	translations := make([]*snomed.ConceptReference, count)
	for i, c := range mapped {
		cr, err := ss.svc.ConceptReference(c, tags)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error: %v", err)
		}
		translations[i] = cr
	}
	result := new(snomed.MapResponse)
	result.Translations = translations
	return result, nil
}

// FromCrossMap translates an external code into SNOMED CT, if possible.
func (ss *coreServer) FromCrossMap(ctx context.Context, r *snomed.TranslateFromRequest) (*snomed.TranslateFromResponse, error) {
	items, err := ss.svc.MapTarget(r.RefsetId, r.S)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, status.Errorf(codes.NotFound, "target '%s' not found in refset %d", r.S, r.RefsetId)
	}
	response := new(snomed.TranslateFromResponse)
	rr := make([]*snomed.TranslateFromResponse_Item, 0)
	for _, item := range items {
		if item.Active == false && r.IncludeInactive == false {
			continue
		}
		c, err := ss.svc.Concept(item.ReferencedComponentId)
		if err != nil {
			return nil, err
		}
		if c.Active == false && r.IncludeInactive == false {
			continue
		}
		rItem := new(snomed.TranslateFromResponse_Item)
		rr = append(rr, rItem)
		rItem.ReferenceSetItem = item
		rItem.Concept = c
		if c.Active == false { // for inactive concepts, help the client by providing associations.
			var err error
			rItem.SameAs, err = ss.svc.GetAssociations(c.Id, snomed.SameAsReferenceSet)
			if err != nil {
				return nil, err
			}
			rItem.PossiblyEquivalentTo, err = ss.svc.GetAssociations(c.Id, snomed.PossiblyEquivalentToReferenceSet)
			if err != nil {
				return nil, err
			}
			rItem.SimilarTo, err = ss.svc.GetAssociations(c.Id, snomed.SimilarToReferenceSet)
			if err != nil {
				return nil, err
			}
			rItem.ReplacedBy, err = ss.svc.GetAssociations(c.Id, snomed.ReplacedByReferenceSet)
			if err != nil {
				return nil, err
			}
		}
	}
	response.Translations = rr
	return response, nil
}

// Subsumes determines whether code A subsumes code B, according to the definition
// in the HL7 FHIR terminology service specification.
// See https://www.hl7.org/fhir/terminology-service.html
func (ss *coreServer) Subsumes(ctx context.Context, r *snomed.SubsumptionRequest) (*snomed.SubsumptionResponse, error) {
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

func (ss *coreServer) Parse(ctx context.Context, r *snomed.ParseRequest) (*snomed.Expression, error) {
	return expression.Parse(r.S)
}

// Refinements determines the appropriate refinements for an arbitrary concept
// It is quite easy to do, we find the relationships and additionally determine
// whether the concept's attributes exist in the lateralisable reference set.
// TODO: this would be better deprecated in favour of using only expressions
// that would mean normalising any concept into an expression and *then* deriving
// possible refinements for that expression, instead.
func (ss *coreServer) Refinements(ctx context.Context, r *snomed.RefinementRequest) (*snomed.RefinementResponse, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	c, err := ss.svc.Concept(r.ConceptId)
	if err != nil {
		if err == terminology.ErrNotFound {
			return nil, status.Errorf(codes.NotFound, "Concept %d not found", r.ConceptId)
		}
		return nil, err
	}
	rels, err := ss.svc.ParentRelationships(c.Id)
	if err != nil {
		return nil, err
	}
	attrs := make([]*snomed.RefinementResponse_Refinement, 0)
	properties := make(map[int64]struct{})
	for _, rel := range rels {
		if rel.Active && rel.TypeId != snomed.IsA {
			if _, done := properties[rel.DestinationId]; done {
				continue
			}
			properties[rel.DestinationId] = struct{}{}
			attr, err := makeRefinement(ctx, ss.svc, rel.TypeId, rel.DestinationId, tags)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)

			if rel.TypeId == snomed.BodyStructure || rel.TypeId == snomed.ProcedureSiteDirect || rel.TypeId == snomed.FindingSite {
				if _, done := properties[snomed.Side]; !done {
					islat, err := isLateralisable(ss.svc, rel.DestinationId)
					if err != nil {
						return nil, err
					}
					if islat {
						lat, err := makeRefinement(ctx, ss.svc, snomed.Laterality, snomed.Side, tags)
						if err != nil {
							return nil, err
						}
						attrs = append(attrs, lat)
					}
				}
			}

		}
	}
	response := new(snomed.RefinementResponse)
	response.Concept = c
	response.Refinements = attrs
	return response, nil
}

func makeRefinement(ctx context.Context, svc *terminology.Svc, attributeID int64, rootID int64, tags []language.Tag) (*snomed.RefinementResponse_Refinement, error) {
	cc, err := svc.Concepts(attributeID, rootID)
	if err != nil {
		return nil, err
	}
	attr := new(snomed.RefinementResponse_Refinement)
	attr.Attribute, err = svc.ConceptReference(attributeID, tags)
	if err != nil {
		return nil, err
	}
	attr.RootValue, err = svc.ConceptReference(rootID, tags)
	if err != nil {
		return nil, err
	}
	attr.Choices = make([]*snomed.ConceptReference, 0)
	valueSet, err := svc.AllChildren(ctx, cc[1], 500)
	if err == nil {
		for _, v := range valueSet {
			if v.Active {
				cr, err := svc.ConceptReference(v.Id, tags)
				if err != nil {
					return nil, err
				}
				attr.Choices = append(attr.Choices, cr)
			}
		}
	}
	return attr, nil
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

func (ss *coreServer) Extract(ctx context.Context, r *snomed.ExtractRequest) (*snomed.ExtractResponse, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	return ss.svc.Extract(r, tags)
}

var _ snomed.SnomedCTServer = (*coreServer)(nil)

// Check is a health check, implementing the grpc-health service
// see https://godoc.org/google.golang.org/grpc/health/grpc_health_v1#HealthServer
func (ss *coreServer) Check(ctx context.Context, r *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	response := new(health.HealthCheckResponse)
	response.Status = health.HealthCheckResponse_SERVING
	return response, nil
}

func (ss *coreServer) Watch(r *health.HealthCheckRequest, w health.Health_WatchServer) error {
	return nil
}

func (ss *coreServer) Search(ctx context.Context, sr *snomed.SearchRequest) (*snomed.SearchResponse, error) {
	tags, err := ss.languageTags(ctx)
	if err != nil {
		return nil, err
	}
	response, err := ss.svc.Search(sr, tags)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ss *coreServer) Synonyms(sr *snomed.SynonymRequest, response snomed.Search_SynonymsServer) error {
	tags, err := ss.languageTags(response.Context())
	if err != nil {
		return err
	}
	search := snomed.SearchRequest{ // turn synonym request into a search request
		IsA:             sr.IsA,
		Fuzzy:           sr.Fuzzy,
		IncludeInactive: sr.IncludeInactive,
		MaximumHits:     sr.MaximumHits,
		S:               sr.S,
	}
	results, err := ss.svc.Search(&search, tags)
	if err != nil {
		return err
	}
	concepts := make(map[int64]struct{})
	maxChildren := 200
	if sr.MaximumHits > 0 {
		maxChildren = int(sr.MaximumHits)
	}
	for _, result := range results.Items {
		concepts[result.ConceptId] = struct{}{}
		if sr.IncludeChildren {
			children, err := ss.svc.AllChildrenIDs(response.Context(), result.ConceptId, maxChildren)
			if err != nil {
				return err
			}
			for _, child := range children {
				concepts[child] = struct{}{}
			}
		}
	}
	for conceptID := range concepts {
		descriptions, err := ss.svc.Descriptions(conceptID)
		if err != nil {
			return err
		}
		for _, d := range descriptions { // TODO: should this limit descriptions by language?
			if d.IsFullySpecifiedName() || (!d.Active && !sr.IncludeInactive) {
				continue
			}
			item := snomed.SynonymResponseItem{
				S: d.Term,
			}
			response.Send(&item)
		}
	}
	return nil
}
