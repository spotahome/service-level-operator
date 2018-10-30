# service-level-operator

Service level operator abstracts and automates the service level of Kubernetes applications by generation SLI & SLOs to be consumed easily by dashboards and alerts and allow that the SLI/SLO's live with the application flow.

This operator interacts with Kubernetes using the CRDs as a way to define application service levels and generating output service level metrics.

Although this operator is though to interact with different backends and generate different output backends, at this moment only uses [Prometheus] as input and output backend.

## Example

For this example the output and input backend will be [Prometheus].

First you will need to define a CRD with your service SLI & SLOs. In this case we have a service that has an SLO on 99.99 availability, and the SLI is that 5xx are considered errors.

```yaml
apiVersion: measure.slok.xyz/v1alpha1
kind: ServiceLevel
metadata:
  name: awesome-service
spec:
  serviceLevelObjectives:
    - name: "9999_http_request_lt_500"
      description: 99.99% of requests must be served with <500 status code.
      disable: false
      availabilityObjectivePercent: 99.99
      serviceLevelIndicator:
        prometheus:
          address: http://myprometheus:9090
          totalQuery: sum(increase(http_request_total{host="awesome_service_io"}[2m]))
          errorQuery: sum(increase(http_request_total{host="awesome_service_io", code=~"5.."}[2m]))
      output:
        prometheus:
          labels:
            team: a-team
            iteration: "3"
```

The Operator will generate the SLI and SLO in this prometheus format:

```text
# HELP service_level_sli_result_count_total Is the number of times an SLI result has been processed.
# TYPE service_level_sli_result_count_total counter
service_level_sli_result_count_total{service_level="awesome-service",slo="9999_http_request_lt_500"} 1708
# HELP service_level_sli_result_error_ratio_total Is the error or failure ratio of an SLI result.
# TYPE service_level_sli_result_error_ratio_total counter
service_level_sli_result_error_ratio_total{service_level="awesome-service",slo="9999_http_request_lt_500"} 0.40508550763795764
# HELP service_level_slo_objective_ratio Is the objective of the SLO in ratio unit.
# TYPE service_level_slo_objective_ratio gauge
service_level_slo_objective_ratio{service_level="awesome-service",slo="9999_http_request_lt_500"} 0.9998999999999999
```

## How does it work

The operator will query and create new metrics based on the SLOs caulculations at regular intervals (see `--resync-seconds` flag).

The approach that has been taken to generate the SLI results is based on [how Google uses and manages SLIs, SLOs and error budgets][sre-book-slo]

In the manifest the SLI is made of 2 prometheus metrics:

- The total of requests: `sum(increase(http_request_total{host="awesome_service_io"}[2m]))`
- The total number of failed requests: `sum(increase(http_request_total{host="awesome_service_io", code=~"5.."}[2m]))`

By expresing what are the total count on SLI result processing and the error ratio processed the operator will generate the SLO metrics for this service.

Like is seen in the above output the operator generates 3 metrics:

- `service_level_sli_result_error_ratio_total`: The _downtime/error_ ratio (0-1) of the service.
- `service_level_sli_result_count_total`: The total count of SLI processed total, in other words, what would be the ratio if the service would be 100% correct all the time becasue ratios are from 0 to 1.
- `service_level_slo_objective_ratio`: The objective of the SLO in ratio. This metrics is't processed at all (only changed to ratio unit), but is important to create error budget quries, alerts...

With these metrics we can build availability graphs based on % and error budget burns.

The approach of using counters (instead of gauges) to store the total counts and the error/downtime total gives us the ability to get SLO/SLI rates, increments, speed... in the different time ranges (check query examples section) and is safer in case of missed scrapes, SLI calculation errors... In other words this approach gives us flexibility and safety.

Is important to note that like every metrics this is not exact and is a aproximation (good one but an approximation after all)

## Supported input/output backends

### Input

- [Prometheus]

### Output

- [Prometheus]

## Query examples

## Availability level rate

This will output the availability rate of a service based.

```text
1 - (
    rate(service_level_sli_result_error_ratio_total[1m])
    /
    rate(service_level_sli_result_count_total[1m])
) * 100
```

## Availability level in the last 24h

This will output the availability rate of a service based.

```text
1 - (
    increase(service_level_sli_result_error_ratio_total[24h])
    /
    increase(service_level_sli_result_count_total[24h])
) * 100
```

## Error budget

Calculating the error budget is a little bit more tricky.

### Context

- Taking the previous example we are calculating error budget based on 1 month, this are 43200m (30 \* 24 \* 60).
- Our SLO objective is 99.99 (in ratio: 0.9998999999999999)
- Error budget is based in a 100% for 30d that decrements when availability is less than 99.99% (like the SLO specifies).

### Query

```text
(
  (
    (1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_sli_result_count_total[1m])
    -
    increase(service_level_sli_result_error_ratio_total[${range}])
  )
  /
  (
    (1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_sli_result_count_total[1m])
  )
) * 100
```

Let's decompose the query.

#### Query explanation

`(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_sli_result_count_total[1m])` is the total ratio measured in 1m (sucess + failures) multiplied by the number of minutes in a month and the error budget ratio(1-0.9998999999999999). In other words this is the total (sum) number of error budget for 1 month we have.

`increase(service_level_sli_result_error_ratio_total[${range}])` this is the SLO error sum that we had in ${range} (range changes over time, the first day of the month will be 1d, the 15th of the month will be 15d).

So `(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_sli_result_count_total[1m]) - increase(service_level_sli_result_error_ratio_total[${range}])` returns the number of remaining error budget we have after `${range}`.

If we take that last part and divide for the total error budget we have for the month (`(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_sli_result_count_total[1m])`) this returns us a ratio of the error budget consumed. Multiply by 100 and we have the percent of error budget consumed after `${range}`.

[sre-book-slo]: https://landing.google.com/sre/book/chapters/service-level-objectives.html
[prometheus]: https://prometheus.io/
