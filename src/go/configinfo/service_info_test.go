// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configinfo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/GoogleCloudPlatform/esp-v2/src/go/options"
	"github.com/GoogleCloudPlatform/esp-v2/src/go/util"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"

	commonpb "github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/http/common"
	scpb "github.com/GoogleCloudPlatform/esp-v2/src/go/proto/api/envoy/http/service_control"
	annotationspb "google.golang.org/genproto/googleapis/api/annotations"
	confpb "google.golang.org/genproto/googleapis/api/serviceconfig"
	apipb "google.golang.org/genproto/protobuf/api"
)

var (
	testProjectName = "bookstore.endpoints.project123.cloud.goog"
	testApiName     = "endpoints.examples.bookstore.Bookstore"
	testConfigID    = "2019-03-02r0"
)

func TestProcessEndpoints(t *testing.T) {
	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		wantedAllowCors   bool
	}{
		{
			desc: "Return true for endpoint name matching service name",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Endpoints: []*confpb.Endpoint{
					{
						Name:      testProjectName,
						AllowCors: true,
					},
				},
			},
			wantedAllowCors: true,
		},
		{
			desc: "Return false for not setting allow_cors",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Endpoints: []*confpb.Endpoint{
					{
						Name: testProjectName,
					},
				},
			},
			wantedAllowCors: false,
		},
		{
			desc: "Return false for endpoint name not matching service name",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Endpoints: []*confpb.Endpoint{
					{
						Name:      "echo.endpoints.project123.cloud.goog",
						AllowCors: true,
					},
				},
			},
			wantedAllowCors: false,
		},
		{
			desc: "Return false for empty endpoint field",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
			},
			wantedAllowCors: false,
		},
	}

	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = "grpc"
		serviceInfo, err := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)
		if err != nil {
			t.Fatal(err)
		}

		if serviceInfo.AllowCors != tc.wantedAllowCors {
			t.Errorf("Test Desc(%d): %s, allow CORS flag got: %v, want: %v", i, tc.desc, serviceInfo.AllowCors, tc.wantedAllowCors)
		}
	}
}

func TestExtractAPIKeyLocations(t *testing.T) {
	testData := []struct {
		desc                   string
		fakeServiceConfig      *confpb.Service
		wantedSystemParameters map[string][]*confpb.SystemParameter
		wantMethods            map[string]*methodInfo
	}{
		{
			desc: "Succeed, only url query",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "echo",
							},
						},
					},
				},
				SystemParameters: &confpb.SystemParameters{
					Rules: []*confpb.SystemParameterRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo",
							Parameters: []*confpb.SystemParameter{
								{
									Name:       "api_key",
									HttpHeader: "header_name",
								},
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo": &methodInfo{
					ShortName: "echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					APIKeyLocations: []*scpb.APIKeyLocation{
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name",
							},
						},
					},
				},
			},
		},
		{
			desc: "Succeed, only header",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "echo",
							},
						},
					},
				},
				SystemParameters: &confpb.SystemParameters{
					Rules: []*confpb.SystemParameterRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo",
							Parameters: []*confpb.SystemParameter{
								{
									Name:              "api_key",
									UrlQueryParameter: "query_name",
								},
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo": &methodInfo{
					ShortName: "echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					APIKeyLocations: []*scpb.APIKeyLocation{
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name",
							},
						},
					},
				},
			},
		},
		{
			desc: "Succeed, url query plus header",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "echo",
							},
						},
					},
				},
				SystemParameters: &confpb.SystemParameters{
					Rules: []*confpb.SystemParameterRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo",
							Parameters: []*confpb.SystemParameter{
								{
									Name:              "api_key",
									HttpHeader:        "header_name_1",
									UrlQueryParameter: "query_name_1",
								},
								{
									Name:              "api_key",
									HttpHeader:        "header_name_2",
									UrlQueryParameter: "query_name_2",
								},
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.echo": &methodInfo{
					ShortName: "echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					APIKeyLocations: []*scpb.APIKeyLocation{
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_2",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_2",
							},
						},
					},
				},
			},
		},

		{
			desc: "Succeed, url query plus header for multiple apis",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "foo",
							},
						},
					},
					{
						Name: "2.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "bar",
							},
						},
					},
				},
				SystemParameters: &confpb.SystemParameters{
					Rules: []*confpb.SystemParameterRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.foo",
							Parameters: []*confpb.SystemParameter{
								{
									Name:              "api_key",
									HttpHeader:        "header_name_1",
									UrlQueryParameter: "query_name_1",
								},
								{
									Name:              "api_key",
									HttpHeader:        "header_name_2",
									UrlQueryParameter: "query_name_2",
								},
							},
						},
						{
							Selector: "2.echo_api_endpoints_cloudesf_testing_cloud_goog.bar",
							Parameters: []*confpb.SystemParameter{
								{
									Name:              "api_key",
									HttpHeader:        "header_name_1",
									UrlQueryParameter: "query_name_1",
								},
								{
									Name:              "api_key",
									HttpHeader:        "header_name_2",
									UrlQueryParameter: "query_name_2",
								},
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.foo": &methodInfo{
					ShortName: "foo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					APIKeyLocations: []*scpb.APIKeyLocation{
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_2",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_2",
							},
						},
					},
				},
				"2.echo_api_endpoints_cloudesf_testing_cloud_goog.bar": &methodInfo{
					ShortName: "bar",
					ApiName:   "2.echo_api_endpoints_cloudesf_testing_cloud_goog",
					APIKeyLocations: []*scpb.APIKeyLocation{
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Query{
								Query: "query_name_2",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_1",
							},
						},
						{
							Key: &scpb.APIKeyLocation_Header{
								Header: "header_name_2",
							},
						},
					},
				},
			},
		},
	}
	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = "grpc"
		serviceInfo, err := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(serviceInfo.Methods) != len(tc.wantMethods) {
			t.Errorf("Test Desc(%d): %s, got: %v, wanted: %v", i, tc.desc, serviceInfo.Methods, tc.wantMethods)
		}
		for key, gotMethod := range serviceInfo.Methods {
			wantMethod := tc.wantMethods[key]
			if eq := cmp.Equal(gotMethod, wantMethod, cmp.Comparer(proto.Equal)); !eq {
				t.Errorf("Test Desc(%d): %s, \ngot: %v,\nwanted: %v", i, tc.desc, gotMethod, wantMethod)
			}
		}
	}
}

