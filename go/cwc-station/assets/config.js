var wpm = document.getElementById("wpm");
var pageStatus = document.getElementById("status");
var socket = new WebSocket("wss://localhost:12345/ws");

function emit(msg) {
    socket.send("fromC:"+msg);
}

socket.onopen = function() {
    emit("status:connected");
    pageStatus.innerHTML = "Connected";
};

socket.onmessage = function(e) {
//    pageStatus.innerHTML = "Received message " + e.data;
    var res = e.data.split(":");
    var direction = res[0]
    var name = res[1]
    var value = res[2]
    if(direction == "fromC:")
        return
    if(name == "wpm") {
        wpm.innerHTML = value
    }
};

function sendWPM() {
    emit("wpm:"+wpm.innerHTML);
}

function increment() {
    var value = wpm.innerHTML;
    value++;
    wpm.innerHTML = value;
    sendWPM();
}

function decrement() {
    var value = wpm.innerHTML;
    value--;
    wpm.innerHTML = value;
    sendWPM();
}