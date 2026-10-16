# Service

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ServiceGetResponse">ServiceGetResponse</a>

Methods:

- <code title="get /">client.Service.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ServiceService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ServiceGetResponse">ServiceGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Verification

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#VerificationGetResponse">VerificationGetResponse</a>

Methods:

- <code title="get /test">client.Verification.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#VerificationService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#VerificationGetResponse">VerificationGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Datasets

Params Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#CaseFieldValueUnionParam">CaseFieldValueUnionParam</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#CaseValueUnionParam">CaseValueUnionParam</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#WandbAudioV1Param">WandbAudioV1Param</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#WandbImageV1Param">WandbImageV1Param</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#WandbVideoV1Param">WandbVideoV1Param</a>

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetNewResponse">DatasetNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/datasets">client.Datasets.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetNewParams">DatasetNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetNewResponse">DatasetNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Cases

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetCaseAddResponse">DatasetCaseAddResponse</a>

Methods:

- <code title="patch /v1/namespaces/{namespace}/scopes/{scope_id}/datasets/{dataset_id}/cases">client.Datasets.Cases.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetCaseService.Add">Add</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, datasetID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetCaseAddParams">DatasetCaseAddParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetCaseAddResponse">DatasetCaseAddResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Versions

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionNewResponse">DatasetVersionNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionGetResponse">DatasetVersionGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/datasets/{dataset_id}/versions">client.Datasets.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, datasetID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionNewParams">DatasetVersionNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionNewResponse">DatasetVersionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/dataset-versions/{dataset_version_id}">client.Datasets.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, datasetVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionGetParams">DatasetVersionGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionGetResponse">DatasetVersionGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Cases

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionCaseListResponse">DatasetVersionCaseListResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/dataset-versions/{dataset_version_id}/cases">client.Datasets.Versions.Cases.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionCaseService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, datasetVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionCaseListParams">DatasetVersionCaseListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionCaseListResponse">DatasetVersionCaseListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Fields

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionFieldListResponse">DatasetVersionFieldListResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/dataset-versions/{dataset_version_id}/fields">client.Datasets.Versions.Fields.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionFieldService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, datasetVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionFieldListParams">DatasetVersionFieldListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#DatasetVersionFieldListResponse">DatasetVersionFieldListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Comparisons

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonNewResponse">ComparisonNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonGetResponse">ComparisonGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/comparisons">client.Comparisons.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonNewParams">ComparisonNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonNewResponse">ComparisonNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/comparisons/{comparison_id}">client.Comparisons.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, comparisonID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonGetParams">ComparisonGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ComparisonGetResponse">ComparisonGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Scorers

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerNewResponse">ScorerNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/scorers">client.Scorers.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerNewParams">ScorerNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerNewResponse">ScorerNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Versions

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionNewResponse">ScorerVersionNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionGetResponse">ScorerVersionGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/scorers/{scorer_id}/versions">client.Scorers.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, scorerID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionNewParams">ScorerVersionNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionNewResponse">ScorerVersionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/scorers/{scorer_id}/versions/{scorer_version_id}">client.Scorers.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, scorerID <a href="https://pkg.go.dev/builtin#string">string</a>, scorerVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionGetParams">ScorerVersionGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScorerVersionGetResponse">ScorerVersionGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Scores

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreNewResponse">ScoreNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreGetResponse">ScoreGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/scores">client.Scores.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreNewParams">ScoreNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreNewResponse">ScoreNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/scores/{score_id}">client.Scores.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, scoreID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreGetParams">ScoreGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#ScoreGetResponse">ScoreGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Evaluations

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationNewResponse">EvaluationNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations">client.Evaluations.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationNewParams">EvaluationNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationNewResponse">EvaluationNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Queries

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryNewResponse">EvaluationQueryNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryListResponse">EvaluationQueryListResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryGetResponse">EvaluationQueryGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations/{evaluation_id}/queries">client.Evaluations.Queries.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryNewParams">EvaluationQueryNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryNewResponse">EvaluationQueryNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations/{evaluation_id}/queries">client.Evaluations.Queries.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryListParams">EvaluationQueryListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryListResponse">EvaluationQueryListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations/{evaluation_id}/queries/{evaluation_query_id}">client.Evaluations.Queries.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, evaluationQueryID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryGetParams">EvaluationQueryGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryGetResponse">EvaluationQueryGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Results

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryResultListResponse">EvaluationQueryResultListResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations/{evaluation_id}/queries/{evaluation_query_id}/results">client.Evaluations.Queries.Results.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryResultService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, evaluationQueryID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryResultListParams">EvaluationQueryResultListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationQueryResultListResponse">EvaluationQueryResultListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Versions

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionNewResponse">EvaluationVersionNewResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionGetResponse">EvaluationVersionGetResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/evaluations/{evaluation_id}/versions">client.Evaluations.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionNewParams">EvaluationVersionNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionNewResponse">EvaluationVersionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluation-versions/{evaluation_version_id}">client.Evaluations.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionGetParams">EvaluationVersionGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionGetResponse">EvaluationVersionGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Results

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionResultListResponse">EvaluationVersionResultListResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluation-versions/{evaluation_version_id}/results">client.Evaluations.Versions.Results.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionResultService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionResultListParams">EvaluationVersionResultListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionResultListResponse">EvaluationVersionResultListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Queries