func TestMethods(t *testing.T) {
	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		backendProtocol   string
		healthz           string
		wantMethods       map[string]*methodInfo
		wantError         string
	}{
		{
			desc: "Succeed for gRPC, no Http rule, with Healthz",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: testApiName,
						Methods: []*apipb.Method{
							{
								Name: "ListShelves",
							},
							{
								Name: "CreateShelf",
							},
						},
					},
				},
			},
			backendProtocol: "gRPC",
			healthz:         "/",
			wantMethods: map[string]*methodInfo{
				fmt.Sprintf("%s.%s", testApiName, "ListShelves"): &methodInfo{
					ShortName: "ListShelves",
					ApiName:   "endpoints.examples.bookstore.Bookstore",
				},
				fmt.Sprintf("%s.%s", testApiName, "CreateShelf"): &methodInfo{
					ShortName: "CreateShelf",
					ApiName:   "endpoints.examples.bookstore.Bookstore",
				},
				"ESPv2.HealthCheck": &methodInfo{
					ShortName:          "HealthCheck",
					ApiName:            "ESPv2",
					SkipServiceControl: true,
					IsGenerated:        true,
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/",
							HttpMethod:  util.GET,
						},
					},
				},
			},
		},
		{
			desc: "Succeed for HTTP, with Healthz",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "Echo",
							},
							{
								Name: "Echo_Auth_Jwt",
							},
						},
					},
				},
				Http: &annotationspb.Http{
					Rules: []*annotationspb.HttpRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt",
							Pattern: &annotationspb.HttpRule_Get{
								Get: "/auth/info/googlejwt",
							},
						},
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/echo",
							},
							Body: "message",
						},
					},
				},
			},
			backendProtocol: "http2",
			healthz:         "/",
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo": &methodInfo{
					ShortName: "Echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/echo",
							HttpMethod:  util.POST,
						},
					},
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt": &methodInfo{
					ShortName: "Echo_Auth_Jwt",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/auth/info/googlejwt",
							HttpMethod:  util.GET,
						},
					},
				},
				"ESPv2.HealthCheck": &methodInfo{
					ShortName:          "HealthCheck",
					ApiName:            "ESPv2",
					SkipServiceControl: true,
					IsGenerated:        true,
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/",
							HttpMethod:  util.GET,
						},
					},
				},
			},
		},
		{
			desc: "Succeed for HTTP with multiple apis",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "Echo",
							},
							{
								Name: "Echo_Auth_Jwt",
							},
						},
					},
					{
						Name: "2.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "Echo",
							},
							{
								Name: "Echo_Auth_Jwt",
							},
						},
					},
				},
				Http: &annotationspb.Http{
					Rules: []*annotationspb.HttpRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt",
							Pattern: &annotationspb.HttpRule_Get{
								Get: "/auth/info/googlejwt",
							},
						},
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/echo",
							},
							Body: "message",
						},
						{
							Selector: "2.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt",
							Pattern: &annotationspb.HttpRule_Get{
								Get: "/auth/info/googlejwt",
							},
						},
						{
							Selector: "2.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/echo",
							},
							Body: "message",
						},
					},
				},
			},
			backendProtocol: "http2",
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo": &methodInfo{
					ShortName: "Echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/echo",
							HttpMethod:  util.POST,
						},
					},
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt": &methodInfo{
					ShortName: "Echo_Auth_Jwt",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/auth/info/googlejwt",
							HttpMethod:  util.GET,
						},
					},
				},
				"2.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo": &methodInfo{
					ShortName: "Echo",
					ApiName:   "2.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/echo",
							HttpMethod:  util.POST,
						},
					},
				},
				"2.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt": &methodInfo{
					ShortName: "Echo_Auth_Jwt",
					ApiName:   "2.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/auth/info/googlejwt",
							HttpMethod:  util.GET,
						},
					},
				},
			},
		},
		{
			desc: "Fail for HTTP",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "Echo",
							},
							{
								Name: "Echo_Auth_Jwt",
							},
						},
					},
				},
				Http: &annotationspb.Http{},
			},
			backendProtocol: "http2",
			wantError:       fmt.Sprintf("no HttpRules generated for the Http service %v", testProjectName),
		},
		{
			desc: "Succeed for HTTP, with OPTIONS, and AllowCors, with Healthz",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
						Methods: []*apipb.Method{
							{
								Name: "Echo",
							},
							{
								Name: "Echo_Auth",
							},
							{
								Name: "Echo_Auth_Jwt",
							},
							{
								Name: "EchoCors",
							},
						},
					},
				},
				Endpoints: []*confpb.Endpoint{
					{
						Name:      testProjectName,
						AllowCors: true,
					},
				},
				Http: &annotationspb.Http{
					Rules: []*annotationspb.HttpRule{
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.EchoCors",
							Pattern: &annotationspb.HttpRule_Custom{
								Custom: &annotationspb.CustomHttpPattern{
									Kind: "OPTIONS",
									Path: "/echo",
								},
							},
						},
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/echo",
							},
							Body: "message",
						},
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt",
							Pattern: &annotationspb.HttpRule_Get{
								Get: "/auth/info/googlejwt",
							},
						},
						{
							Selector: "1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/auth/info/googlejwt",
							},
						},
					},
				},
			},
			backendProtocol: "http1",
			healthz:         "/healthz",
			wantMethods: map[string]*methodInfo{
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.EchoCors": &methodInfo{
					ShortName: "EchoCors",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/echo",
							HttpMethod:  util.OPTIONS,
						},
					},
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo": &methodInfo{
					ShortName: "Echo",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/echo",
							HttpMethod:  util.POST,
						},
					},
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.CORS_auth_info_googlejwt": &methodInfo{
					ShortName: "CORS_auth_info_googlejwt",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/auth/info/googlejwt",
							HttpMethod:  util.OPTIONS,
						},
					},
					IsGenerated: true,
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth_Jwt": &methodInfo{
					ShortName: "Echo_Auth_Jwt",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{{
						UriTemplate: "/auth/info/googlejwt",
						HttpMethod:  util.GET,
					},
					},
				},
				"1.echo_api_endpoints_cloudesf_testing_cloud_goog.Echo_Auth": &methodInfo{
					ShortName: "Echo_Auth",
					ApiName:   "1.echo_api_endpoints_cloudesf_testing_cloud_goog",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/auth/info/googlejwt",
							HttpMethod:  util.POST,
						},
					},
				},
				"ESPv2.HealthCheck": &methodInfo{
					ShortName:          "HealthCheck",
					ApiName:            "ESPv2",
					SkipServiceControl: true,
					IsGenerated:        true,
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/healthz",
							HttpMethod:  util.GET,
						},
					},
				},
			},
		},
		{
			desc: "Succeed for multiple url Pattern",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "endpoints.examples.bookstore.Bookstore",
						Methods: []*apipb.Method{
							{
								Name:            "CreateBook",
								RequestTypeUrl:  "type.googleapis.com/endpoints.examples.bookstore.CreateBookRequest",
								ResponseTypeUrl: "type.googleapis.com/endpoints.examples.bookstore.Book",
							},
						},
					},
				},
				Http: &annotationspb.Http{
					Rules: []*annotationspb.HttpRule{
						{
							Selector: "endpoints.examples.bookstore.Bookstore.CreateBook",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/v1/shelves/{shelf}/books/{book.id}/{book.author}",
							},
							Body: "book.title",
						},
						{
							Selector: "endpoints.examples.bookstore.Bookstore.CreateBook",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/v1/shelves/{shelf}/books",
							},
							Body: "book",
						},
					},
				},
			},
			backendProtocol: "grpc",
			wantMethods: map[string]*methodInfo{
				"endpoints.examples.bookstore.Bookstore.CreateBook": &methodInfo{
					ShortName: "CreateBook",
					ApiName:   "endpoints.examples.bookstore.Bookstore",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/v1/shelves/{shelf}/books/{book.id}/{book.author}",
							HttpMethod:  util.POST,
						},
						{
							UriTemplate: "/v1/shelves/{shelf}/books",
							HttpMethod:  util.POST,
						},
					},
				},
			},
		},
		{
			desc: "Succeed for additional binding",
			fakeServiceConfig: &confpb.Service{
				Name: testProjectName,
				Apis: []*apipb.Api{
					{
						Name: "endpoints.examples.bookstore.Bookstore",
						Methods: []*apipb.Method{
							{
								Name:            "CreateBook",
								RequestTypeUrl:  "type.googleapis.com/endpoints.examples.bookstore.CreateBookRequest",
								ResponseTypeUrl: "type.googleapis.com/endpoints.examples.bookstore.Book",
							},
						},
					},
				},
				Http: &annotationspb.Http{
					Rules: []*annotationspb.HttpRule{
						{
							Selector: "endpoints.examples.bookstore.Bookstore.CreateBook",
							Pattern: &annotationspb.HttpRule_Post{
								Post: "/v1/shelves/{shelf}/books/{book.id}/{book.author}",
							},
							Body: "book.title",
							AdditionalBindings: []*annotationspb.HttpRule{
								{
									Pattern: &annotationspb.HttpRule_Post{
										Post: "/v1/shelves/{shelf}/books/foo",
									},
									Body: "book",
								},
								{
									Pattern: &annotationspb.HttpRule_Post{
										Post: "/v1/shelves/{shelf}/books/bar",
									},
									Body: "book",
								},
							},
						},
					},
				},
			},
			backendProtocol: "grpc",
			wantMethods: map[string]*methodInfo{
				"endpoints.examples.bookstore.Bookstore.CreateBook": &methodInfo{
					ShortName: "CreateBook",
					ApiName:   "endpoints.examples.bookstore.Bookstore",
					HttpRule: []*commonpb.Pattern{
						{
							UriTemplate: "/v1/shelves/{shelf}/books/{book.id}/{book.author}",
							HttpMethod:  util.POST,
						},
						{
							UriTemplate: "/v1/shelves/{shelf}/books/foo",
							HttpMethod:  util.POST,
						},
						{
							UriTemplate: "/v1/shelves/{shelf}/books/bar",
							HttpMethod:  util.POST,
						},
					},
				},
			},
		},
	}

	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = tc.backendProtocol
		opts.Healthz = tc.healthz
		serviceInfo, err := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)
		if tc.wantError != "" {
			if err == nil || err.Error() != tc.wantError {
				t.Errorf("Test Desc(%d): %s, got Errors : %v, want: %v", i, tc.desc, err, tc.wantError)
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if len(serviceInfo.Methods) != len(tc.wantMethods) {
			t.Errorf("Test Desc(%d): %s, got Methods: %v, want: %v", i, tc.desc, serviceInfo.Methods, tc.wantMethods)
		}
		for key, gotMethod := range serviceInfo.Methods {
			wantMethod := tc.wantMethods[key]
			if eq := cmp.Equal(gotMethod, wantMethod, cmp.Comparer(proto.Equal)); !eq {
				t.Errorf("Test Desc(%d): %s, got Method: %v, want: %v", i, tc.desc, gotMethod, wantMethod)
			}
		}
	}
}

