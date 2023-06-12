package main

import (
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Contact struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Relations []string `json:"relations"`
}

func main() {
	// Create a Neo4j driver
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("", "", ""))
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	item, err := insertItem(driver)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Printf("%v\n", item.Message)

}

func insertItem(driver neo4j.Driver) (*Item, error) {
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	result, err := session.WriteTransaction(createContacts)
	if err != nil {
		return nil, err
	}
	return result.(*Item), nil
}
func createContacts(tx neo4j.Transaction) (interface{}, error) {
	// Create a query to create contacts and relationships
	contacts := []Contact{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john.doe@example.com",
			Phone: "123-456-7890",
			Relations: []string{
				"2", "3",
			},
		},
		{
			ID:    2,
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
			Phone: "987-654-3210",
			Relations: []string{
				"3",
			},
		},
		{
			ID:    3,
			Name:  "Alice Johnson",
			Email: "alice.johnson@example.com",
			Phone: "555-555-5555",
			Relations: []string{
				"1", "2",
			},
		},
		{
			ID:    4,
			Name:  "Dali Calla",
			Email: "test.johnson@example.com",
			Phone: "505-555-5555",
			Relations: []string{
				"2",
			},
		},
	}
	query := `
		UNWIND $contacts AS contact
		MERGE (c:Contact {id: contact.id})
		SET c.name = contact.name, c.email = contact.email, c.phone = contact.phone
		WITH c, contact.relations AS relations
		UNWIND relations AS relId
		MERGE (r:Contact {id: toInteger(relId)})
		MERGE (c)-[:RELATED_TO]->(r)
		RETURN r.id
	`

	// Create parameters map
	params := map[string]interface{}{
		"contacts": convertToMap(contacts),
	}
	// Run the query
	records, err := tx.Run(query, params)
	if err != nil {
		log.Fatal(err)
	}

	summary, err := records.Consume()
	if err != nil {
		return nil, err
	}
	// You can also retrieve values by name, with e.g. `id, found := record.Get("n.id")`
	return &Item{
		Message: fmt.Sprintf("Nodes Created: %d\nNumber of relationship created: %d", summary.Counters().NodesCreated(), summary.Counters().RelationshipsCreated()),
	}, nil
}
func convertToMap(contacts []Contact) []map[string]interface{} {
	var contactsMap []map[string]interface{}
	for _, contact := range contacts {
		contactMap := map[string]interface{}{
			"id":        contact.ID,
			"name":      contact.Name,
			"email":     contact.Email,
			"phone":     contact.Phone,
			"relations": contact.Relations,
		}
		contactsMap = append(contactsMap, contactMap)
	}
	return contactsMap
}

type Item struct {
	Message string
}
