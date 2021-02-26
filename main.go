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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	pb "combination.se/apm/helloworld"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	port = ":50001"
)

var c pb.GreeterClient

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return c.SayHello(ctx, in)
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello2(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// 172.16.31.75:5000
	// localhost:50001
	conn, err := grpc.Dial("172.16.31.75:5000", grpc.WithInsecure(), grpc.WithChainUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewGreeterClient(conn)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor()))
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		md, _ := metadata.FromOutgoingContext(ctx)
		md2 := metadata.New(nil)
		//md2.Set("elastic-apm-traceparent", md.Get("elastic-apm-traceparent")...)
		md2.Set("traceparent", md.Get("traceparent")...)
		md2.Set("grpc-accept-encoding", "identity,gzip")
		ctx = metadata.NewOutgoingContext(ctx, md2)

		err := invoker(ctx, method, req, resp, cc, opts...)

		return err
	}
}