func TestProcessBackendRule(t *testing.T) {
	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		wantedAllowCors   bool
		wantedErr         string
	}{
		{
			desc: "Failed for dynamic routing only supports HTTPS",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Backend: &confpb.Backend{
					Rules: []*confpb.BackendRule{
						{
							Address:         "http://192.168.0.1/api/",
							PathTranslation: confpb.BackendRule_CONSTANT_ADDRESS,
						},
					},
				},
			},
			wantedErr: "Failed for dynamic routing only supports HTTPS",
		},
		{
			desc: "Fail, dynamic routing only supports domain name, got IP address: 192.168.0.1",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Backend: &confpb.Backend{
					Rules: []*confpb.BackendRule{
						{
							Address:         "https://192.168.0.1/api/",
							PathTranslation: confpb.BackendRule_CONSTANT_ADDRESS,
						},
					},
				},
			},
			wantedErr: "dynamic routing only supports domain name, got IP address: 192.168.0.1",
		},
	}

	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = "grpc"
		_, err := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)
		if (err == nil && tc.wantedErr != "") || (err != nil && tc.wantedErr == "") {
			t.Errorf("Test Desc(%d): %s, extract backend address got: %v, want: %v", i, tc.desc, err, tc.wantedErr)
		}
	}
}

func TestProcessQuota(t *testing.T) {
	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		wantMethods       map[string]*methodInfo
	}{
		{
			desc: "Succeed, simple case",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
						Methods: []*apipb.Method{
							{
								Name: "ListShelves",
							},
						},
					},
				},
				Quota: &confpb.Quota{
					MetricRules: []*confpb.MetricRule{
						{
							Selector: "endpoints.examples.bookstore.Bookstore.ListShelves",
							MetricCosts: map[string]int64{
								"metric_a": 2,
								"metric_b": 1,
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				fmt.Sprintf("%s.%s", testApiName, "ListShelves"): &methodInfo{
					ShortName: "ListShelves",
					ApiName:   testApiName,
					MetricCosts: []*scpb.MetricCost{
						{
							Name: "metric_a",
							Cost: 2,
						},
						{
							Name: "metric_b",
							Cost: 1,
						},
					},
				},
			},
		},
		{
			desc: "Succeed, two metric cost items",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
						Methods: []*apipb.Method{
							{
								Name: "ListShelves",
							},
						},
					},
				},
				Quota: &confpb.Quota{
					MetricRules: []*confpb.MetricRule{
						{
							Selector: "endpoints.examples.bookstore.Bookstore.ListShelves",
							MetricCosts: map[string]int64{
								"metric_c": 2,
								"metric_a": 3,
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				fmt.Sprintf("%s.%s", testApiName, "ListShelves"): &methodInfo{
					ShortName: "ListShelves",
					ApiName:   testApiName,
					MetricCosts: []*scpb.MetricCost{
						{
							Name: "metric_a",
							Cost: 3,
						},
						{
							Name: "metric_c",
							Cost: 2,
						},
					},
				},
			},
		},
	}

	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = "grpc"
		serviceInfo, _ := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)

		for key, gotMethod := range serviceInfo.Methods {
			wantMethod := tc.wantMethods[key]

			sort.Slice(gotMethod.MetricCosts, func(i, j int) bool { return gotMethod.MetricCosts[i].Name < gotMethod.MetricCosts[j].Name })
			if eq := cmp.Equal(gotMethod, wantMethod, cmp.Comparer(proto.Equal)); !eq {
				t.Errorf("Test Desc(%d): %s,\ngot Method: %v,\nwant Method: %v", i, tc.desc, gotMethod, wantMethod)
			}
		}
	}
}

