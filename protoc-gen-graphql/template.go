package main

var goTemplate = `
// Code generated by proroc-gen-graphql, DO NOT EDIT.
package {{ .RootPackage.Name }}

import (
	"context"

	"github.com/graphql-go/graphql"
{{- if .Services }}
	"github.com/ysugimoto/grpc-graphql-gateway/runtime"
	"google.golang.org/grpc"
{{- end }}

{{- range .Packages }}
	{{ if .Path }}{{ .Name }} "{{ .Path }}"{{ end }}
{{ end }}
)

var (
	{{- range .Enums }}
	gql__enum_{{ .Name }} *graphql.Enum // enum {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Interfaces }}
	gql__interface_{{ .Name }} *graphql.Interface // message {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Types }}
	gql__type_{{ .TypeName }} *graphql.Object // message {{ .Name }} in {{ .Filename }}
	{{- end }}
	{{- range .Inputs }}
	gql__input_{{ .TypeName }} *graphql.InputObject // message {{ .Name }} in {{ .Filename }}
	{{- end }}
)

{{ range .Enums -}}
func Gql__enum_{{ .Name }}() *graphql.Enum {
	if gql__enum_{{ .Name }} == nil {
		gql__enum_{{ .Name }} =  graphql.NewEnum(graphql.EnumConfig{
			Name: "{{ .Name }}",
			Values: graphql.EnumValueConfigMap{
{{- range .Values }}
				"{{ .Name }}": &graphql.EnumValueConfig{
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					Value: {{ .Number }},
				},
{{- end }}
			},
		})
	}
	return gql__enum_{{ .Name }}
}

{{ end }}

{{ range .Interfaces -}}
func Gql__interface_{{ .TypeName }}() *graphql.Interface {
	if gql__interface_{{ .TypeName }} == nil {
		gql__interface_{{ .TypeName }} =  graphql.NewInterface(graphql.InterfaceConfig{
			Name: "{{ .TypeName }}Interface",
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Fields: graphql.Fields{
{{- range .Fields }}
			{{- if not .IsCyclic }}
				"{{ .Name }}": &graphql.Field{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
				},
			{{- end }}
{{- end }}
			},
			ResolveType: func(p graphql.ResolveTypeParams) *graphql.Object {
				return Gql__type_{{ .TypeName }}()
			},
		})
	}
	return gql__interface_{{ .TypeName }}
}

{{ end }}

{{ range .Types -}}
func Gql__type_{{ .TypeName }}() *graphql.Object {
	if gql__type_{{ .TypeName }} == nil {
		gql__type_{{ .TypeName }} =  graphql.NewObject(graphql.ObjectConfig{
			Name: "{{ .TypeName }}",
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Fields: graphql.Fields {
{{- range .Fields }}
				"{{ .Name }}": &graphql.Field{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
				},
{{- end }}
			},
			{{- if .Interfaces }}
			Interfaces: []*graphql.Interface{
{{- range .Interfaces }}
			Gql__interface_{{ .TypeName }}(),
{{- end }}
			},
			{{- end }}
		})
	}
	return gql__type_{{ .TypeName }}
}

{{ end }}

{{ range .Inputs -}}
func Gql__input_{{ .TypeName }}() *graphql.InputObject {
	if gql__input_{{ .TypeName }} == nil {
		gql__input_{{ .TypeName }} =  graphql.NewInputObject(graphql.InputObjectConfig{
			Name: "{{ .TypeName }}",
			Fields: graphql.InputObjectConfigFieldMap{
{{- range .Fields }}
				"{{ .Name }}": &graphql.InputObjectFieldConfig{
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					Type: {{ .FieldTypeInput $.RootPackage.Path }},
				},
{{- end }}
			},
		})
	}
	return gql__input_{{ .TypeName }}
}

{{ end }}

{{ range $_, $service := .Services -}}
// graphql__resolver_{{ $service.Name }} is a struct for making query, mutation and resolve fields.
// This struct must be implemented runtime.SchemaBuilder interface.
type graphql__resolver_{{ $service.Name }} struct {

	// Automatic connection host
	host string

	// grpc dial options
	dialOptions []grpc.DialOption

	// grpc client connection.
	// this connection may be provided by user
	conn *grpc.ClientConn
}

// CreateConnection() returns grpc connection which user specified or newly connected and closing function
func (x *graphql__resolver_{{ $service.Name }}) CreateConnection(ctx context.Context) (*grpc.ClientConn, func(), error) {
	// If x.conn is not nil, user injected their own connection
	if x.conn != nil {
		return x.conn, func() {}, nil
	}

	// Otherwise, this handler opens connection with specified host
	conn, err := grpc.DialContext(ctx, x.host, x.dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	return conn, func() { conn.Close() }, nil
}

// GetQueries returns acceptable graphql.Fields for Query.
func (x *graphql__resolver_{{ $service.Name }}) GetQueries(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
{{- range .Queries }}
		"{{ .QueryName }}": &graphql.Field{
			Type: {{ .QueryType }},
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{- end }}
			Args: graphql.FieldConfigArgument{
			{{- range .Args }}
				"{{ .Name }}": &graphql.ArgumentConfig{
					Type: {{ .FieldType $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					{{- if .DefaultValue }}
					DefaultValue: {{ .DefaultValue }},
					{{- end }}
				},
			{{- end }}
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var req *{{ .InputType }}
				if err := runtime.MarshalRequest(p.Args, &req); err != nil {
					return nil, err
				}
				client := {{ .Package }}New{{ $service.Name }}Client(conn)
				resp, err := client.{{ .Method.Name }}(p.Context, req)
				if err != nil {
					return nil, err
				}
				{{- if .IsPluckResponse }}
				return resp.Get{{ .PluckResponseFieldName }}(), nil
				{{- else }}
				return resp, nil
				{{- end }}
			},
		},
{{- end }}
	}
}

// GetMutations returns acceptable graphql.Fields for Mutation.
func (x *graphql__resolver_{{ $service.Name }}) GetMutations(conn *grpc.ClientConn) graphql.Fields {
	return graphql.Fields{
{{- range .Mutations }}
		"{{ .MutationName }}": &graphql.Field{
			Type: {{ .MutationType }},
			{{- if .Comment }}
			Description: ` + "`" + `{{ .Comment }}` + "`" + `,
			{{ end }}
			Args: graphql.FieldConfigArgument{
			{{- range .Args }}
				"{{ .Name }}": &graphql.ArgumentConfig{
					Type: {{ .FieldTypeInput $.RootPackage.Path }},
					{{- if .Comment }}
					Description: ` + "`" + `{{ .Comment }}` + "`" + `,
					{{- end }}
					{{- if .DefaultValue }}
					DefaultValue: {{ .DefaultValue }},
					{{- end }}
				},
			{{- end }}
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var req *{{ .InputType }}
				if err := runtime.MarshalRequest(p.Args, &req); err != nil {
					return nil, err
				}
				client := {{ .Package }}New{{ $service.Name }}Client(conn)
				resp, err := client.{{ .Method.Name }}(p.Context, req)
				if err != nil {
					return nil, err
				}
				{{- if .IsPluckResponse }}
				return resp.Get{{ .PluckResponseFieldName }}(), nil
				{{- else }}
				return resp, nil
				{{- end }}
			},
		},
{{ end }}
	}
}

// Register package divided graphql handler "without" *grpc.ClientConn,
// therefore gRPC connection will be opened and closed automatically.
// Occasionally you may worry about open/close performance for each handling graphql request,
// then you can call Register{{ .Name }}GraphqlHandler with *grpc.ClientConn manually.
func Register{{ .Name }}Graphql(mux *runtime.ServeMux) error {
	return Register{{ .Name }}GraphqlHandler(mux, nil)
}

// Register package divided graphql handler "with" *grpc.ClientConn.
// this function accepts your defined grpc connection, so that we reuse that and never close connection inside.
// You need to close it maunally when application will terminate.
// Otherwise, you can specify automatic opening connection with ServiceOption directive:
//
// service {{ .Name }} {
//    option (graphql.service) = {
//        host: "host:port"
//        insecure: true or false
//    };
//
//    ...with RPC definitions
// }
func Register{{ .Name }}GraphqlHandler(mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return mux.AddHandler(&graphql__resolver_{{ .Name }}{
		conn: conn,
		host: "{{ if .Host }}{{ .Host }}{{ else }}localhost:50051{{ end }}",
		dialOptions: []grpc.DialOption{
		{{- if .Insecure }}
			grpc.WithInsecure(),
		{{- end }}
		},
	})
}

{{ end }}
`
