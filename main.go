package main

import (
	"net/http"

	logging "github.com/bkpeh/httpsvr/util"
	hsvr "github.com/bkpeh/httpsvr/web"
)

func main() {
	hsvr.SetLog("logs/http.log")
	http.HandleFunc("/", hsvr.Index)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		logging.LogError("logs/http.log", "Main:"+err.Error())
	}
	/*
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

		a := types.AttributeValueMemberS{Value: `"S":"easter@gmail.com"`}
		//a := types.AttributeValueMemberS{S: "easter@gmail.com"}
		output, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String("People"),
			Key:       map[string]types.AttributeValue{"Email": &a},
		})

		fmt.Println("output:", output)

		// Build the request with its input parameters
		resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit: aws.Int32(5)})

		if err != nil {
			log.Fatalf("failed to list tables, %v", err)
		}

		fmt.Println("Tables:")
		for _, tableName := range resp.TableNames {
			fmt.Println(tableName)
		}
	*/
}
