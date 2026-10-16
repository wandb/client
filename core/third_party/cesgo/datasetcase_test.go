// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package coreweaveevaluations_test

import (
	"context"
	"errors"
	"os"
	"testing"

	coreweaveevaluations "github.com/wandb/core/services/evaluations/generated/go"
	"github.com/wandb/core/services/evaluations/generated/go/internal/testutil"
	"github.com/wandb/core/services/evaluations/generated/go/option"
)

func TestDatasetCaseAdd(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := coreweaveevaluations.NewClient(
		option.WithBaseURL(baseURL),
		option.WithBearerToken("My Bearer Token"),
		option.WithUsername("My Username"),
		option.WithAPIKey("My API Key"),
	)
	_, err := client.Datasets.Cases.Add(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		coreweaveevaluations.DatasetCaseAddParams{
			Namespace: coreweaveevaluations.F("wandb"),
			ScopeID:   coreweaveevaluations.F("scope_id"),
			Cases: coreweaveevaluations.F([]coreweaveevaluations.DatasetCaseAddParamsCase{
				{
					Input: coreweaveevaluations.F[coreweaveevaluations.CaseValueUnionParam](
						coreweaveevaluations.CaseValueMapParam(map[string]coreweaveevaluations.CaseFieldValueUnionParam{
							"foo": coreweaveevaluations.WandbImageV1Param{
								ExtensionType: coreweaveevaluations.F(
									coreweaveevaluations.WandbImageV1ExtensionTypeWandbImage,
								),
								Format:        coreweaveevaluations.F("x"),
								SchemaVersion: coreweaveevaluations.F(coreweaveevaluations.WandbImageV1SchemaVersion1),
								Sha256: coreweaveevaluations.F(
									"210b9798eb53baa4e69d31c1071cf03d212b8ad0ca30cf321e0ea82e120aac26",
								),
								Size: coreweaveevaluations.F(int64(0)),
								Uri:  coreweaveevaluations.F("x"),
								WbMediaType: coreweaveevaluations.F(
									coreweaveevaluations.WandbImageV1WbMediaTypeImageFile,
								),
								Boxes: coreweaveevaluations.F(map[string]interface{}{
									"foo": "bar",
								}),
								Caption: coreweaveevaluations.F("caption"),
								Height:  coreweaveevaluations.F(int64(0)),
								Masks: coreweaveevaluations.F(map[string]interface{}{
									"foo": "bar",
								}),
								Width: coreweaveevaluations.F(int64(0)),
							},
						}),
					),
					Output: coreweaveevaluations.F[coreweaveevaluations.CaseValueUnionParam](
						coreweaveevaluations.CaseValueMapParam(map[string]coreweaveevaluations.CaseFieldValueUnionParam{
							"foo": coreweaveevaluations.WandbImageV1Param{
								ExtensionType: coreweaveevaluations.F(
									coreweaveevaluations.WandbImageV1ExtensionTypeWandbImage,
								),
								Format:        coreweaveevaluations.F("x"),
								SchemaVersion: coreweaveevaluations.F(coreweaveevaluations.WandbImageV1SchemaVersion1),
								Sha256: coreweaveevaluations.F(
									"210b9798eb53baa4e69d31c1071cf03d212b8ad0ca30cf321e0ea82e120aac26",
								),
								Size: coreweaveevaluations.F(int64(0)),
								Uri:  coreweaveevaluations.F("x"),
								WbMediaType: coreweaveevaluations.F(
									coreweaveevaluations.WandbImageV1WbMediaTypeImageFile,
								),
								Boxes: coreweaveevaluations.F(map[string]interface{}{
									"foo": "bar",
								}),
								Caption: coreweaveevaluations.F("caption"),
								Height:  coreweaveevaluations.F(int64(0)),
								Masks: coreweaveevaluations.F(map[string]interface{}{
									"foo": "bar",
								}),
								Width: coreweaveevaluations.F(int64(0)),
							},
						}),
					),
					Attributes: coreweaveevaluations.F(map[string]interface{}{
						"foo": "bar",
					}),
					Subject: coreweaveevaluations.F("x"),
				},
			}),
		},
	)
	if err != nil {
		var apierr *coreweaveevaluations.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
