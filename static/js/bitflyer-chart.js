google.charts.load("current", { packages: ["corechart", "controls"] });

const config = {
  api: {
    enable: true,
    interval: 1000 * 3,
  },
  candleStick: {
    productCode: "BTC_JPY",
    duration: "1m",
    limit: 365,
    indexOfColumn: 5,
  },
  dataTable: {
    index: 0,
    value: null,
  },
  sma: {
    enable: false,
    periods: [],
    values: [],
    indexes: [],
  },
  ema: {
    enable: false,
    periods: [],
    values: [],
    indexes: [],
  },
  bbands: {
    enable: false,
    period: 20,
    k: 2,
    maType: "sma",
    values: [],
    indexes: [],
  },
  volume: {
    enable: false,
    rendered: false,
  },
  rsi: {
    enable: false,
    period: 14,
    upper: 70,
    lower: 30,
    values: [],
    indexes: [],
    rendered: false,
  },
  macd: {
    enable: false,
    fastPeriod: 9,
    slowPeriod: 26,
    signalPeriod: 12,
    values: [],
    indexes: [],
    rendered: false,
  },
  hv: {
    enable: false,
    periods: [],
    values: [],
    indexes: [],
    rendered: false,
  },
  event: {
    enable: false,
    signals: [],
    indexes: [],
  },
};