#### Results

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionQueryResultListResponse">EvaluationVersionQueryResultListResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluation-versions/{evaluation_version_id}/queries/{evaluation_query_id}/results">client.Evaluations.Versions.Queries.Results.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionQueryResultService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, evaluationQueryID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionQueryResultListParams">EvaluationVersionQueryResultListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionQueryResultListResponse">EvaluationVersionQueryResultListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Rows

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowListResponse">EvaluationVersionRowListResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowSearchResponse">EvaluationVersionRowSearchResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/evaluation-versions/{evaluation_version_id}/rows">client.Evaluations.Versions.Rows.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowListParams">EvaluationVersionRowListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowListResponse">EvaluationVersionRowListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/evaluation-versions/{evaluation_version_id}/rows/search">client.Evaluations.Versions.Rows.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowService.Search">Search</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowSearchParams">EvaluationVersionRowSearchParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvaluationVersionRowSearchResponse">EvaluationVersionRowSearchResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# Trajectories

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#TrajectoryGetResponse">TrajectoryGetResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/trajectories/{trajectory_id}">client.Trajectories.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#TrajectoryService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, trajectoryID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#TrajectoryGetParams">TrajectoryGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#TrajectoryGetResponse">TrajectoryGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

# EvalTables

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableNewResponse">EvalTableNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/eval-tables">client.EvalTables.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableNewParams">EvalTableNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableNewResponse">EvalTableNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Columns

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableColumnNewResponse">EvalTableColumnNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/eval-tables/{evaluation_id}/columns">client.EvalTables.Columns.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableColumnService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableColumnNewParams">EvalTableColumnNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableColumnNewResponse">EvalTableColumnNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Rows

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableRowAddResponse">EvalTableRowAddResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/eval-tables/{evaluation_id}/rows">client.EvalTables.Rows.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableRowService.Add">Add</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableRowAddParams">EvalTableRowAddParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableRowAddResponse">EvalTableRowAddResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

## Versions

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionNewResponse">EvalTableVersionNewResponse</a>

Methods:

- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/eval-tables/{evaluation_id}/versions">client.EvalTables.Versions.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionService.New">New</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationID <a href="https://pkg.go.dev/builtin#string">string</a>, body <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionNewParams">EvalTableVersionNewParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionNewResponse">EvalTableVersionNewResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Columns

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionColumnGetResponse">EvalTableVersionColumnGetResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/eval-table-versions/{evaluation_version_id}/columns">client.EvalTables.Versions.Columns.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionColumnService.Get">Get</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, query <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionColumnGetParams">EvalTableVersionColumnGetParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionColumnGetResponse">EvalTableVersionColumnGetResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>

### Rows

Response Types:

- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowListResponse">EvalTableVersionRowListResponse</a>
- <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowSearchResponse">EvalTableVersionRowSearchResponse</a>

Methods:

- <code title="get /v1/namespaces/{namespace}/scopes/{scope_id}/eval-table-versions/{evaluation_version_id}/rows">client.EvalTables.Versions.Rows.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowService.List">List</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowListParams">EvalTableVersionRowListParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowListResponse">EvalTableVersionRowListResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
- <code title="post /v1/namespaces/{namespace}/scopes/{scope_id}/eval-table-versions/{evaluation_version_id}/rows/search">client.EvalTables.Versions.Rows.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowService.Search">Search</a>(ctx <a href="https://pkg.go.dev/context">context</a>.<a href="https://pkg.go.dev/context#Context">Context</a>, evaluationVersionID <a href="https://pkg.go.dev/builtin#string">string</a>, params <a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowSearchParams">EvalTableVersionRowSearchParams</a>) (\*<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go">coreweaveevaluations</a>.<a href="https://pkg.go.dev/github.com/wandb/core/services/evaluations/generated/go#EvalTableVersionRowSearchResponse">EvalTableVersionRowSearchResponse</a>, <a href="https://pkg.go.dev/builtin#error">error</a>)</code>
