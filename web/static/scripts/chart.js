var WowTrends = WowTrends || {};

WowTrends.Chart = (function() {
    var BASE_ICON_URL = "https://wow.zamimg.com/images/wow/icons/large/";
    var _chart        = null;
    var _data         = null;
    var _iconUrl      = null;
    var _itemId       = new URL(window.location.href).searchParams.get("itemId");
    var _maxVolume    = null;

    var callback = function (evt) {
        _data      = JSON.parse(this.response, parseDate);
        _iconUrl   = BASE_ICON_URL + _data.icon + ".jpg";
        _maxVolume = Math.max(..._data.periods.map(function(i) { return i.volume; }));

        draw();

        var element = document.createElement("img");
        var parent  = document.querySelector("#chart-main");
        var sibling = document.querySelector("#chart-container");

        element.setAttribute("src", _iconUrl);
        element.classList.add("chart-icon");
        parent.insertBefore(element, sibling);

        populateTableData();
    };

    var draw = function() {
        _chart = new CanvasJS.Chart("chart-container",
            {
                colorSet: "spectreColorSet",
                toolTip: WowTrends.Chart.toolTip,
                title: {
                    fontColor: "#454d5d",
                    horizontalAlign: "left",
                    padding: { top: 0, left: 20, right: 0, bottom: 0 },
                    text: _data.name,
                    fontFamily: 'Segoe UI, Roboto, Helvetica Neue, sans-serif',
                    fontSize: 30,
                    fontWeight: "bold"
                },
                subtitles: WowTrends.Chart.subtitles.build(),
                zoomEnabled: true,
                axisY: {
                    includeZero: false,
                    title: "Price",
                    prefix: "",
                    labelFormatter: function(e) { return formatPrice(e.value); }
                },
                axisY2: {
                    title: "Volume",
                    includeZero: true,
                    maximum: 3 * _maxVolume
                },
                axisX: {
                    scaleBreaks: {
                        autoCalculate: true,
                        collapsibleThreshold: "10%"
                    },
                    interval: 24,
                    intervalType: "hour",
                    valueFormatString: "DD MMM",
                    labelAngle: -45,
                },
                zoomEnabled: true,
                data: [
                    {
                        type: "candlestick",
                        dataPoints: formatData()
                    },
                    {
                        type: "line",
                        axisYType: "secondary",
                        dataPoints: formatVolumeData()
                    },
                    {
                        type: "line",
                        axisYType: "primary",
                        dataPoints: formatAverageData()
                    }
                ]
            });
        _chart.render();
    }

    var formatData = function () {
        return _data.periods.map(function (i) {
            return { x: i.created_at, y: [i.open, i.high, i.low, i.close], label: i.created_at };
        });
    };

    var formatAverageData = function() {
        return _data.periods.map(function(i) {
            return { x: i.created_at, y: i.average }
        });
    };

    var formatVolumeData = function() {
        return _data.periods.map(function(i) {
            return { x: i.created_at, y: i.volume };
        });
    };

    var formatPrice = function(copper) {
        if (copper > 10000) return CanvasJS.formatNumber(Math.floor(copper / 10000)) + "G";
        if (copper > 100) return Math.floor(copper / 100) + "S";
        return copper + "C";
    }

    /**
     * Format a price as "X gold, Y silver, Z copper"
     * @param {number} copper - The price in copper pieces.
     * @returns {string}
     */
    var formatPriceLong = function(copper) {
        var gold = CanvasJS.formatNumber(Math.floor(copper / 10000));

        copper = copper % 10000

        var silver = Math.floor(copper / 100) % 100;

        copper = copper % 100;

        return gold + " gold, " + silver + " silver, " + copper + " copper";
    };

    var init = function () {
        if (_itemId === null) return;

        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", callback);
        oReq.open("GET", "details?itemId=" + _itemId);
        oReq.send();
    };

    var parseDate = function (key, value) {
        var a;

        if (key === 'created_at' && typeof value === 'string') {
            a = /(\d+)-(\d+)-(\d+)T(\d+):(\d+):(\d+)/.exec(value)
            if (a) {
                return new Date(a[1], a[2] - 1, a[3], a[4], a[5], a[6]);
            }
        }

        return value;
    }

    var populateTableData = function() {
        var highFourteen = Math.max(..._data.periods.slice(0,14).map(function(i) { return i.high }));
        var lowFourteen = Math.max(..._data.periods.slice(0,14).map(function(i) { return i.low }));

        document.getElementById("data-table-open").innerText = formatPrice(_data.periods[0].open);
        document.getElementById("data-table-close").innerText = formatPrice(_data.periods[0].close);
        document.getElementById("data-table-high").innerText = formatPrice(_data.periods[0].high);
        document.getElementById("data-table-low").innerText = formatPrice(_data.periods[0].low);
        document.getElementById("data-table-high-14").innerText = formatPrice(highFourteen);
        document.getElementById("data-table-low-14").innerText = formatPrice(lowFourteen);
    };

    return {
        formatPriceLong: formatPriceLong,
        getChart: function () { return _chart; },
        getData: function () { return _data; },
        init: init
    };
})();

window.addEventListener("load", function() {
    WowTrends.Chart.init();
});

CanvasJS.addColorSet("spectreColorSet", 
                     [
                       "#5755d9",
                       "#32b643",
                       "#ffb700",
                       "#f1f1fc",
                       "#e85600",
                       "#454d5d",
                       "#727e96",
                       "#acb3c2",
                       "#e7e9ed",
                       "#f0f1f4",
                       "#f8f9fa",
                       "#fff"
                     ]);
