var Chart = (function() {
    var baseIconUrl  = "https://wow.zamimg.com/images/wow/icons/large/";
    var chart        = null;
    var currentPrice = null
    var data         = null;
    var iconUrl      = null;
    var itemId       = new URL(window.location.href).searchParams.get("itemId");
    var itemName     = null;
    var maxVolume    = null;

    var callback = function (evt) {
        data = JSON.parse(this.response, parseDate);
        itemName = data.name;
        iconUrl = baseIconUrl + data.icon + ".jpg";
        maxVolume = Math.max(...data.periods.map(function(i) { return i.volume; }));
        currentPrice = formatPriceLong(data.current);
        drawChart();

        //TODO: Something else
        var element = document.createElement("img");
        element.setAttribute("src", iconUrl);
        document.querySelector("#chart-main").appendChild(element);
    };

    var drawChart = function() {
        chart = new CanvasJS.Chart("chartContainer",
            {
                toolTip: {
                    shared: true,
                    contentFormatter: function(e) {
                        var point = e.entries[0].dataPoint;
                        var str   = "";

                        str += "<strong>" + point.label + "</strong>";
                        str += "<br />";
                        str += "<strong>Open:</strong> " + formatPriceLong(point.y[0]) + "<br />";
                        str += "<strong>High:</strong> " + formatPriceLong(point.y[1]) + "<br />";
                        str += "<strong>Low:</strong> " + formatPriceLong(point.y[2]) + "<br />";
                        str += "<strong>Close:</strong> " + formatPriceLong(point.y[3]) + "<br />";
                        str += "<strong>Volume:</strong> " + e.entries[1].dataPoint.y;

                        return str;
                    }
                },
                title: {
                    horizontalAlign: "left",
                    padding: {
                        top: 0,
                        left: 20,
                        right: 0,
                        bottom: 0
                    },
                    text: itemName,
                    fontSize: 30
                },
                subtitles: [
                    {
                        text: currentPrice,
                        horizontalAlign: "left",
                        fontSize: 20,
                        padding: {
                            top: 0,
                            left: 20,
                            right: 0,
                            bottom: 10
                        }
                    }
                ],
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
                    maximum: 3 * maxVolume
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
        chart.render();
    }

    var formatData = function () {
        return data.periods.map(function (i) {
            return { x: i.created_at, y: [i.open, i.high, i.low, i.close], label: i.created_at };
        });
    };

    var formatVolumeData = function() {
        return data.periods.map(function(i) {
            return { x: i.created_at, y: i.volume };
        });
    };

    var formatPriceLong = function(copper) {
        var gold = CanvasJS.formatNumber(Math.floor(copper / 10000));

        copper = copper % 10000

        var silver = Math.floor(copper / 100) % 100;

        copper = copper % 100;

        return gold + " gold, " + silver + " silver, " + copper + " copper";
    }

    var init = function () {
        if (itemId === null) return;

        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", callback);
        oReq.open("GET", "details?itemId=" + itemId);
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
        getAuctions: function () { return auctions; },
        getChart: function () { return chart; },
        getData: function () { return data; },
        init: init
    };
})();
