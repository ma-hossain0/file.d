# Meric plugin
Metric plugin.

**Example:**
```yaml
pipelines:
  example_pipeline:
    ...
    actions:
    - type: metric
	  metric_name: errors_total
	  metric_labels:
	  	- level

    ...
```


### Config params
**`metric_name`** *`string`* *`default=total`* 

The metric name.

<br>

**`metric_labels`** *`[]string`* 

Lists the event fields to add to the metric. Blank list means no labels.
Important note: labels metrics are not currently being cleared.

<br>

**`metric_name`** *`string`* 

The metric name.

<br>

**`metric_labels`** *`[]string`* 

Lists the event fields to add to the metric. Blank list means no labels.
Important note: labels metrics are not currently being cleared.

<br>


<br>*Generated using [__insane-doc__](https://github.com/vitkovskii/insane-doc)*