function drawChart(dataTable) {
  const dashboardDiv = document.getElementById("dashboard_div");
  const charts = [];
  const dashboard = new google.visualization.Dashboard(dashboardDiv);
  const mainChart = new google.visualization.ChartWrapper({
    chartType: "ComboChart",
    containerId: "chart_div",
    options: {
      hAxis: { slantedText: false },
      legend: { position: "none" },
      candlestick: {
        fallingColor: { strokeWidth: 0, fill: "#a52714" },
        risingColor: { strokeWidth: 0, fill: "#0f9d58" },
      },
      seriesType: "candlesticks",
      series: {},
    },
    view: {
      columns: [
        {
          calc: function (d, rowIndex) {
            return d.getFormattedValue(rowIndex, 0);
          },
          type: "string",
        },
        1,
        2,
        3,
        4,
      ],
    },
  });
  charts.push(mainChart);

  const options = mainChart.getOptions();
  const view = mainChart.getView();

  for (let index of config.sma.indexes) {
    options.series[index] = { type: "line" };
    view.columns.push(config.candleStick.indexOfColumn + index);
  }

  for (let index of config.ema.indexes) {
    options.series[index] = { type: "line" };
    view.columns.push(config.candleStick.indexOfColumn + index);
  }

  for (let index of config.bbands.indexes) {
    options.series[index] = { type: "line" };
    view.columns.push(config.candleStick.indexOfColumn + index);
  }

  if (config.event.enable) {
    options.series[config.event.indexes[0]] = {
      type: "line",
      tooltip: "none",
      enableInteractivity: false,
      lineWidth: 0,
    };
    view.columns.push(
      config.candleStick.indexOfColumn + config.event.indexes[0]
    );
    view.columns.push(
      config.candleStick.indexOfColumn + config.event.indexes[1]
    );
  }

  if (config.volume.enable && !config.volume.rendered) {
    document.getElementById("technical_div").innerHTML += `
        <div id="volume_div" class="bottom_chart">
          <span>Volume</span>
          <div id="volume_chart"></div>
        </div>
      `;

    const volumeChart = new google.visualization.ChartWrapper({
      chartType: "ColumnChart",
      containerId: "volume_chart",
      options: {
        hAxis: { slantedText: false },
        legend: { position: "none" },
        series: {},
      },
      view: {
        columns: [{ type: "string" }, 5],
      },
    });
    charts.push(volumeChart);

    config.volume.rendered = true;
  }

  if (config.rsi.enable && !config.rsi.rendered) {
    document.getElementById("technical_div").innerHTML += `
        <div id="rsi_div" class="bottom_chart">
          <span>RSI</span>
          <div id="rsi_chart"></div>
        </div>
      `;

    const upperIndex = config.candleStick.indexOfColumn + config.rsi.indexes[0];
    const valueIndes = config.candleStick.indexOfColumn + config.rsi.indexes[1];
    const lowerIndex = config.candleStick.indexOfColumn + config.rsi.indexes[2];

    var rsiChart = new google.visualization.ChartWrapper({
      chartType: "LineChart",
      containerId: "rsi_chart",
      options: {
        hAxis: { slantedText: false },
        legend: { position: "none" },
        series: {
          0: { color: "black", lineWidth: 1 },
          1: { color: "#e2431e" },
          2: { color: "black", lineWidth: 1 },
        },
      },
      view: {
        columns: [{ type: "string" }, upperIndex, valueIndes, lowerIndex],
      },
    });
    charts.push(rsiChart);
    config.rsi.rendered = true;
  }

  if (config.macd.enable && !config.macd.rendered) {
    document.getElementById("technical_div").innerHTML += `
        <div id="macd_div" class="bottom_chart">
          <span>MACD</span>
          <div id="macd_chart"></div>
        </div>
      `;

    const macdIndex = config.candleStick.indexOfColumn + config.macd.indexes[0];
    const macdSignalIndex =
      config.candleStick.indexOfColumn + config.macd.indexes[1];
    const macdHistIndex =
      config.candleStick.indexOfColumn + config.macd.indexes[2];

    const macdChart = new google.visualization.ChartWrapper({
      chartType: "ComboChart",
      containerId: "macd_chart",
      options: {
        legend: { position: "none" },
        seriesType: "bars",
        series: {
          1: { type: "line", lineWidth: 1 },
          2: { type: "line", lineWidth: 1 },
        },
      },
      view: {
        columns: [
          { type: "string" },
          macdIndex,
          macdSignalIndex,
          macdHistIndex,
        ],
      },
    });

    charts.push(macdChart);
    config.macd.rendered = true;
  }

  if (config.hv.enable && !config.hv.rendered) {
    document.getElementById("technical_div").innerHTML += `
        <div id="hv_div" class="bottom_chart">
          <span>HV</span>
          <div id="hv_chart"></div>
        </div>
      `;

    const series = {};
    const columns = [{ type: "string" }];
    for (let index of config.hv.indexes) {
      series[index] = { lineWidth: 1 };
      columns.push(config.candleStick.indexOfColumn + index);
    }

    const hvChart = new google.visualization.ChartWrapper({
      chartType: "LineChart",
      containerId: "hv_chart",
      options: {
        legend: { position: "none" },
        series: series,
      },
      view: {
        columns: columns,
      },
    });
    charts.push(hvChart);
    config.hv.rendered = true;
  }

  const controlWrapper = new google.visualization.ControlWrapper({
    controlType: "ChartRangeFilter",
    containerId: "filter_div",
    options: {
      filterColumnIndex: 0,
      ui: {
        chartType: "LineChart",
        chartView: {
          columns: [0, 4],
        },
      },
    },
  });

  dashboard.bind(controlWrapper, charts);
  dashboard.draw(dataTable);
}

function initConfig() {
  config.dataTable.index = 0;
  config.dataTable.value = null;
  config.sma.values = [];
  config.sma.indexes = [];
  config.ema.values = [];
  config.ema.indexes = [];
  config.bbands.values = [];
  config.bbands.indexes = [];
  config.rsi.values = [];
  config.rsi.indexes = [];
  config.macd.values = [];
  config.macd.indexes = [];
  config.hv.values = [];
  config.hv.indexes = [];
  config.event.indexes = [];
  config.event.signals = [];
}

