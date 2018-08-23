var WowTrends = WowTrends || {};

WowTrends.Chart.subtitles = (function() {
    var priceWithChange = function() {
        var data        = WowTrends.Chart.getData();
        var difference  = data.periods[0].ask - data.periods[0].open;
        var plusOrMinus = difference < 0 ? "" : "+";
        var percentage  = ((difference / data.periods[0].ask) * 100).toFixed(1);
        var text        = WowTrends.Chart.formatPriceLong(data.current) + " (" + plusOrMinus + (difference / 10000) + "G / " + percentage + "%)";

        return {
            text: text,
            horizontalAlign: "left",
            fontFamily: 'Segoe UI, Roboto, Helvetica Neue, sans-serif',
            fontSize: 20,
            padding: { top: 0, left: 20, right: 0, bottom: 10 }
        };
    };

    return {
        build: function() { return [priceWithChange()] }
    };
})();
