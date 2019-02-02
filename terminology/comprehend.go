package terminology

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehendmedical"
	"github.com/wardle/go-terminology/snomed"
	"golang.org/x/text/language"
	"strings"
)

// typeRootMap gives the root SNOMED identifier for the type of entity found by the comprehend API
var typeRootMap = map[string]int64{
	"MEDICATION-GENERIC_NAME":            373873005,
	"MEDICAL_CONDITION-DX_NAME":          404684003, // clinical finding
	"TEST_TREATMENT_PROCEDURE-TEST_NAME": 103693007, // diagnostic procedures
	"ANATOMY-SYSTEM_ORGAN_SITE":          123037004, // body structure
}

// Extract using NLP to extract entities from a block of text using Amazon Comprehend's API
func (svc *Svc) Extract(r *snomed.ExtractRequest, tags []language.Tag) (*snomed.ExtractResponse, error) {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"), // TODO: allow configuration
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		return nil, err
	}
	client := comprehendmedical.New(s)
	input := comprehendmedical.DetectEntitiesInput{}
	input.SetText(r.S)
	result, err := client.DetectEntities(&input)
	if err != nil {
		return nil, err
	}
	entities := make([]*snomed.ExtractResponse_Entity, 0)
	for _, entity := range result.Entities {
		key := *entity.Category + "-" + *entity.Type
		fmt.Printf("entity: %s, category: %s, type: %s\n", *entity.Text, *entity.Category, *entity.Type)
		roots := make([]int64, 0)
		root := typeRootMap[key]
		if root != 0 {
			roots = append(roots, root)
		}
		responseEntity := new(snomed.ExtractResponse_Entity)
		responseEntity.Text = *entity.Text
		responseEntity.Score = *entity.Score
		for _, trait := range entity.Traits {
			if *trait.Name == "NEGATION" {
				responseEntity.Negated = true
			}
		}
		sr, err := svc.Search(&snomed.SearchRequest{
			S:           *entity.Text,
			IsA:         roots,
			MaximumHits: 5,
		}, tags)
		if err != nil {
			return nil, err
		}
		concepts := make([]*snomed.ConceptReference, 0)
		// can we find any exact matches for the entity text - only use those if so
		for _, item := range sr.Items {
			if strings.EqualFold(item.Term, *entity.Text) {
				concepts = append(concepts, &snomed.ConceptReference{
					ConceptId: item.ConceptId,
					Term:      item.Term,
				})
			}
		}
		// if we have no exact matches, find all possible matches (client will have to show to user)
		if len(concepts) == 0 {
			for _, item := range sr.Items {
				concepts = append(concepts, &snomed.ConceptReference{
					ConceptId: item.ConceptId,
					Term:      item.Term,
				})
			}
		}
		if len(concepts) > 0 {
			responseEntity.Concepts = concepts
			responseEntity.BestMatch = concepts[0].ConceptId
		}
		if len(concepts) > 1 {
			// TODO: need to calculate the most generic of the list of concepts (very useful for analytics)
		}
		entities = append(entities, responseEntity)
	}
	response := &snomed.ExtractResponse{
		Entities: entities,
	}
	return response, nil
}
