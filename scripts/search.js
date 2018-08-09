var Search = (function() {
    var init = function() {
        var form = document.querySelector('form')
        form.onsubmit = function(event) {
            event.preventDefault();
            console.log(event)
            console.log('hello')
        }
    }

    return {
        init: init
    };
})();