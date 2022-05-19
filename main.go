package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"log"
	"os"
)

const version = "v0.0.1"

var (
	showVersion = flag.Bool("version", false, "print the version and exit")
	omitempty   = flag.Bool("omitempty", true, "omit if google.api is empty")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gin %v\n", version)
		return
	}

	//ops := protogen.Options{
	//	//ParamFunc: flag.CommandLine.Set,
	//}
	//err := run(ops, func(gen *protogen.Plugin) error {
	//	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	//	for _, f := range gen.Files {
	//		if !f.Generate {
	//			continue
	//		}
	//		generateFile(gen, f, *omitempty)
	//	}
	//	return nil
	//})
	//if err != nil {
	//	log.Println(err)
	//}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, *omitempty)
		}
		return nil
	})
}

func run(opts protogen.Options, f func(*protogen.Plugin) error) error {
	req := &pluginpb.CodeGeneratorRequest{}
	bi, err := os.Open("./code_generator_request.pb.bin")
	if err != nil {
		log.Fatalln(err)
	}
	in, _ := ioutil.ReadAll(bi)
	if err := proto.Unmarshal(in, req); err != nil {
		return err
	}
	gen, err := opts.New(req)
	if err != nil {
		return err
	}
	if err := f(gen); err != nil {
		gen.Error(err)
	}
	resp := gen.Response()
	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}
