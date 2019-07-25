package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

const ()

func existInArray(inputMethod string) bool {
	execuldeMethodFromInterceptor := []string{
		"/todo.Todo/RegistUser",
	}

	for num := range execuldeMethodFromInterceptor {
		if inputMethod == execuldeMethodFromInterceptor[num] {
			return true
		}
	}
	return false
}

/*
func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	s, ok := info.Server.(*TodoServer)
	if !ok {
		return nil, nil
	}
	if !existInArray(info.FullMethod) {
		//		ensureJWT(s.client_auth, req)
	}
	// 2. リクエストされたgrpcメソッドを実行
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, nil
	}

	// 3. response bodyをprint
	fmt.Printf("Interceptor End\n")
	return resp, nil

}*/

/*func StreamServerInterceptor(authFunc AuthFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := srv.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, err = authFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}*/

type Auth struct {
	client *auth.Client
}

func newAuth(ctx context.Context, app *firebase.App) (*Auth, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Error occured while set auth. %v\n", err)
		return nil, err
	}
	return &Auth{client: client}, nil
}

func (a *Auth) ensureJWT(ctx context.Context, idToken string) (string, error) {
	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return "", err
	}
	return token.UID, nil
}

func (a *Auth) createUser(ctx context.Context, username string, password string, email string) (string, error) {
	// [START create_user_golang]
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Password(password).
		DisplayName(username)
	u, err := a.client.CreateUser(ctx, params)
	if err != nil {
		return "", err
	}

	uName := u.UserInfo.DisplayName
	// [END create_user_golang]
	return uName, nil
}
