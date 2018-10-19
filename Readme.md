# service-level-operator

Service level operator is a operator that based on CRDs creates Service level metrics for services.

This operator doesn't interact with Kubernetes except using the CRDs as a way to define application service levels and generating output metrics based on the CRD set metrics.

Although this operator is though to interact with different backends, at this moment only uses [Prometheus] as input and output.

## How does it work

For this example the output and input backend will be [Prometheus].

First you will need to define a CRD with your service SLOs. In this case we have a service that has an SLO on 99.99 availability, and the SLI is that 5xx are considered errors for this SLO.

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

The operator will query and create new metrics based on the SLOs caulculations at regular intervals (see `--resync-seconds` flag).

The approach that has been taken to generate the SLO metrics is based on [how Google uses and manages SLIs, SLOs and error budgets][sre-book-slo]

In the manifest the SLI is made of 2 prometheus metrics:

- The total of requests: `sum(increase(http_request_total{host="awesome_service_io"}[2m]))`
- The total number of failed requests: `sum(increase(http_request_total{host="awesome_service_io", code=~"5.."}[2m]))`

By expresing what are the total (success and failed) amount and the error amount the operator will generate the SLO metrics for this service.

Output example:

```text
# HELP service_level_slo_error_ratio_total Is the total error ratio counter for SLOs.
# TYPE service_level_slo_error_ratio_total counter
service_level_slo_error_ratio_total{service_level="awesome-service",slo="9999_http_request_lt_500"} 0.40508550763795764
# HELP service_level_slo_full_ratio_total Is the full SLOs ratio counter in time.
# TYPE service_level_slo_full_ratio_total counter
service_level_slo_full_ratio_total{service_level="awesome-service",slo="9999_http_request_lt_500"} 1708
# HELP service_level_slo_objective_ratio Is the objective of the SLO.
# TYPE service_level_slo_objective_ratio gauge
service_level_slo_objective_ratio{service_level="awesome-service",slo="9999_http_request_lt_500"} 0.9998999999999999
```

Like is seen in the above output the operator generates 3 metrics:

- `service_level_slo_error_ratio_total`: The _downtime/error_ ratio of the service.
- `service_level_slo_full_ratio_total`: The total ratio if the service, in other words what would be the ratio if the service would be 100% correct all the time.
- `service_level_slo_objective_ratio`: The objective of the SLO in ratio.

Every metric is based on ratios (0-1).

With this metrics we can build availability graphs based on % and error budget burns.

**Is important to note that like every metrics this is not exact and is a aproximation (good one but an approximation after all)**

## Supported input/output backends

### Input

- [Prometheus]

### Output

- [Prometheus]

## Querie examples

## Availability level

This will output the availability rate of a service based on its SLO.

```text
1 - (
    rate(service_level_slo_error_ratio_total[1m])
    /
    rate(service_level_slo_full_ratio_total[1m])
) * 100
```

## Error budget

Calculating the error budget is a little bit more tricky but it can be get with this kind of queries.

### Context

- Taking the previous example we are calculating error budget based on 1 month, this are 43200m (30 \* 24 \* 60).
- Our SLO objective is 99.99 (in ratio: 0.9998999999999999)
- Error budget is based a 100% that decrements as availability ratio is not 100%.

### Query

```text
(
  (
    (1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_slo_full_ratio_total[1m])
    -
    increase(service_level_slo_error_ratio_total[${range}])
  )
  /
  (
    (1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_slo_full_ratio_total[1m])
  )
) * 100
```

Let's decompose the query.

#### Query explanation

`(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_slo_full_ratio_total[1m])` is the total ratio measured in 1m (sucess + failures) multiplied by the number of minutes in a month and the error budget ratio(1-0.9998999999999999). In other words this is the total (sum) number of error budget for 1 month we have.

`increase(service_level_slo_error_ratio_total[${range}])` this is the SLO error sum that we had in ${range} (range changes over time, the first day of the month will be 1d, the 15th of the month will be 15d).

So `(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_slo_full_ratio_total[1m]) - increase(service_level_slo_error_ratio_total[${range}])` returns the number of remaining error budget we have after `${range}`.

If we take that last part and divide for the total error budget we have for the month (`(1 - service_level_slo_objective_ratio) * 43200 * increase(service_level_slo_full_ratio_total[1m])`) this returns us a ratio of the error budget consumed. Multiply by 100 and we have the percent of error budget consumed after `${range}`.

[sre-book-slo]: https://landing.google.com/sre/book/chapters/service-level-objectives.html
[prometheus]: https://prometheus.io/
