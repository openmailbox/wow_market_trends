var WowTrends = WowTrends || {};

WowTrends.Search = (function() {
    var BASE_ICON_URL = "https://wow.zamimg.com/images/wow/icons/large/";
    var timer = null;

    var callback = function(event) {
        data = JSON.parse(this.response);

        var list = document.querySelector("#search-results");

        while (list.firstChild !== null) {
            list.firstChild.remove();
        }

        if (data === null) {
            hideSpinner();
            return;
        }

        list.style.display = "inherit";

        for (var i = 0; i < data.length; i++) {
            var element        = document.createElement("li");
            var link           = document.createElement("a");
            var contentWrapper = document.createElement("div");
            var contentImage   = document.createElement("img");
            var contentText    = document.createElement("p");

            element.classList.add("menu-item");
            link.setAttribute("href", "history.html?itemId=" + data[i].id);
            contentWrapper.classList.add("tile", "tile-centered")
            contentImage.setAttribute("src", BASE_ICON_URL + data[i].icon + ".jpg");
            contentImage.classList.add("tile-icon");
            contentText.classList.add("tile-content");
            contentText.textContent = data[i].name;

            contentWrapper.appendChild(contentImage);
            contentWrapper.appendChild(contentText);
            link.appendChild(contentWrapper);
            element.appendChild(link);
            list.appendChild(element);
            element.appendChild(link);
        }

        hideSpinner();
    };

    var hideSpinner = function() {
        document.querySelector(".form-icon").classList.remove("loading");
    }

    var init = function() {
        var textField = document.querySelector('#search-text')

        textField.addEventListener("input", function(event) {
            event.preventDefault();

            if (timer) clearTimeout(timer);

            document.querySelector(".form-icon").classList.add("loading");

            // TODO: Trigger sendQuery on page load if text is in the input? Or on focus?
            timer = setTimeout(sendQuery, 1000);
        });
    };

    var sendQuery = function() {
        timer = null;
        var query = document.querySelector("#search-text").value;

        if (query.length < 3) {
            hideSpinner();
            return;
        }

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

window.addEventListener("load", function() {
    WowTrends.Search.init();
});