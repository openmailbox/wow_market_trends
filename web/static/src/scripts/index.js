import "../styles/main.css"

import Chart from './chart.js'
import Search from './search.js'

window.WowTrends = {};
window.WowTrends.Search = Search();
window.WowTrends.Chart = Chart();

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