import "../styles/main.css"

import Search from './search.js'
import Chart from './chart.js'
import ChartSubtitles from './chart/subtitles.js'
import ChartTooltips from './chart/tool_tip.js'

window.WowTrends = {};
window.WowTrends.Search = Search();
window.WowTrends.Chart = Chart();
window.WowTrends.Chart.subtitles = ChartSubtitles();
window.WowTrends.Chart.toolTip = ChartTooltips();

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

window.addEventListener("load", function() {
  window.WowTrends.Search.init();
  window.WowTrends.Chart.init();
});
