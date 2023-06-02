// Copyright 2023 Google LLC. All Rights Reserved.
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
// GENERATED BY gen_go_data.go
// gen_go_data -package logging -var YAML_log_metric blaze-out/k8-fastbuild/genfiles/cloud/graphite/mmv2/services/google/logging/log_metric.yaml

package logging

// blaze-out/k8-fastbuild/genfiles/cloud/graphite/mmv2/services/google/logging/log_metric.yaml
var YAML_log_metric = []byte("info:\n  title: Logging/LogMetric\n  description: The Logging LogMetric resource\n  x-dcl-struct-name: LogMetric\n  x-dcl-has-iam: false\npaths:\n  get:\n    description: The function used to get information about a LogMetric\n    parameters:\n    - name: logMetric\n      required: true\n      description: A full instance of a LogMetric\n  apply:\n    description: The function used to apply information about a LogMetric\n    parameters:\n    - name: logMetric\n      required: true\n      description: A full instance of a LogMetric\n  delete:\n    description: The function used to delete a LogMetric\n    parameters:\n    - name: logMetric\n      required: true\n      description: A full instance of a LogMetric\n  deleteAll:\n    description: The function used to delete all LogMetric\n    parameters:\n    - name: project\n      required: true\n      schema:\n        type: string\n  list:\n    description: The function used to list information about many LogMetric\n    parameters:\n    - name: project\n      required: true\n      schema:\n        type: string\ncomponents:\n  schemas:\n    LogMetric:\n      title: LogMetric\n      x-dcl-id: projects/{{project}}/metrics/{{name}}\n      x-dcl-uses-state-hint: true\n      x-dcl-parent-container: project\n      x-dcl-has-create: true\n      x-dcl-has-iam: false\n      x-dcl-read-timeout: 0\n      x-dcl-apply-timeout: 0\n      x-dcl-delete-timeout: 0\n      type: object\n      required:\n      - name\n      - filter\n      - project\n      properties:\n        bucketOptions:\n          type: object\n          x-dcl-go-name: BucketOptions\n          x-dcl-go-type: LogMetricBucketOptions\n          description: Optional. The `bucket_options` are required when the logs-based\n            metric is using a DISTRIBUTION value type and it describes the bucket\n            boundaries used to create a histogram of the extracted values.\n          properties:\n            explicitBuckets:\n              type: object\n              x-dcl-go-name: ExplicitBuckets\n              x-dcl-go-type: LogMetricBucketOptionsExplicitBuckets\n              description: The explicit buckets.\n              x-dcl-conflicts:\n              - linearBuckets\n              - exponentialBuckets\n              properties:\n                bounds:\n                  type: array\n                  x-dcl-go-name: Bounds\n                  description: The values must be monotonically increasing.\n                  x-dcl-send-empty: true\n                  x-dcl-list-type: list\n                  items:\n                    type: number\n                    format: double\n                    x-dcl-go-type: float64\n            exponentialBuckets:\n              type: object\n              x-dcl-go-name: ExponentialBuckets\n              x-dcl-go-type: LogMetricBucketOptionsExponentialBuckets\n              description: The exponential buckets.\n              x-dcl-conflicts:\n              - linearBuckets\n              - explicitBuckets\n              properties:\n                growthFactor:\n                  type: number\n                  format: double\n                  x-dcl-go-name: GrowthFactor\n                  description: Must be greater than 1.\n                numFiniteBuckets:\n                  type: integer\n                  format: int64\n                  x-dcl-go-name: NumFiniteBuckets\n                  description: Must be greater than 0.\n                scale:\n                  type: number\n                  format: double\n                  x-dcl-go-name: Scale\n                  description: Must be greater than 0.\n            linearBuckets:\n              type: object\n              x-dcl-go-name: LinearBuckets\n              x-dcl-go-type: LogMetricBucketOptionsLinearBuckets\n              description: The linear bucket.\n              x-dcl-conflicts:\n              - exponentialBuckets\n              - explicitBuckets\n              properties:\n                numFiniteBuckets:\n                  type: integer\n                  format: int64\n                  x-dcl-go-name: NumFiniteBuckets\n                  description: Must be greater than 0.\n                offset:\n                  type: number\n                  format: double\n                  x-dcl-go-name: Offset\n                  description: Lower bound of the first bucket.\n                width:\n                  type: number\n                  format: double\n                  x-dcl-go-name: Width\n                  description: Must be greater than 0.\n        createTime:\n          type: string\n          format: date-time\n          x-dcl-go-name: CreateTime\n          readOnly: true\n          description: Output only. The creation timestamp of the metric. This field\n            may not be present for older metrics.\n          x-kubernetes-immutable: true\n        description:\n          type: string\n          x-dcl-go-name: Description\n          description: Optional. A description of this metric, which is used in documentation.\n            The maximum length of the description is 8000 characters.\n        disabled:\n          type: boolean\n          x-dcl-go-name: Disabled\n          description: Optional. If set to True, then this metric is disabled and\n            it does not generate any points.\n        filter:\n          type: string\n          x-dcl-go-name: Filter\n          description: 'Required. An [advanced logs filter](https://cloud.google.com/logging/docs/view/advanced_filters)\n            which is used to match log entries. Example: \"resource.type=gae_app AND\n            severity>=ERROR\" The maximum length of the filter is 20000 characters.'\n        labelExtractors:\n          type: object\n          additionalProperties:\n            type: string\n          x-dcl-go-name: LabelExtractors\n          description: Optional. A map from a label key string to an extractor expression\n            which is used to extract data from a log entry field and assign as the\n            label value. Each label key specified in the LabelDescriptor must have\n            an associated extractor expression in this map. The syntax of the extractor\n            expression is the same as for the `value_extractor` field. The extracted\n            value is converted to the type defined in the label descriptor. If the\n            either the extraction or the type conversion fails, the label will have\n            a default value. The default value for a string label is an empty string,\n            for an integer label its 0, and for a boolean label its `false`. Note\n            that there are upper bounds on the maximum number of labels and the number\n            of active time series that are allowed in a project.\n        metricDescriptor:\n          type: object\n          x-dcl-go-name: MetricDescriptor\n          x-dcl-go-type: LogMetricMetricDescriptor\n          description: Optional. The metric descriptor associated with the logs-based\n            metric. If unspecified, it uses a default metric descriptor with a DELTA\n            metric kind, INT64 value type, with no labels and a unit of \"1\". Such\n            a metric counts the number of log entries matching the `filter` expression.\n            The `name`, `type`, and `description` fields in the `metric_descriptor`\n            are output only, and is constructed using the `name` and `description`\n            field in the LogMetric. To create a logs-based metric that records a distribution\n            of log values, a DELTA metric kind with a DISTRIBUTION value type must\n            be used along with a `value_extractor` expression in the LogMetric. Each\n            label in the metric descriptor must have a matching label name as the\n            key and an extractor expression as the value in the `label_extractors`\n            map. The `metric_kind` and `value_type` fields in the `metric_descriptor`\n            cannot be updated once initially configured. New labels can be added in\n            the `metric_descriptor`, but existing labels cannot be modified except\n            for their description.\n          properties:\n            description:\n              type: string\n              x-dcl-go-name: Description\n              readOnly: true\n              description: A detailed description of the metric, which can be used\n                in documentation.\n            displayName:\n              type: string\n              x-dcl-go-name: DisplayName\n              description: A concise name for the metric, which can be displayed in\n                user interfaces. Use sentence case without an ending period, for example\n                \"Request count\". This field is optional but it is recommended to be\n                set for any metrics associated with user-visible concepts, such as\n                Quota.\n            labels:\n              type: array\n              x-dcl-go-name: Labels\n              description: The set of labels that can be used to describe a specific\n                instance of this metric type. For example, the `appengine.googleapis.com/http/server/response_latencies`\n                metric type has a label for the HTTP response code, `response_code`,\n                so you can look at latencies for successful responses or just for\n                responses that failed.\n              x-dcl-send-empty: true\n              x-dcl-list-type: set\n              items:\n                type: object\n                x-dcl-go-type: LogMetricMetricDescriptorLabels\n                properties:\n                  description:\n                    type: string\n                    x-dcl-go-name: Description\n                    description: A human-readable description for the label.\n                    x-kubernetes-immutable: true\n                  key:\n                    type: string\n                    x-dcl-go-name: Key\n                    description: The label key.\n                    x-kubernetes-immutable: true\n                  valueType:\n                    type: string\n                    x-dcl-go-name: ValueType\n                    x-dcl-go-type: LogMetricMetricDescriptorLabelsValueTypeEnum\n                    description: 'The type of data that can be assigned to the label.\n                      Possible values: STRING, BOOL, INT64, DOUBLE, DISTRIBUTION,\n                      MONEY'\n                    x-kubernetes-immutable: true\n                    enum:\n                    - STRING\n                    - BOOL\n                    - INT64\n                    - DOUBLE\n                    - DISTRIBUTION\n                    - MONEY\n            launchStage:\n              type: string\n              x-dcl-go-name: LaunchStage\n              x-dcl-go-type: LogMetricMetricDescriptorLaunchStageEnum\n              description: 'Optional. The launch stage of the metric definition. Possible\n                values: UNIMPLEMENTED, PRELAUNCH, EARLY_ACCESS, ALPHA, BETA, GA, DEPRECATED'\n              enum:\n              - UNIMPLEMENTED\n              - PRELAUNCH\n              - EARLY_ACCESS\n              - ALPHA\n              - BETA\n              - GA\n              - DEPRECATED\n              x-dcl-mutable-unreadable: true\n            metadata:\n              type: object\n              x-dcl-go-name: Metadata\n              x-dcl-go-type: LogMetricMetricDescriptorMetadata\n              description: Optional. Metadata which can be used to guide usage of\n                the metric.\n              x-dcl-mutable-unreadable: true\n              properties:\n                ingestDelay:\n                  type: string\n                  x-dcl-go-name: IngestDelay\n                  description: The delay of data points caused by ingestion. Data\n                    points older than this age are guaranteed to be ingested and available\n                    to be read, excluding data loss due to errors.\n                samplePeriod:\n                  type: string\n                  x-dcl-go-name: SamplePeriod\n                  description: The sampling period of metric data points. For metrics\n                    which are written periodically, consecutive data points are stored\n                    at this time interval, excluding data loss due to errors. Metrics\n                    with a higher granularity have a smaller sampling period.\n            metricKind:\n              type: string\n              x-dcl-go-name: MetricKind\n              x-dcl-go-type: LogMetricMetricDescriptorMetricKindEnum\n              description: 'Whether the metric records instantaneous values, changes\n                to a value, etc. Some combinations of `metric_kind` and `value_type`\n                might not be supported. Possible values: GAUGE, DELTA, CUMULATIVE'\n              x-kubernetes-immutable: true\n              enum:\n              - GAUGE\n              - DELTA\n              - CUMULATIVE\n            monitoredResourceTypes:\n              type: array\n              x-dcl-go-name: MonitoredResourceTypes\n              readOnly: true\n              description: Read-only. If present, then a time series, which is identified\n                partially by a metric type and a MonitoredResourceDescriptor, that\n                is associated with this metric type can only be associated with one\n                of the monitored resource types listed here.\n              x-kubernetes-immutable: true\n              x-dcl-list-type: list\n              items:\n                type: string\n                x-dcl-go-type: string\n            name:\n              type: string\n              x-dcl-go-name: Name\n              readOnly: true\n              description: The resource name of the metric descriptor.\n              x-kubernetes-immutable: true\n            type:\n              type: string\n              x-dcl-go-name: Type\n              readOnly: true\n              description: 'The metric type, including its DNS name prefix. The type\n                is not URL-encoded. All user-defined metric types have the DNS name\n                `custom.googleapis.com` or `external.googleapis.com`. Metric types\n                should use a natural hierarchical grouping. For example: \"custom.googleapis.com/invoice/paid/amount\"\n                \"external.googleapis.com/prometheus/up\" \"appengine.googleapis.com/http/server/response_latencies\"'\n              x-kubernetes-immutable: true\n            unit:\n              type: string\n              x-dcl-go-name: Unit\n              description: 'The units in which the metric value is reported. It is\n                only applicable if the `value_type` is `INT64`, `DOUBLE`, or `DISTRIBUTION`.\n                The `unit` defines the representation of the stored metric values.\n                Different systems might scale the values to be more easily displayed\n                (so a value of `0.02kBy` _might_ be displayed as `20By`, and a value\n                of `3523kBy` _might_ be displayed as `3.5MBy`). However, if the `unit`\n                is `kBy`, then the value of the metric is always in thousands of bytes,\n                no matter how it might be displayed. If you want a custom metric to\n                record the exact number of CPU-seconds used by a job, you can create\n                an `INT64 CUMULATIVE` metric whose `unit` is `s{CPU}` (or equivalently\n                `1s{CPU}` or just `s`). If the job uses 12,005 CPU-seconds, then the\n                value is written as `12005`. Alternatively, if you want a custom metric\n                to record data in a more granular way, you can create a `DOUBLE CUMULATIVE`\n                metric whose `unit` is `ks{CPU}`, and then write the value `12.005`\n                (which is `12005/1000`), or use `Kis{CPU}` and write `11.723` (which\n                is `12005/1024`). The supported units are a subset of [The Unified\n                Code for Units of Measure](https://unitsofmeasure.org/ucum.html) standard:\n                **Basic units (UNIT)** * `bit` bit * `By` byte * `s` second * `min`\n                minute * `h` hour * `d` day * `1` dimensionless **Prefixes (PREFIX)**\n                * `k` kilo (10^3) * `M` mega (10^6) * `G` giga (10^9) * `T` tera (10^12)\n                * `P` peta (10^15) * `E` exa (10^18) * `Z` zetta (10^21) * `Y` yotta\n                (10^24) * `m` milli (10^-3) * `u` micro (10^-6) * `n` nano (10^-9)\n                * `p` pico (10^-12) * `f` femto (10^-15) * `a` atto (10^-18) * `z`\n                zepto (10^-21) * `y` yocto (10^-24) * `Ki` kibi (2^10) * `Mi` mebi\n                (2^20) * `Gi` gibi (2^30) * `Ti` tebi (2^40) * `Pi` pebi (2^50) **Grammar**\n                The grammar also includes these connectors: * `/` division or ratio\n                (as an infix operator). For examples, `kBy/{email}` or `MiBy/10ms`\n                (although you should almost never have `/s` in a metric `unit`; rates\n                should always be computed at query time from the underlying cumulative\n                or delta value). * `.` multiplication or composition (as an infix\n                operator). For examples, `GBy.d` or `k{watt}.h`. The grammar for a\n                unit is as follows: Expression = Component: { \".\" Component } { \"/\"\n                Component } ; Component = ( [ PREFIX ] UNIT | \"%\" ) [ Annotation ]\n                | Annotation | \"1\" ; Annotation = \"{\" NAME \"}\" ; Notes: * `Annotation`\n                is just a comment if it follows a `UNIT`. If the annotation is used\n                alone, then the unit is equivalent to `1`. For examples, `{request}/s\n                == 1/s`, `By{transmitted}/s == By/s`. * `NAME` is a sequence of non-blank\n                printable ASCII characters not containing `{` or `}`. * `1` represents\n                a unitary [dimensionless unit](https://en.wikipedia.org/wiki/Dimensionless_quantity)\n                of 1, such as in `1/s`. It is typically used when none of the basic\n                units are appropriate. For example, \"new users per day\" can be represented\n                as `1/d` or `{new-users}/d` (and a metric value `5` would mean \"5\n                new users). Alternatively, \"thousands of page views per day\" would\n                be represented as `1000/d` or `k1/d` or `k{page_views}/d` (and a metric\n                value of `5.3` would mean \"5300 page views per day\"). * `%` represents\n                dimensionless value of 1/100, and annotates values giving a percentage\n                (so the metric values are typically in the range of 0..100, and a\n                metric value `3` means \"3 percent\"). * `10^2.%` indicates a metric\n                contains a ratio, typically in the range 0..1, that will be multiplied\n                by 100 and displayed as a percentage (so a metric value `0.03` means\n                \"3 percent\").'\n              x-dcl-server-default: true\n            valueType:\n              type: string\n              x-dcl-go-name: ValueType\n              x-dcl-go-type: LogMetricMetricDescriptorValueTypeEnum\n              description: 'Whether the measurement is an integer, a floating-point\n                number, etc. Some combinations of `metric_kind` and `value_type` might\n                not be supported. Possible values: STRING, BOOL, INT64, DOUBLE, DISTRIBUTION,\n                MONEY'\n              x-kubernetes-immutable: true\n              enum:\n              - STRING\n              - BOOL\n              - INT64\n              - DOUBLE\n              - DISTRIBUTION\n              - MONEY\n        name:\n          type: string\n          x-dcl-go-name: Name\n          description: 'Required. The client-assigned metric identifier. Examples:\n            `\"error_count\"`, `\"nginx/requests\"`. Metric identifiers are limited to\n            100 characters and can include only the following characters: `A-Z`, `a-z`,\n            `0-9`, and the special characters `_-.,+!*'',()%/`. The forward-slash\n            character (`/`) denotes a hierarchy of name pieces, and it cannot be the\n            first character of the name. The metric identifier in this field must\n            not be [URL-encoded](https://en.wikipedia.org/wiki/Percent-encoding).\n            However, when the metric identifier appears as the `[METRIC_ID]` part\n            of a `metric_name` API parameter, then the metric identifier must be URL-encoded.\n            Example: `\"projects/my-project/metrics/nginx%2Frequests\"`.'\n          x-kubernetes-immutable: true\n        project:\n          type: string\n          x-dcl-go-name: Project\n          description: The resource name of the project in which to create the metric.\n          x-kubernetes-immutable: true\n          x-dcl-references:\n          - resource: Cloudresourcemanager/Project\n            field: name\n            parent: true\n        updateTime:\n          type: string\n          format: date-time\n          x-dcl-go-name: UpdateTime\n          readOnly: true\n          description: Output only. The last update timestamp of the metric. This\n            field may not be present for older metrics.\n          x-kubernetes-immutable: true\n        valueExtractor:\n          type: string\n          x-dcl-go-name: ValueExtractor\n          description: 'Optional. A `value_extractor` is required when using a distribution\n            logs-based metric to extract the values to record from a log entry. Two\n            functions are supported for value extraction: `EXTRACT(field)` or `REGEXP_EXTRACT(field,\n            regex)`. The argument are: 1. field: The name of the log entry field from\n            which the value is to be extracted. 2. regex: A regular expression using\n            the Google RE2 syntax (https://github.com/google/re2/wiki/Syntax) with\n            a single capture group to extract data from the specified log entry field.\n            The value of the field is converted to a string before applying the regex.\n            It is an error to specify a regex that does not include exactly one capture\n            group. The result of the extraction must be convertible to a double type,\n            as the distribution always records double values. If either the extraction\n            or the conversion to double fails, then those values are not recorded\n            in the distribution. Example: `REGEXP_EXTRACT(jsonPayload.request, \".*quantity=(d+).*\")`'\n")

// 22100 bytes
// MD5: b55014c3141d88528a1a93f25cd97b16
