var callback = function() {
  console.log(this.responseText);
}

window.onload = function () {
    var oReq = new XMLHttpRequest();
    oReq.addEventListener("load", callback);
    //oReq.open("GET", "http://wow.open-mailbox.com/history?itemId=55550");
    oReq.open("GET", "http://localhost:8081/history?itemId=36211");
    oReq.send();

	var chart = new CanvasJS.Chart("chartContainer", {
		title:{
			text: "My First Chart in CanvasJS"              
		},
		data: [              
		{
			// Change type to "doughnut", "line", "splineArea", etc.
			type: "column",
			dataPoints: [
				{ label: "apple",  y: 10  },
				{ label: "orange", y: 15  },
				{ label: "banana", y: 25  },
				{ label: "mango",  y: 30  },
				{ label: "grape",  y: 28  }
			]
		}
		]
	});
	chart.render();
}
