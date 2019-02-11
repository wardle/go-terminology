package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

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
	RPCPort:         8080,
	RESTPort:        8081,
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
	result.PreferredDescription, _, err = ss.svc.PreferredSynonym(conceptID.Identifier, tags)
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

func (ss *coreServer) GetDescription(ctx context.Context, id *snomed.SctID) (*snomed.Description, error) {
	return ss.svc.Description(id.Identifier)
}

// CrossMap translates a SNOMED CT concept into an external code system, as defined by the map reference
// set specified in this request.
func (ss *coreServer) CrossMap(tr *snomed.TranslateToRequest, stream snomed.SnomedCT_CrossMapServer) error {
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
func (ss *coreServer) Map(ctx context.Context, tr *snomed.TranslateToRequest) (*snomed.Concept, error) {
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
func (ss *coreServer) FromCrossMap(ctx context.Context, r *snomed.TranslateFromRequest) (*snomed.TranslateFromResponse, error) {
	items, err := ss.svc.MapTarget(r.RefsetId, r.S)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, status.Errorf(codes.NotFound, "target '%s' not found in refset %d", r.S, r.RefsetId)
	}

	response := new(snomed.TranslateFromResponse)

	bestMatch, err := findBestMatch(ss.svc, r.RefsetId, items)
	if err != nil {
		return nil, err
	}
	if bestMatch != 0 {
		response.BestMatch = bestMatch
	}

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
		if c.Active == false {
			var err error
			rr[i].SameAs, err = getAssociations(ss.svc, snomed.SameAsReferenceSet, c.Id)
			if err != nil {
				return nil, err
			}
			rr[i].PossiblyEquivalentTo, err = getAssociations(ss.svc, snomed.PossiblyEquivalentToReferenceSet, c.Id)
			if err != nil {
				return nil, err
			}
			rr[i].SimilarTo, err = getAssociations(ss.svc, snomed.SimilarToReferenceSet, c.Id)
			if err != nil {
				return nil, err
			}
			rr[i].ReplacedBy, err = getAssociations(ss.svc, snomed.ReplacedByReferenceSet, c.Id)
			if err != nil {
				return nil, err
			}

		}
	}
	return response, nil
}

func findBestMatch(svc *terminology.Svc, refsetID int64, items []*snomed.ReferenceSetItem) (int64, error) {
	var candidateConcepts []*snomed.Concept
	for _, item := range items {
		c, err := svc.Concept(item.ReferencedComponentId)
		if err != nil {
			return 0, err
		}

		// Exclude inactive concepts
		if c.Active == false {
			continue
		}

		// Only for complex mappings
		complexMap := item.GetComplexMap()
		if complexMap != nil {
			// Get all mappings for concept to determine if mapTarget mapPriority is highest and so a most appropriate mapping
			mappings, err := svc.ComponentFromReferenceSet(refsetID, c.Id)
			if err != nil {
				return 0, err
			}

			skipItem := false
			for _, mapping := range mappings {
				mappingComplexMap := mapping.GetComplexMap()
				if mappingComplexMap.MapPriority > complexMap.MapPriority && mappingComplexMap.MapTarget != complexMap.MapTarget {
					skipItem = true
				}
			}
			if skipItem {
				continue
			}
		}

		// Build map of candidate concepts
		candidateConcepts = append(candidateConcepts, c)
	}

	var subsumptionMap = make(map[int64][]int64)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	// Determine which candidateConcepts are subsumed by each other
	// (This is quite expensive as is uses AllParents which is recursive and unoptimised)
	for _, concept := range candidateConcepts {
		wg.Add(1)
		go func(concept *snomed.Concept) {
			defer wg.Done()
			parents, err := svc.AllParents(concept)
			if err != nil {
				// TODO: Better error handling e.g return via channel
				return
			}
			parentIdMap := make(map[int64]bool)
			for _, p := range parents {
				parentIdMap[p.Id] = true
			}
			for _, test := range candidateConcepts {
				// Skip current concept
				if test.Id == concept.Id {
					continue
				}
				if _, ok := parentIdMap[test.Id]; ok {
					mutex.Lock()
					subsumptionMap[test.Id] = append(subsumptionMap[test.Id], concept.Id)
					mutex.Unlock()
				}
			}
		}(concept)
	}
	wg.Wait()

	// fmt.Printf("%#v", subsumptionMap)

	var bestmatch int64
	var subsumptionCount int
	// Find candidateConcept which subsumes the most other concepts
	// TODO: better deal with occasions where no clear winner
	for conceptId, subsumes := range subsumptionMap {
		if len(subsumes) > subsumptionCount {
			bestmatch = conceptId
			subsumptionCount = len(subsumes)
		}
	}

	// If no candidateConcept subsumes any other pick concept with shortest path to root
	if len(subsumptionMap) == 0 {
		for _, concept := range candidateConcepts {
			wg.Add(1)
			go func(concept *snomed.Concept) {
				defer wg.Done()
				path, err := svc.ShortestPathToRoot(concept)
				if err != nil {
					// TODO: Better error handling e.g return via channel
					return
				}
				mutex.Lock()
				if subsumptionCount > len(path) || subsumptionCount == 0 {
					bestmatch = concept.Id
					subsumptionCount = len(path)
				}
				mutex.Unlock()
			}(concept)
		}
		wg.Wait()
	}

	return bestmatch, nil
}

func getAssociations(svc *terminology.Svc, refsetID int64, conceptID int64) ([]int64, error) {
	items, err := svc.ComponentFromReferenceSet(refsetID, conceptID)
	if err != nil {
		return nil, err
	}
	result := make([]int64, len(items))
	for i, item := range items {
		result[i] = item.GetAssociation().GetTargetComponentId()
	}
	return result, nil
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
	return expression.ParseExpression(r.S)
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
			cc, err := ss.svc.Concepts(rel.TypeId, rel.DestinationId)
			if err != nil {
				return nil, err
			}
			attr := new(snomed.RefinementResponse_Refinement)
			attr.Attribute = makeConceptReference(ss.svc, cc[0], tags)
			attr.RootValue = makeConceptReference(ss.svc, cc[1], tags)
			attr.Choices = make([]*snomed.ConceptReference, 0)
			valueSet, err := ss.svc.AllChildren(cc[1], 1000)
			if err == nil {
				for _, v := range valueSet {
					if v.Active {
						attr.Choices = append(attr.Choices, makeConceptReference(ss.svc, v, tags))
					}
				}
			}
			attrs = append(attrs, attr)
			if rel.TypeId == snomed.BodyStructure || rel.TypeId == snomed.ProcedureSiteDirect || rel.TypeId == snomed.FindingSite {
				if _, done := properties[snomed.Side]; !done {
					islat, err := isLateralisable(ss.svc, rel.DestinationId)
					if err != nil {
						return nil, err
					}
					if islat {
						lat := new(snomed.RefinementResponse_Refinement)
						ll, err := ss.svc.Concepts(snomed.Laterality, snomed.Side)
						if err != nil {
							return nil, err
						}
						lat.Attribute = makeConceptReference(ss.svc, ll[0], tags)
						lat.RootValue = makeConceptReference(ss.svc, ll[1], tags)
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
	d := svc.MustGetPreferredSynonym(c.Id, tags)
	r.Term = d.Term
	return r
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
	return ss.svc.Search(sr, tags)
}
