var WowTrends = WowTrends || {};

WowTrends.Chart.toolTip = (function() {
    /**
     * @param {Object} event - The chart data provided by CanvasJS.
     * @param {Object[]} event.entries - All of the data series available to the chart.
     * @param {Object} event.entries[].dataPoint - An individual data point.
     */
    var getToolTip = function(event) {
        var point = event.entries[0].dataPoint;
        var str = "";

        str += "<strong>" + point.label + "</strong>";
        str += "<br />";
        str += "<strong>Open:</strong> " + point.y[0] / 10000 + "G<br />";
        str += "<strong>High:</strong> " + point.y[1] / 10000 + "G<br />";
        str += "<strong>Low:</strong> " + point.y[2] / 10000 + "G<br />";
        str += "<strong>Close:</strong> " + point.y[3] / 10000 + "G<br />";
        str += "<strong>Volume:</strong> " + event.entries[1].dataPoint.y;

        return str;
    }

    return {
        shared: true,
        contentFormatter: getToolTip
    };
})();