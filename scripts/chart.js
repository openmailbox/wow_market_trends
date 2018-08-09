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
                title: {
                    text: "Price History for Item " + itemName
                },
                zoomEnabled: true,
                axisY: {
                    includeZero: false,
                    title: "Prices",
                    prefix: ""
                },
                axisX: {
                    interval: 24,
                    intervalType: "hour",
                    valueFormatString: "MMM-DD",
                    labelAngle: -45
                },
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
            return { x: i.created_at, y: [i.open, i.high, i.low, i.close] }
        });
    };

    var init = function () {
        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", callback);
        //oReq.open("GET", "https://yara.open-mailbox.com/wow?itemId=" + itemId);
        oReq.open("GET", "http://localhost:8081/history?itemId=" + itemId);
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

window.onload = function () {
    Chart.init();
}
