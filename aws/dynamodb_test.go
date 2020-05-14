package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type stubDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *stubDynamoDB) CreateTable(input *dynamodb.CreateTableInput) (output *dynamodb.CreateTableOutput, err error) {
	return
}

func (m *stubDynamoDB) PutItem(input *dynamodb.PutItemInput) (output *dynamodb.PutItemOutput, err error) {
	return
}

func (m *stubDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	resp := map[string]*dynamodb.AttributeValue{
		"pk": &dynamodb.AttributeValue{
			S: aws.String("key"),
		},
		"msg": &dynamodb.AttributeValue{
			S: aws.String("sample-msg"),
		},
	}
	output := &dynamodb.GetItemOutput{
		Item: resp,
	}
	return output, nil
}

func (m *stubDynamoDB) TransactGetItems(input *dynamodb.TransactGetItemsInput) (*dynamodb.TransactGetItemsOutput, error) {
	output := &dynamodb.TransactGetItemsOutput{
		Responses: []*dynamodb.ItemResponse{
			{
				Item: map[string]*dynamodb.AttributeValue{
					"pk": &dynamodb.AttributeValue{
						S: aws.String("key"),
					},
					"msg": &dynamodb.AttributeValue{
						S: aws.String("sample-msg"),
					},
				},
			},
		},
	}
	return output, nil
}

func (m *stubDynamoDB) TransactWriteItems(input *dynamodb.TransactWriteItemsInput) (output *dynamodb.TransactWriteItemsOutput, err error) {
	return
}

func (m *stubDynamoDB) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	resp := []map[string]*dynamodb.AttributeValue{
		{
			"pk": &dynamodb.AttributeValue{
				S: aws.String("key"),
			},
			"msg": &dynamodb.AttributeValue{
				S: aws.String("sample-msg"),
			},
		},
		{
			"pk": &dynamodb.AttributeValue{
				S: aws.String("key2"),
			},
			"msg": &dynamodb.AttributeValue{
				S: aws.String("sample-msg2"),
			},
		},
	}
	output := &dynamodb.ScanOutput{
		Items: resp,
	}
	return output, nil
}

type data struct {
	PK  string
	Msg string
}

func TestDDB(t *testing.T) {
	ddb := DynamoDB{
		table: "testTable",
		conn:  &stubDynamoDB{},
	}
	err := ddb.createTable(ddb.table, []DynamoDBTableAttribute{
		{
			Name:    "pk",
			Type:    "S",
			KeyType: "HASH",
		},
		{
			Name:    "type",
			Type:    "S",
			KeyType: "RANGE",
		},
	})

	if err != nil {
		t.Errorf("Unable to create a table " + err.Error())
	}

	if err := ddb.Store(data{}); err != nil {
		t.Errorf("Unable to save " + err.Error())
	}

	var dd []data
	if err := ddb.List("foo", "bar", &dd); err != nil {
		t.Errorf("Unable to get list " + err.Error())
	}
	if len(dd) != 2 {
		t.Errorf("Incorrect amount of data returned")
	}

	var d data
	if err := ddb.Get("key", &d); err != nil {
		t.Errorf("Unable to get data " + err.Error())
	}
	if d.Msg != "sample-msg" {
		t.Errorf("Incorrect data returned " + d.Msg)
	}

	if err := ddb.TxStore(data{}); err != nil {
		t.Errorf("Unable to save " + err.Error())
	}

	if err := ddb.TxGet([]string{}, dd); err != nil {
		t.Errorf("Unable to get " + err.Error())
	}
	if dd[0].PK != "key" {
		t.Errorf("Incorrect data")
	}
}
