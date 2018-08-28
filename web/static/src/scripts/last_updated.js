var callback = function(_evt) {
    var pattern = /(\d+)-(\d+)-(\d+)T(\d+):(\d+):(\d+)/.exec(this.response)
    var date    = null;

    if (pattern) {
        date = new Date(pattern[1], pattern[2] - 1, pattern[3], pattern[4], pattern[5], pattern[6]);
    }

    document.querySelector("#last-updated").textContent = "Last updated " + date.toGMTString();
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

export default function FetchTimestamp() {
    var oReq = new XMLHttpRequest();
    oReq.addEventListener("load", callback);
    oReq.open("GET", "lastUpdated");
    oReq.send();
};