/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"io"
	"log"

	"time"

	pb "todo/todo"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:9090"
	accessToken = "{your_idtoken}"
)

func main() {
	// サーバーへの接続のセットアップ
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoClient(conn)

	// ユーザー登録
	log.Printf("### Register User ###")
	ctxReg, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r3, err := c.RegistUser(ctxReg, &pb.UserRequest{Username: "aaa", Password: "password", Email: "test104@example.com"})
	if err != nil {
		log.Printf("could not register user: %v", err)
	}
	log.Printf("Token: %T", r3)

	// 単一Todoの登録
	token := accessToken
	log.Printf("### Register Todo ###")
	defer cancel()
	ctxRegTodo, cancel := context.WithTimeout(context.Background(), time.Second*5)
	r4, err := c.RegistTodo(ctxRegTodo, &pb.TodoRequest{Content: "test2", Token: token})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Message: %v", r4.Message)
	tid := r4.Id

	// 単一Todoの検索
	log.Printf("### Retrive Todo ###")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RetrieveTodo(ctx, &pb.SearchRequest{Token: token, Id: tid})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Content: %s", r.Content)

	// 複数Todoの登録
	log.Printf("### Register multiple Todos ###")
	ctxRegMulti, cancel := context.WithTimeout(context.Background(), time.Second*3)
	streamReg, err := c.RegistTodos(ctxRegMulti)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	contents := [3]string{"test-11", "test-12", "test-13"}
	for _, content := range contents {
		if err := streamReg.Send(&pb.TodoRequest{Content: content, Token: token}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", streamReg, "", err)
		}
	}
	resRegTodos, err := streamReg.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", streamReg, err, nil)
	}
	log.Printf("Greeting Client Stream: %s", resRegTodos.Message)

	// 複数Todoの検索
	log.Printf("### List Todos ###")
	ctxList, cancel := context.WithTimeout(context.Background(), time.Second*3)
	stream, err := c.ListTodos(ctxList, &pb.SearchRequest{Token: token})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Println(feature.Content)
	}

	// Todoの一件削除
	log.Printf("### Delete Todo ###")
	defer cancel()
	ctxDelTodo, cancel := context.WithTimeout(context.Background(), time.Second*5)
	resDel, err := c.DeleteTodo(ctxDelTodo, &pb.DeleteRequest{Id: tid, Token: token})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Message: %v", resDel.Message)

	//Todoの全件削除
	log.Printf("### Delete Todo ###")
	defer cancel()
	ctxDelTodos, cancel := context.WithTimeout(context.Background(), time.Second*5)
	resDelAll, err := c.DeleteTodoAll(ctxDelTodos, &pb.DeleteRequest{Token: token})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Message: %v", resDelAll.Message)

}
