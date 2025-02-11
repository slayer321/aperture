local lib = import '../../grafana/grafana.libsonnet';
local grafana = import 'github.com/grafana/grafonnet-lib/grafonnet/grafana.libsonnet';

local dashboard = grafana.dashboard;
local row = grafana.row;
local prometheus = grafana.prometheus;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local tablePanel = grafana.tablePanel;
local barGaugePanel = grafana.barGaugePanel;
local statPanel = grafana.statPanel;
local annotation = grafana.annotation;
local timeSeriesPanel = lib.TimeSeriesPanel;

local defaults = {
  policyName: error 'policyName must be set',
  refreshInterval: '10s',
  datasourceName: '$datasource',
  datasourceFilterRegex: '',
};

local newTimeSeriesPanel(graphTitle, datasource, graphQuery, axisLabel='', unit='') =
  local target =
    prometheus.target(graphQuery, intervalFactor=1) +
    {
      range: true,
      editorMode: 'code',
    };

  timeSeriesPanel.new(
    title=graphTitle,
    datasource=datasource,
    span=24,
    min_span=24,
    axis_label=axisLabel,
  ) +
  timeSeriesPanel.withTarget(target) +
  timeSeriesPanel.withOptions(
    timeSeriesPanel.options.withLegend() +
    timeSeriesPanel.options.withTooltip()
  ) +
  timeSeriesPanel.withDefaults(
    timeSeriesPanel.defaults.withColorMode('palette-classic') +
    timeSeriesPanel.defaults.withCustom({
      drawStyle: 'line',
      lineInterpolation: 'linear',
      lineWidth: 1,
      pointSize: 5,
      scaleDistribution: {
        type: 'linear',
      },
      showPoints: 'auto',
      spanNulls: false,
    })
  ) +
  timeSeriesPanel.defaults.withThresholds({
    mode: 'absolute',
    steps: [
      { color: 'green', value: null },
      { color: 'red', value: 80 },
    ],
  }) +
  timeSeriesPanel.withFieldConfig(
    timeSeriesPanel.fieldConfig.withDefaults(
      timeSeriesPanel.fieldConfig.defaults.withUnit(unit) +
      timeSeriesPanel.fieldConfig.defaults.withColor({
        mode: 'palette-classic',
      }) +
      timeSeriesPanel.fieldConfig.defaults.withCustom({
        [if axisLabel != '' then 'axisLabel']: axisLabel,
        axisPlacement: 'auto',
        barAlignment: 0,
        fillOpacity: 0,
        gradientMode: 'none',
        drawStyle: 'line',
        hideFrom: {
          legend: false,
          tooltip: false,
          viz: false,
        },
        lineInterpolation: 'linear',
        lineWidth: 1,
        pointSize: 5,
        scaleDistribution: {
          type: 'linear',
        },
        showPoints: 'auto',
        spanNulls: false,
        stacking: {
          group: 'A',
          mode: 'none',
        },
        thresholdsStyle: {
          mode: 'off',
        },
      }) +
      timeSeriesPanel.fieldConfig.defaults.withThresholds({
        mode: 'absolute',
        steps: [
          { color: 'green', value: null },
          { color: 'red', value: 80 },
        ],
      })
    )
  );

local newBarGaugePanel(graphTitle, datasource, graphQuery) =
  local target =
    prometheus.target(graphQuery) +
    {
      legendFormat: '{{ instance }} - {{ policy_name }}',
      format: 'time_series',
      instant: false,
      range: true,
    };

  barGaugePanel.new(
    title=graphTitle,
    datasource=datasource,
  ).addTarget(target) +
  {
    fieldConfig: {
      defaults: {
        color: {
          mode: 'thresholds',
        },
        mappings: [],
        thresholds: {
          mode: 'absolute',
          steps: [
            { color: 'green', value: null },
          ],
        },
      },
      overrides: [],
    },
    options: {
      displayMode: 'gradient',
      minVizHeight: 10,
      minVizWidth: 0,
      orientation: 'horizontal',
      reduceOptions: {
        calcs: ['lastNotNull'],
        fields: '',
        values: false,
      },
      showUnfilled: true,
    },
  };

local newStatPanel(graphTitle, datasource, graphQuery) =
  local target =
    prometheus.target(graphQuery) +
    {
      legendFormat: '{{ instance }} - {{ policy_name }}',
      editorMode: 'code',
      range: true,
    };

  statPanel.new(
    title=graphTitle,
    datasource=datasource,
  ).addTarget(target) +
  {
    fieldConfig: {
      defaults: {
        color: {
          mode: 'thresholds',
        },
        mappings: [],
        thresholds: {
          mode: 'absolute',
          steps: [
            { color: 'green', value: null },
          ],
        },
      },
      overrides: [],
    },
    options: {
      colorMode: 'value',
      graphMode: 'none',
      justifyMode: 'center',
      orientation: 'horizontal',
      reduceOptions: {
        calcs: ['lastNotNull'],
        fields: '',
        values: false,
      },
      textMode: 'auto',
    },
  };

