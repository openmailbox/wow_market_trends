var Chart = (function() {
    var data     = [];
    var chart    = null;
    var itemId   = new URL(window.location.href).searchParams.get("itemId");
    var itemName = null;

    var callback = function (evt) {
        data = JSON.parse(this.response, parseDate);
        itemName = data[0].name;
        drawChart();
    };

    var drawChart = function() {
        chart = new CanvasJS.Chart("chartContainer",
            {
                toolTip: {
                    contentFormatter: function(e) {
                        var point = e.entries[0].dataPoint;
                        var str   = "";

                        str += "<strong>" + point.label + "</strong>";
                        str += "<br />";
                        str += "<strong>Open:</strong> " + formatPriceLong(point.y[0]) + "<br />";
                        str += "<strong>High:</strong> " + formatPriceLong(point.y[1]) + "<br />";
                        str += "<strong>Low:</strong> " + formatPriceLong(point.y[2]) + "<br />";
                        str += "<strong>Close:</strong> " + formatPriceLong(point.y[3]);

                        return str;
                    }
                },
                title: {
                    horizontalAlign: "left",
                    padding: 20,
                    text: itemName,
                },
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
                    }
                ]
            });
        chart.render();
    }

    var formatData = function () {
        return data.map(function (i) {
            return { x: i.created_at, y: [i.open, i.high, i.low, i.close], label: i.created_at }
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
        oReq.open("GET", "history?itemId=" + itemId);
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
        getChart: function () { return chart; },
        getData: function () { return data; },
        init: init
    };
})();
