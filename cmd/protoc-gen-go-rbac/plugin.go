package main

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	rbac "github.com/casnerano/protoc-gen-go-rbac/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	outputFileSuffix = ".rbac.go"
)

//go:embed plugin.go.tmpl
var templateFS embed.FS

type TemplateData struct {
	Meta     Meta
	File     File
	Services []Service
}

type Meta struct {
	ModuleVersion string
	ProtocVersion string
}

type File struct {
	Name    string
	Package string
	Source  string
}

type Service struct {
	Name    string
	Rules   *rbac.Rules
	Methods []Method
}

type Method struct {
	Name  string
	Rules *rbac.Rules
}

func execute(plugin *protogen.Plugin) error {
	tmpl, err := parseTemplate()
	if err != nil {
		return err
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		services := collectServices(file.Services)
		if len(services) == 0 {
			continue
		}

		templateData := TemplateData{
			Meta: Meta{
				ProtocVersion: func() string {
					if ver := plugin.Request.CompilerVersion; ver != nil {
						return fmt.Sprintf("v%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
					}
					return "(unknown)"
				}(),
				ModuleVersion: "(unknown)",
			},
			File: File{
				Name:    filepath.Base(file.GeneratedFilenamePrefix),
				Package: string(file.GoPackageName),
				Source:  file.Desc.Path(),
			},
			Services: services,
		}

		filename := file.GeneratedFilenamePrefix + outputFileSuffix
		if err = tmpl.Execute(plugin.NewGeneratedFile(filename, file.GoImportPath), templateData); err != nil {
			return fmt.Errorf("failed execute template: %w", err)
		}
	}

	return nil
}

func collectServices(protoServices []*protogen.Service) []Service {
	var services []Service
	for _, protoService := range protoServices {
		if options := protoService.Desc.Options().(*descriptorpb.ServiceOptions); options != nil {
			if serviceRules := proto.GetExtension(options, rbac.E_ServiceRules).(*rbac.Rules); serviceRules != nil {
				services = append(services, Service{
					Name:    string(protoService.Desc.Name()),
					Rules:   serviceRules,
					Methods: collectMethods(protoService.Methods),
				})
			}
		}
	}

	return services
}

func collectMethods(protoMethods []*protogen.Method) []Method {
	var methods []Method
	for _, protoMethod := range protoMethods {
		if options := protoMethod.Desc.Options().(*descriptorpb.MethodOptions); options != nil {
			if methodRules := proto.GetExtension(options, rbac.E_MethodRules).(*rbac.Rules); methodRules != nil {
				methods = append(methods, Method{
					Name:  string(protoMethod.Desc.Name()),
					Rules: methodRules,
				})
			}
		}
	}

	return methods
}

func parseTemplate() (*template.Template, error) {
	templateContent, err := templateFS.ReadFile("plugin.go.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	tmpl := template.New("plugin.rbac").
		Funcs(template.FuncMap{
			"toLower": strings.ToLower,
		})

	return tmpl.Parse(string(templateContent))
}
