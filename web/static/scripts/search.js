var Search = (function() {
    var timer = null;

    var callback = function(event) {
        data = JSON.parse(this.response);

        var list = document.querySelector("#search-results");

        while (list.firstChild !== null) {
            list.firstChild.remove();
        }

        if (data === null) return;

        list.style.display = "inherit";

        for (var i = 0; i < data.length; i++) {
            var element = document.createElement("li");
            var link    = document.createElement("a");

            link.setAttribute("href", "?itemId=" + data[i].id);
            link.textContent = data[i].name;

            element.appendChild(link);
            list.appendChild(element);
        }
    };

    var init = function() {
        var textField = document.querySelector('#search-text')

        textField.addEventListener("input", function(event) {
            event.preventDefault();

            if (timer) clearTimeout(timer);

            // TODO: Trigger sendQuery on page load if text is in the input? Or on focus?
            timer = setTimeout(sendQuery, 1000);
        });
    };

    var sendQuery = function() {
        timer = null;
        var query = document.querySelector("#search-text").value;

        if (query.length < 3) return;

        var oReq = new XMLHttpRequest();
        oReq.addEventListener("load", callback);
        oReq.open("GET", "names?search=" + query);
        oReq.send();
    }

    return {
        init: init,
        getTimer: function() { return timer; }
    };
})();
