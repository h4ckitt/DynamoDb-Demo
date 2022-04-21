package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {

	/*cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.E)
	if err != nil {
		log.Fatal(err)
	}*/

	opts := dynamodb.Options{
		Region:           "us-east-1",
		EndpointResolver: dynamodb.EndpointResolverFromURL("http://localhost:4566"),
	}

	//sess := session.Must(session.NewSessionWithOptions(session.Options{SharedConfigState: session.SharedConfigEnable}))

	svc := dynamodb.New(opts)

	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})

	if err != nil {
		log.Fatalf("Failed To List Tables: %v", err)
	}

	if len(resp.TableNames) == 0 {
		log.Fatalln("Nope, No Tables In Sight Chief")
	}

	fmt.Println("Tables:")
	for _, name := range resp.TableNames {
		fmt.Println(name)
	}

	anotherResp, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(resp.TableNames[0]),
	})

	if err != nil {
		log.Fatalln("A Fatal Error Occurred")
	}

	x := anotherResp.Table

	for _, x := range x.AttributeDefinitions {
		fmt.Printf("%s -> %v \n", *(x.AttributeName), x.AttributeType)
	}

	for _, key := range x.KeySchema {
		fmt.Printf("%s -> %v\n", *(key.AttributeName), key.KeyType)
	}

	testInput := struct {
		Artist     string
		SongTitle  string
		AlbumTitle string
		Awards     int
	}{
		Artist:     "The Beach Boys",
		SongTitle:  "Does This Work",
		AlbumTitle: "Some Random Album",
		Awards:     10,
	}

	input, err := attributevalue.MarshalMap(testInput)

	fmt.Println("Successfully Promoted To Marshal")

	if err != nil {
		log.Fatalln("Invalid Input Somewhere In Your String")
	}

	finalInput := &dynamodb.PutItemInput{
		Item:                input,
		TableName:           aws.String(resp.TableNames[0]),
		ConditionExpression: aws.String("attribute_not_exists(SongTitle)"),
	}

	fmt.Printf("%v\n", finalInput)

	fmt.Println("Successfully Created The Final Input")

	_, err = svc.PutItem(context.TODO(), finalInput)

	if err != nil {
		log.Fatalln("Welp Some Shit Happened: ", err.Error())
	}

	fmt.Println("Suksexful")
	//fmt.Println(x.AttributeDefinitions)
}
