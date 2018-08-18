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

        // TODO: Better icon handling
        var element = document.createElement("img");
        element.setAttribute("src", _iconUrl);
        document.querySelector("#chart-main").appendChild(element);
    };

    var draw = function() {
        _chart = new CanvasJS.Chart("chart-container",
            {
                toolTip: WowTrends.Chart.toolTip,
                title: {
                    horizontalAlign: "left",
                    padding: { top: 0, left: 20, right: 0, bottom: 0 },
                    text: _data.name,
                    fontSize: 30
                },
                subtitles: WowTrends.Chart.subtitles.build(),
                zoomEnabled: true,
                axisY: {
                    includeZero: false,
                    title: "Price",
                    prefix: "",
                    labelFormatter: function(e) {
                        if (e.value > 10000) return CanvasJS.formatNumber(Math.floor(e.value / 10000)) + "G";
                        if (e.value > 100) return Math.floor(e.value / 100) + "S";
                        return e.value + "C";
                    }
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
                    interval: 8,
                    intervalType: "hour",
                    valueFormatString: "DD MMM H:mm",
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

    var formatVolumeData = function() {
        return _data.periods.map(function(i) {
            return { x: i.created_at, y: i.volume };
        });
    };

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

    return {
        formatPriceLong: formatPriceLong,
        getChart: function () { return _chart; },
        getData: function () { return _data; },
        init: init
    };
})();
