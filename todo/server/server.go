package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync"
	pb "todo/todo"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

const (
	port = ":9090"
)

type TodoServer struct {
	mu   sync.Mutex
	auth *Auth
	todo *Todo
}

func initFirebase() (*Auth, *Todo) {
	opt := option.WithCredentialsFile("{/path/to/your-service-account-file.json}")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error occured while set app. %v\n", err)
	}
	auth, err := newAuth(ctx, app)
	if err != nil {
		log.Fatalf("Error occured while set auth. %v\n", err)
	}
	todo, err := newTodo(ctx, app)
	if err != nil {
		log.Fatalf("Error occured while set firestore. %v\n", err)
	}
	return auth, todo
}

func (s *TodoServer) RegistUser(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	uName := in.GetUsername()
	password := in.GetPassword()
	addr := in.GetEmail()
	retUName, err := s.auth.createUser(ctx, uName, password, addr)
	if err != nil {
		return &pb.UserReply{Username: ""}, err
	}
	return &pb.UserReply{Username: retUName}, nil
}

func (s *TodoServer) RegistTodo(ctx context.Context, in *pb.TodoRequest) (*pb.TodoReply, error) {
	token := in.GetToken()
	uid, err := s.auth.ensureJWT(ctx, token)
	if err != nil {
		return &pb.TodoReply{Id: "", Content: "", Message: "JWT検証失敗による登録エラー"}, err
	}
	rcvContent := in.GetContent()
	tid := in.GetTodoId()
	tid, result, err := s.todo.registTodo(ctx, rcvContent, uid, tid)
	if err != nil {
		return &pb.TodoReply{Id: "", Content: "", Message: result}, err
	}
	return &pb.TodoReply{Id: tid, Content: rcvContent, Message: result}, nil
}

func (s *TodoServer) RetrieveTodo(ctx context.Context, in *pb.SearchRequest) (*pb.TodoReply, error) {
	token := in.GetToken()
	uid, err := s.auth.ensureJWT(ctx, token)
	if err != nil {
		return &pb.TodoReply{Content: "", Id: "", Message: "JWT検証失敗による検索エラー"}, err
	}
	tid := in.GetId()
	content, tid, err := s.todo.getTodo(ctx, uid, tid)
	if err != nil {
		return &pb.TodoReply{Content: content, Id: tid, Message: ""}, nil
	}
	return &pb.TodoReply{Content: content, Id: tid, Message: ""}, nil
}

func (s *TodoServer) ListTodos(in *pb.SearchRequest, stream pb.Todo_ListTodosServer) error {
	ctx := context.Background()
	token := in.GetToken()
	uid, err := s.auth.ensureJWT(ctx, token)
	if err != nil {
		if errStream := stream.Send(&pb.TodoReply{Content: "", Id: "", Message: "JWT検証失敗による検索エラー"}); errStream != nil {
			return errStream
		}
		return nil
	}
	todoList, err := s.todo.listTodos(ctx, uid)
	for _, content := range todoList {
		if errStream := stream.Send(&pb.TodoReply{Content: content}); errStream != nil {
			return errStream
		}
	}
	return nil
}

func (s *TodoServer) RegistTodos(stream pb.Todo_RegistTodosServer) error {
	ctx := context.Background()
	for {
		todo, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ResultReply{Message: "登録完了"})
		}
		token := todo.GetToken()
		uid, err := s.auth.ensureJWT(ctx, token)
		if err != nil {
			if errStream := stream.SendAndClose(&pb.ResultReply{Message: "JWT検証失敗による登録エラー"}); errStream != nil {
				return errStream
			}
			return nil
		}
		rcvContent := todo.GetContent()
		_, _, err = s.todo.registTodo(ctx, rcvContent, uid, "")
		if err != nil {
			return err
		}
	}
}

func (s *TodoServer) DeleteTodo(ctx context.Context, in *pb.DeleteRequest) (*pb.ResultReply, error) {
	token := in.GetToken()
	uid, err := s.auth.ensureJWT(ctx, token)
	if err != nil {
		return &pb.ResultReply{Message: "JWT検証失敗による登録エラー"}, err
	}
	tid := in.GetId()
	if tid == "" {
		return &pb.ResultReply{Message: "TodoのIDは必須です"}, err
	}
	result, err := s.todo.deleteTodo(ctx, uid, tid)
	if err != nil {
		return &pb.ResultReply{Message: result}, err
	}
	return &pb.ResultReply{Message: result}, nil
}

func (s *TodoServer) DeleteTodoAll(ctx context.Context, in *pb.DeleteRequest) (*pb.ResultReply, error) {
	token := in.GetToken()
	uid, err := s.auth.ensureJWT(ctx, token)
	if err != nil {
		return &pb.ResultReply{Message: "JWT検証失敗による登録エラー"}, err
	}
	result, err := s.todo.deleteTodo(ctx, uid, "")
	if err != nil {
		return &pb.ResultReply{Message: result}, err
	}
	return &pb.ResultReply{Message: result}, nil
}

func (s *TodoServer) DeleteUser(ctx context.Context, in *pb.DeleteRequest) (*pb.ResultReply, error) {
	return nil, nil
}

func (s *TodoServer) LoginUser(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	return nil, nil
}

func newServer() *TodoServer {
	auth, todo := initFirebase()
	s := &TodoServer{auth: auth, todo: todo}
	return s
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//	s := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor))
	s := grpc.NewServer()
	pb.RegisterTodoServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