function(params) {
  _config:: defaults + params,

  local p = 'service_latency',
  local ds = $._config.datasourceName,

  local fluxMeterPanel =
    newTimeSeriesPanel('FluxMeter',
                       ds,
                       |||
                         sum(increase(flux_meter_sum{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="%(policyName)s"}[$__rate_interval]))
                         / sum(increase(flux_meter_count{decision_type!="DECISION_TYPE_REJECTED", response_status="OK", flux_meter_name="%(policyName)s"}[$__rate_interval]))
                       ||| % { policyName: $._config.policyName },
                       'Latency (ms)',
                       'ms'),
  local WFQSchedulerFlows =
    newBarGaugePanel('WFQ Scheduler Flows', ds, 'avg(wfq_flows_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local WFQSchedulerHeapRequests =
    newBarGaugePanel('WFQ Scheduler Heap Requests', ds, 'avg(wfq_requests_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local TotalBucketLoadSchedFactor =
    newStatPanel('Total Bucket Load Sched Factor', ds, 'avg(token_bucket_lsf_ratio{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }),

  local TokenBucketBucketCapacity =
    newStatPanel('Token Bucket Bucket Capacity', ds, 'avg(token_bucket_capacity_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName })
    + {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketBucketFillRate =
    newStatPanel('Token Bucket Bucket FillRate', ds, 'avg(token_bucket_fill_rate{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local TokenBucketAvailableTokens =
    newStatPanel('Token Bucket Available Tokens', ds, 'avg(token_bucket_available_tokens_total{policy_name="%(policyName)s"})' % { policyName: $._config.policyName }) +
    {
      options+: {
        orientation: 'auto',
      },
    },

  local IncomingConcurrency =
    newTimeSeriesPanel('Incoming Concurrency', ds, 'sum(rate(incoming_concurrency_ms{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Concurrency', 'ms'),

  local AcceptedConcurrency =
    newTimeSeriesPanel('Accepted Concurrency', ds, 'sum(rate(accepted_concurrency_ms{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Concurrency', 'ms'),

  local WorkloadDecisions =
    newTimeSeriesPanel('Workload Decisions', ds, 'sum by(workload_index, decision_type) (rate(workload_latency_ms_count{policy_name="%(policyName)s"}[$__rate_interval]))' % { policyName: $._config.policyName }, 'Decisions', 'reqps'),

  local WorkloadLatency =
    newTimeSeriesPanel('Workload Latency (Auto Tokens)', ds, '(sum by (workload_index) (increase(workload_latency_ms_sum{policy_name="%(policyName)s",decision_type!="DECISION_TYPE_REJECTED"}[$__rate_interval])))/(sum by (workload_index) (increase(workload_latency_ms_count{policy_name="%(policyName)s",decision_type!="DECISION_TYPE_REJECTED"}[$__rate_interval])))' % { policyName: $._config.policyName }, 'Latency', 'ms'),


  local dashboardDef =
    dashboard.new(
      title='Jsonnet / FluxNinja',
      editable=true,
      schemaVersion=18,
      refresh=$._config.refreshInterval,
      time_from='now-5m',
      time_to='now'
    )
    .addTemplate(
      {
        current: {
          text: 'default',
          value: $._config.datasourceName,
        },
        hide: 0,
        label: 'Data Source',
        name: 'datasource',
        options: [],
        query: 'prometheus',
        refres: 1,
        regex: $._config.datasourceFilterRegex,
        type: 'datasource',
      }
    )
    .addPanel(fluxMeterPanel, gridPos={ h: 10, w: 24, x: 0, y: 0 })
    .addPanel(WorkloadDecisions, gridPos={ h: 10, w: 24, x: 0, y: 10 })
    .addPanel(WorkloadLatency, gridPos={ h: 10, w: 24, x: 0, y: 20 })
    .addPanel(IncomingConcurrency, gridPos={ h: 8, w: 12, x: 0, y: 30 })
    .addPanel(AcceptedConcurrency, gridPos={ h: 8, w: 12, x: 12, y: 30 })
    .addPanel(WFQSchedulerFlows, gridPos={ h: 3, w: 8, x: 0, y: 40 })
    .addPanel(TotalBucketLoadSchedFactor, gridPos={ h: 6, w: 4, x: 8, y: 40 })
    .addPanel(TokenBucketBucketCapacity, gridPos={ h: 6, w: 4, x: 12, y: 40 })
    .addPanel(TokenBucketBucketFillRate, gridPos={ h: 6, w: 4, x: 16, y: 40 })
    .addPanel(TokenBucketAvailableTokens, gridPos={ h: 6, w: 4, x: 20, y: 40 })
    .addPanel(WFQSchedulerHeapRequests, gridPos={ h: 3, w: 8, x: 0, y: 40 }),

  dashboard: dashboardDef,
}
