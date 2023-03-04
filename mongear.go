package mongear

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

// NewFilter returns a Mongo filter document (bson.D) from a given
// correctly formatted filter string.
func NewFilter(str string) (bson.D, error) {
	str = strings.TrimSpace(str)
	if strings.Index(str, "{") == 0 {
		var filter bson.D
		err := bson.UnmarshalExtJSON([]byte(str), false, &filter)
		if err != nil {
			return nil, err
		}
		return filter, nil
	} else {
		return nil, errors.New("Not a valid filter string")
	}
}

// NewStage returns a Mongo aggregation pipeline stage document(bsond.D) from a
// correctly formatted aggregation stage string.  These can be appended to an
// existing mongo.Pipeline and the whole pipeline passed to the collection's
// aggregate function.
//
// If you pass in the pipeline, the new stage will be appended to the pipeline.
// If you pass in nil for the pipeline parameter, the pipeline will be unaffected
//
// In both cases, the pipeline and the newly created stage are returned.
//
// This is function is useful for building aggreation pipelines a stage at a time
func NewStage(p mongo.Pipeline, str string) (mongo.Pipeline, bson.D, error) {
	stage, err := NewFilter(str)
	if err == nil && p != nil {
		p = append(p, stage)
	}
	return p, stage, err
}

// NewAggregate returns a Mongo aggregation pipeline ([]bsond.D) from a
// correctly formatted aggregation string.
func NewAggregate(str string) (mongo.Pipeline, error) {
	var aggregation = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") == 0 {
		var aggregate bson.D
		err := bson.UnmarshalExtJSON([]byte(str), false, &aggregate)
		if err != nil {
			return nil, err
		}
		aggregation = append(aggregation, aggregate)
	} else {
		return nil, errors.New("Not a valid aggregation string")
	}
	return aggregation, nil
}
