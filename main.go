package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	logging "github.com/bkpeh/httpsvr/util"
	hsvr "github.com/bkpeh/httpsvr/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	hsvr.SetLog("logs/http.log")
	http.HandleFunc("/", hsvr.Index)

	//defer logging.GetLogFile().Close()

	//cfg, err := external.LoadDefaultAWSConfig()

	//config, err := config.LoadDefaultConfig(context.TODO())

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-west-2" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:8000",
				SigningRegion: "us-west-2",
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"), config.WithEndpointResolver(customResolver))

	svc := dynamodb.NewFromConfig(cfg)

	output, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 new(string),
		AttributesToGet:           []string{},
		ConditionalOperator:       "",
		ConsistentRead:            new(bool),
		ExclusiveStartKey:         map[string]types.AttributeValue{},
		ExpressionAttributeNames:  map[string]string{},
		ExpressionAttributeValues: map[string]types.AttributeValue{":i": &types.AttributeValueMemberS{"1020"}},
		FilterExpression:          new(string),
		IndexName:                 new(string),
		KeyConditionExpression:    aws.String("EID = :i"),
		KeyConditions:             map[string]types.Condition{},
		Limit:                     new(int32),
		ProjectionExpression:      new(string),
		QueryFilter:               map[string]types.Condition{},
		ReturnConsumedCapacity:    "",
		ScanIndexForward:          new(bool),
		Select:                    "",
	})
	// Build the request with its input parameters
	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5)})

	fmt.Println("output:", output)
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}
	/*
	   	// Set the AWS Region that the service clients should use
	   	//cfg.Region = endpoints.UsWest2RegionID
	   	sess := session.Must(session.NewSessionWithOptions(session.Options{
	   		Config: aws.Config{
	   			Region:           "us-west-2",
	   			EndpointResolver: endpoints.ResolverFunc(customResolver),
	   		},
	   	}))
	   	// Using the Config value, create the DynamoDB client
	   	svc := dynamodb.New(sess)

	   	// Build the request with its input parameters
	   	req := svc.DescribeTableRequest(&dynamodb.DescribeTableInput{
	       TableName: aws.String("myTable"),
	   })
	*/
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		logging.LogError("logs/http.log", "Main:"+err.Error())
	}
}
