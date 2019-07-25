package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

type Todo struct {
	client *firestore.Client
}

func newTodo(ctx context.Context, app *firebase.App) (*Todo, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error occured while set firestore. %v\n", err)
		return nil, err
	}
	return &Todo{client: client}, nil
}

func (t *Todo) createTodo(ctx context.Context, content string, uid string) (string, string, error) {
	return t.registTodo(ctx, content, uid, "")

}

func (t *Todo) registTodo(ctx context.Context, content string, uid string, tid string) (string, string, error) {
	var err error
	switch tid {
	case "":
		ref := t.client.Collection("users").Doc(uid).Collection("todos").NewDoc()
		_, err = ref.Set(ctx, map[string]interface{}{
			"content": content,
		}, firestore.MergeAll)
		tid = ref.ID
	default:
		_, err = t.client.Collection("users").Doc(uid).Collection("todos").Doc(tid).Set(ctx, map[string]interface{}{
			"content": content,
		})
	}
	if err != nil {
		log.Printf("Failed adding alovelace: %v", err)
		return tid, "登録失敗", err
	}
	log.Printf("Add Todo: %v\n", content)
	return tid, "登録成功", err
}

func (t *Todo) deleteTodo(ctx context.Context, uid string, tid string) (string, error) {
	var err error
	switch tid {
	case "":
		ref := t.client.Collection("users").Doc(uid).Collection("todos")
		limit := 100
		for {
			// Get a batch of documents
			iter := ref.Limit(limit).Documents(ctx)
			numDeleted := 0

			// Iterate through the documents, adding
			// a delete operation for each one to a
			// WriteBatch.
			batch := t.client.Batch()
			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					return "", err
				}

				batch.Delete(doc.Ref)
				numDeleted++
			}

			// If there are no documents to delete,
			// the process is over.
			if numDeleted == 0 {
				return "削除成功", nil
			}

			_, err := batch.Commit(ctx)
			if err != nil {
				return "削除失敗", err
			}
		}
	default:
		_, err = t.client.Collection("users").Doc(uid).Collection("todos").Doc(tid).Delete(ctx)
	}
	if err != nil {
		log.Printf("Failed adding alovelace: %v", err)
		return "削除失敗", err
	}

	return "削除成功", err
}

func (t *Todo) getTodo(ctx context.Context, uid string, tid string) (string, string, error) {
	dsnap, err := t.client.Collection("users").Doc(uid).Collection("todos").Doc(tid).Get(ctx)
	if err != nil {
		fmt.Printf("error_occured %v", err)
		return "", "", err
	}
	m := dsnap.Data()
	content, ok := m["content"].(string)
	if !ok {
		return "", "", nil
	}
	return content, tid, nil
}

func (t *Todo) listTodos(ctx context.Context, uid string) ([]string, error) {
	contents := []string{}
	iter := t.client.Collection("users").Doc(uid).Collection("todos").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		content, ok := doc.Data()["content"].(string)
		if ok {
			contents = append(contents, content)
		}
	}
	return contents, nil
}
