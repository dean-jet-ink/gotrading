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
  config.sma.indexes = [];
  config.sma.values = [];
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

      const candles = data["candles"];
      const rows = [];
      for (let i = 0; i < candles.length; i++) {
        const date = new Date(candles[i].time);
        const row = [
          date,
          candles[i].low,
          candles[i].open,
          candles[i].close,
          candles[i].high,
          candles[i].volume,
        ];

        for (let sma of config.sma.values) {
          if (sma[i] == 0) {
            row.push(null);
          } else {
            row.push(sma[i]);
          }
        }

        rows.push(row);
      }

      dataTable.addRows(rows);
      drawChart(dataTable);
    });
}

window.onload = () => {
  send();

  setInterval(send, config.api.interval);

  const dashboard = document.getElementById("dashboard_div");
  dashboard.addEventListener("mouseenter", () => {
    config.api.enable = false;
  });
  dashboard.addEventListener("mouseleave", () => {
    config.api.enable = true;
  });

  document.getElementById("sma").addEventListener("change", (event) => {
    if (event.target.checked) {
      config.sma.enable = true;
    } else {
      config.sma.enable = false;
    }
    send();
  });
  for (let i = 0; i < 3; i++) {
    const period = document.getElementById(`smaPeriod${i + 1}`);
    period.addEventListener("input", (event) => {
      config.sma.periods[i] = event.target.value;
      send();
    });
    config.sma.periods[i] = period.value;
  }
};