function send() {
  if (!config.api.enable) {
    return;
  }

  initConfig();

  const params = new URLSearchParams();
  params.append("product_code", config.candleStick.productCode);
  params.append("duration", config.candleStick.duration);
  params.append("limit", config.candleStick.limit);

  if (config.sma.enable) {
    params.append("sma", true);
    for (let i = 0; i < config.sma.periods.length; i++) {
      params.append(`sma_period${i + 1}`, config.sma.periods[i]);
    }
  }

  if (config.ema.enable) {
    params.append("ema", true);
    for (let i = 0; i < config.ema.periods.length; i++) {
      params.append(`ema_period${i + 1}`, config.ema.periods[i]);
    }
  }

  if (config.bbands.enable) {
    params.append("bbands", true);
    params.append("bbands_period", config.bbands.period);
    params.append("bbands_k", config.bbands.k);
    params.append("bbands_maType", config.bbands.maType);
  }

  if (config.rsi.enable) {
    params.append("rsi", true);
    params.append("rsi_period", config.rsi.period);
  }

  if (config.macd.enable) {
    params.append("macd", true);
    params.append("macd_fastPeriod", config.macd.fastPeriod);
    params.append("macd_slowPeriod", config.macd.slowPeriod);
    params.append("macd_signalPeriod", config.macd.signalPeriod);
  }

  if (config.hv.enable) {
    params.append("hv", true);
    for (let i = 0; i < config.hv.periods.length; i++) {
      params.append(`hv_period${i + 1}`, config.hv.periods[i]);
    }
  }

  if (config.event.enable) {
    params.append("event", true);
  }

  fetch(`/api/candle/?${params}`)
    .then((response) => response.json())
    .then((data) => {
      if (data["error"] != undefined) {
        console.log(`${data["code"]}: ${data["error"]}`);
        return;
      }

      const dataTable = new google.visualization.DataTable();
      dataTable.addColumn("date", "Date");
      dataTable.addColumn("number", "Low");
      dataTable.addColumn("number", "Open");
      dataTable.addColumn("number", "Close");
      dataTable.addColumn("number", "High");
      dataTable.addColumn("number", "Volume");

      if (data["smas"] != undefined) {
        const smas = data["smas"];
        for (let sma of smas) {
          dataTable.addColumn("number", `SMA${sma.period}`);
          config.sma.indexes.push(++config.dataTable.index);
          config.sma.values.push(sma.values);
        }
      }

      if (data["emas"] != undefined) {
        const emas = data["emas"];
        for (let ema of emas) {
          dataTable.addColumn("number", `EMA${ema.period}`);
          config.ema.indexes.push(++config.dataTable.index);
          config.ema.values.push(ema.values);
        }
      }

      if (data["bbands"] != undefined) {
        const bbands = data["bbands"];

        config.bbands.indexes.push(++config.dataTable.index);
        config.bbands.indexes.push(++config.dataTable.index);
        config.bbands.indexes.push(++config.dataTable.index);
        dataTable.addColumn("number", `BBands-Upper${bbands.k}`);
        dataTable.addColumn("number", `BBands-Mid${bbands.period}`);
        dataTable.addColumn("number", `BBands-Lower${bbands.k}`);
        config.bbands.values.push(bbands.upper);
        config.bbands.values.push(bbands.mid);
        config.bbands.values.push(bbands.lower);
      }

      if (data["events"] != undefined) {
        const events = data["events"];

        dataTable.addColumn("number", "Marker");
        dataTable.addColumn({ type: "string", role: "annotation" });

        config.event.indexes.push(++config.dataTable.index);
        config.event.indexes.push(++config.dataTable.index);
        config.event.signals = events.signals;

        if (events["profit"] != undefined) {
          const profit = Math.floor(events["profit"] * 100) / 100;

          document.getElementById("profit").innerHTML = `
            <div>Profit  ${profit}円</div>
          `;
        }
      }

      if (data["rsi"] != undefined) {
        const rsi = data["rsi"];

        config.rsi.indexes.push(++config.dataTable.index);
        config.rsi.indexes.push(++config.dataTable.index);
        config.rsi.indexes.push(++config.dataTable.index);
        dataTable.addColumn("number", `RSI-Upper${config.rsi.upper}%`);
        dataTable.addColumn("number", `RSI-Values${config.rsi.period}`);
        dataTable.addColumn("number", `RSI-Lower${config.rsi.lower}%`);

        config.rsi.values = rsi.values;
      }

      if (data["macd"] != undefined) {
        const macd = data["macd"];

        config.macd.indexes.push(++config.dataTable.index);
        config.macd.indexes.push(++config.dataTable.index);
        config.macd.indexes.push(++config.dataTable.index);
        dataTable.addColumn(
          "number",
          `MACD(${macd.fastPeriod}, ${macd.slowPeriod})`
        );
        dataTable.addColumn("number", `MACD-Signal${macd.signalPeriod}`);
        dataTable.addColumn("number", `MACD-Histgram`);

        config.macd.values.push(macd["values"]);
        config.macd.values.push(macd["signal_values"]);
        config.macd.values.push(macd["histgram"]);
      }

      if (data["hvs"] != undefined) {
        const hvs = data["hvs"];

        for (let i = 0; i < config.hv.periods.length; i++) {
          config.hv.indexes.push(++config.dataTable.index);
          dataTable.addColumn("number", `HV${config.hv.periods[i]}`);
          config.hv.values.push(hvs[i]["values"]);
        }
      }

      const candles = data["candles"];
      const rows = [];
      for (let i = 0; i < candles.length; i++) {
        let candle = candles[i];
        const date = new Date(candle.time);
        const row = [
          date,
          candle.low,
          candle.open,
          candle.close,
          candle.high,
          candle.volume,
        ];

        for (let sma of config.sma.values) {
          if (sma[i] == 0) {
            row.push(null);
          } else {
            row.push(sma[i]);
          }
        }

        for (let ema of config.ema.values) {
          if (ema[i] == 0) {
            row.push(null);
          } else {
            row.push(ema[i]);
          }
        }

        for (let bband of config.bbands.values) {
          if (bband[i] == 0) {
            row.push(null);
          } else {
            row.push(bband[i]);
          }
        }

        if (config.rsi.enable) {
          row.push(config.rsi.upper);
          if (config.rsi.values[i] == 0) {
            row.push(null);
          } else {
            row.push(config.rsi.values[i]);
          }
          row.push(config.rsi.lower);
        }

        for (let value of config.macd.values) {
          if (value[i] == 0) {
            row.push(null);
          } else {
            row.push(value[i]);
          }
        }

        for (let value of config.hv.values) {
          if (value[i] == 0) {
            row.push(null);
          } else {
            row.push(value[i]);
          }
        }

        if (config.event.enable) {
          let signals = config.event.signals;

          if (signals.length != 0 && signals[0].time == candle.time) {
            row.push(candle.high * 1.001);
            row.push(signals[0].side);
            config.event.signals.shift();
          } else {
            row.push(null);
            row.push(null);
          }
        }

        rows.push(row);
      }

      dataTable.addRows(rows);
      config.dataTable.value = dataTable;
      drawChart(dataTable);
    });
}