func TestProcessEmptyJwksUriByOpenID(t *testing.T) {
	r := mux.NewRouter()
	jwksUriEntry, _ := json.Marshal(map[string]string{"jwks_uri": "this-is-jwksUri"})
	r.Path(util.OpenIDDiscoveryCfgURLSuffix).Methods("GET").Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(jwksUriEntry)
		}))
	openIDServer := httptest.NewServer(r)

	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		wantedJwksUri     string
		wantErr           bool
	}{
		{
			desc: "Empty jwksUri, use jwksUri acquired by openID",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Authentication: &confpb.Authentication{
					Providers: []*confpb.AuthProvider{
						{
							Id:     "auth_provider",
							Issuer: openIDServer.URL,
						},
					},
				},
			},
			wantedJwksUri: "this-is-jwksUri",
		},
		{
			desc: "Empty jwksUri and Open ID Connect Discovery failed",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: testApiName,
					},
				},
				Authentication: &confpb.Authentication{
					Providers: []*confpb.AuthProvider{
						{
							Id:     "auth_provider",
							Issuer: "aaaaa.bbbbbb.ccccc/inaccessible_uri/",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for i, tc := range testData {
		opts := options.DefaultConfigGeneratorOptions()
		opts.BackendProtocol = "grpc"
		serviceInfo, err := NewServiceInfoFromServiceConfig(tc.fakeServiceConfig, testConfigID, opts)

		if tc.wantErr {
			if err == nil {
				t.Errorf("Test Desc(%d): %s, process jwksUri got: no err, but expected err", i, tc.desc)
			}
		} else if err != nil {
			t.Errorf("Test Desc(%d): %s, process jwksUri got: %v, but expected no err", i, tc.desc, err)
		} else if jwksUri := serviceInfo.serviceConfig.Authentication.Providers[0].JwksUri; jwksUri != tc.wantedJwksUri {
			t.Errorf("Test Desc(%d): %s, process jwksUri got: %v, want: %v", i, tc.desc, jwksUri, tc.wantedJwksUri)
		}
	}
}

func TestProcessApis(t *testing.T) {
	testData := []struct {
		desc              string
		fakeServiceConfig *confpb.Service
		wantMethods       map[string]*methodInfo
		wantApiNames      []string
	}{
		{
			desc: "Succeed, process multiple apis",
			fakeServiceConfig: &confpb.Service{
				Apis: []*apipb.Api{
					{
						Name: "api-1",
						Methods: []*apipb.Method{
							{
								Name: "foo",
							},
							{
								Name: "bar",
							},
						},
					},
					{
						Name: "api-2",
						Methods: []*apipb.Method{
							{
								Name: "foo",
							},
							{
								Name: "bar",
							},
						},
					},
					{
						Name:    "api-3",
						Methods: []*apipb.Method{},
					},
					{
						Name: "api-4",
						Methods: []*apipb.Method{
							{
								Name: "bar",
							},
							{
								Name: "baz",
							},
						},
					},
				},
			},
			wantMethods: map[string]*methodInfo{
				"api-1.foo": &methodInfo{
					ShortName: "foo",
					ApiName:   "api-1",
				},
				"api-1.bar": &methodInfo{
					ShortName: "bar",
					ApiName:   "api-1",
				},
				"api-2.foo": &methodInfo{
					ShortName: "foo",
					ApiName:   "api-2",
				},
				"api-2.bar": &methodInfo{
					ShortName: "bar",
					ApiName:   "api-2",
				},
				"api-4.bar": &methodInfo{
					ShortName: "bar",
					ApiName:   "api-4",
				},
				"api-4.baz": &methodInfo{
					ShortName: "baz",
					ApiName:   "api-4",
				},
			},
			wantApiNames: []string{
				"api-1",
				"api-2",
				"api-3",
				"api-4",
			},
		},
	}

	for i, tc := range testData {

		serviceInfo := &ServiceInfo{
			serviceConfig: tc.fakeServiceConfig,
		}
		serviceInfo.processApis()

		for key, gotMethod := range serviceInfo.Methods {
			wantMethod := tc.wantMethods[key]
			if eq := cmp.Equal(gotMethod, wantMethod, cmp.Comparer(proto.Equal)); !eq {
				t.Errorf("Test Desc(%d): %s,\ngot Method: %v,\nwant Method: %v", i, tc.desc, gotMethod, wantMethod)
			}
		}
		for idx, gotApiName := range serviceInfo.ApiNames {
			wantApiName := tc.wantApiNames[idx]
			if gotApiName != wantApiName {
				t.Errorf("Test Desc(%d): %s,\ngot ApiName: %v,\nwant Apiname: %v", i, tc.desc, gotApiName, wantApiName)
			}
		}
	}
}
