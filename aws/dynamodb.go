package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// DynamoDBService expose
type DynamoDBService interface {
	List(field string, value string, castTo interface{}) error
	Get(key string, castTo interface{}) error
	Store(item interface{}) error
	TxStore(item interface{}) error
	TxGet(keys []string, castTo interface{}) error
}

type (
	// DynamoDB aggregator
	DynamoDB struct {
		table string
		conn  dynamodbiface.DynamoDBAPI
	}
	// DynamoDBTableAttribute composing attributes and elements
	DynamoDBTableAttribute struct {
		Name    string
		Type    string
		KeyType string
	}
)

// NewDynamoDB - creates new dynamodb instance with connection
func NewDynamoDB(region, endpoint, table string) (ddb *DynamoDB, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		return nil, err
	}
	ddb = &DynamoDB{
		conn: dynamodb.New(sess), table: table,
	}

	// Single table design https://www.trek10.com/blog/dynamodb-single-table-relational-modeling/
	err = ddb.createTable(table, []DynamoDBTableAttribute{
		{
			Name:    "PK",
			Type:    "S",
			KeyType: "HASH",
		},
		{
			Name:    "SK",
			Type:    "S",
			KeyType: "RANGE",
		},
	})

	return
}

// List gets a collection of resources
func (ddb *DynamoDB) List(field string, value string, castTo interface{}) error {
	results, err := ddb.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(ddb.table),
		ScanFilter: map[string]*dynamodb.Condition{
			field: {
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(value),
					},
				},
				ComparisonOperator: aws.String("EQ"),
			},
		},
	})
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalListOfMaps(results.Items, &castTo)
}

// Store an Item
func (ddb *DynamoDB) Store(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ddb.table),
	}
	_, err = ddb.conn.PutItem(input)
	return err
}

// Get an item
func (ddb *DynamoDB) Get(key string, castTo interface{}) error {
	result, err := ddb.conn.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ddb.table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return err
	}
	return dynamodbattribute.UnmarshalMap(result.Item, &castTo)
}

// TxStore writes an item with transaction
func (ddb *DynamoDB) TxStore(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	_, err = ddb.conn.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			&dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					TableName: aws.String(ddb.table),
					Item:      av,
				},
			},
		},
	})
	return err
}

// TxGet gets items with transaction
func (ddb *DynamoDB) TxGet(keys []string, castTo interface{}) error {
	var avs []*dynamodb.TransactGetItem
	for _, key := range keys {
		avs = append(avs, &dynamodb.TransactGetItem{
			Get: &dynamodb.Get{
				TableName: aws.String(ddb.table),
				Key: map[string]*dynamodb.AttributeValue{
					"PK": {
						S: aws.String(key),
					},
				},
			},
		})
	}

	res, err := ddb.conn.TransactGetItems(&dynamodb.TransactGetItemsInput{
		TransactItems: avs,
	})
	if err != nil {
		return err
	}

	var outputItems []map[string]*dynamodb.AttributeValue
	for _, item := range res.Responses {
		outputItems = append(outputItems, item.Item)
	}

	return dynamodbattribute.UnmarshalListOfMaps(outputItems, &castTo)
}

// CreateTable a DynamoDB table
func (ddb *DynamoDB) createTable(tableName string, atts []DynamoDBTableAttribute) (err error) {
	var (
		eles []*dynamodb.KeySchemaElement
		defs []*dynamodb.AttributeDefinition
	)

	for _, att := range atts {
		keyType := strings.ToUpper(att.KeyType)
		// TODO: check KeyType
		defs = append(defs,
			&dynamodb.AttributeDefinition{
				AttributeName: aws.String(att.Name),
				AttributeType: aws.String(att.Type),
			},
		)
		eles = append(eles,
			&dynamodb.KeySchemaElement{
				AttributeName: aws.String(att.Name),
				KeyType:       aws.String(keyType),
			},
		)
	}
	_, err = ddb.conn.CreateTable(&dynamodb.CreateTableInput{
		TableName:            aws.String(tableName),
		BillingMode:          aws.String("PAY_PER_REQUEST"),
		KeySchema:            eles,
		AttributeDefinitions: defs,
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{{
			IndexName: aws.String("gsi_1"),
			KeySchema: eles,
			Projection: &dynamodb.Projection{
				ProjectionType: aws.String("ALL"),
			},
		}},
	})
	// Skip when table exists
	if err != nil {
		if strings.Contains(err.Error(), "ResourceInUseException") {
			return nil
		}
	}
	return
}