function changeDuration(duration) {
  config.candleStick.duration = duration;
  send();
}

function initPeriods(indicator) {
  const periods = document.querySelectorAll(`.${indicator}Period`);
  for (let i = 0; i < periods.length; i++) {
    config[indicator].periods[i] = periods[i].value;
  }
}

function switchIndicator(indicator, isChecked) {
  config[indicator].enable = isChecked;
  send();
}

function changePeriods(indicator, period, index) {
  config[indicator].periods[index] = period;
  send();
}

function changePeriod(indicator, period) {
  config[indicator].period = period;
  send();
}

function switchIndicatorOfTechnicalChart(indicator, isChecked, divId) {
  config[indicator].enable = isChecked;
  if (isChecked) {
    send();
  } else {
    config[indicator].rendered = false;
    document.getElementById(divId).remove();
  }
}

window.onload = () => {
  send();

  const initialization = ["sma", "ema", "hv"];

  for (let indicator of initialization) {
    initPeriods(indicator);
  }

  setInterval(send, config.api.interval);

  const dashboard = document.getElementById("dashboard_div");
  dashboard.addEventListener("mouseenter", () => {
    config.api.enable = false;
  });
  dashboard.addEventListener("mouseleave", () => {
    config.api.enable = true;
  });

  document.getElementById("bbandsK").addEventListener("input", (event) => {
    config.bbands.k = event.target.value;
  });

  document.getElementById("volume").addEventListener("change", (event) => {
    config.volume.enable = event.target.checked;
    if (!event.target.checked) {
      document.getElementById("volume_div").remove();
      config.volume.rendered = false;
    } else {
      drawChart(config.dataTable.value);
    }
  });

  document.getElementById("event").addEventListener("change", (event) => {
    config.event.enable = event.target.checked;
    if (!config.event.enable) {
      document.getElementById("profit").innerHTML = "";
    }
    send();
  });
};
