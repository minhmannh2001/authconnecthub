import htmx from "htmx.org"
window.htmx = htmx;

htmx.on("htmx:beforeSend", function(evt) {
    evt.detail.xhr.setRequestHeader("HX-Request", "true");
})