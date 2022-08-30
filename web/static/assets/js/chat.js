window.addEventListener("load", function(evt) {
    evt.preventDefault();
    url = 'ws://localhost:8079/ws';
    c = new WebSocket(url);

    send = function(data){
        $("#output").append((new Date())+ " ==> "+data+"\n")
        c.send(data)
    }

    c.onmessage = function(msg){

        $("#output").append((new Date())+ " <== "+msg.data+"\n")
        if (msg.data.includes("sender")) {
            var data = JSON.parse(msg.data)
            if (document.getElementById("conversation_name").innerText != data["sender"]) {
                document.getElementById("chats").insertAdjacentHTML("beforeend", "<div class=\"conversation inside\"> <div class=\"message inside-msg\"> <span> "+data["message"]+"</span> <div class=\"info\"> <div class=\"time\">Just now</div><div class=\"seen-time\">Seen 1:03PM</div></div></div><div class=\"avatar\"> <img style='border-radius: 50%;' src=\""+document.getElementById("avatarpng").src+"\"> </div></div>");
            }else{
                document.getElementById("chats").insertAdjacentHTML("beforeend","<div class=\"conversation outside margin-bottom30\"> <div class=\"avatar\"> <img style=\"border-radius: 50%\" src=\""+document.getElementById("otherpfp").src+"\"/> </div><div class=\"message outside-msg\"> <span> "+data["message"]+" </span> <div class=\"info\"> <div class=\"time\">Just now</div><div class=\"seen-time\"></div></div></div></div>")

            }
            document.getElementById("chats").scrollIntoView(false, {behavior: "smooth"});
        }
        console.log(msg);
    }

    c.onopen = function(){
        send('getMessages[{"conversation": '+getQueryVariable("id")+'}]')
    }
    function getQueryVariable(variable)
    {
        var query = window.location.search.substring(1);
        var vars = query.split("&");
        for (var i=0;i<vars.length;i++) {
            var pair = vars[i].split("=");
            if(pair[0] == variable){return pair[1];}
        }
        return(false);
    }
    var botconv = false;
    if (document.getElementById("conversation_name").innerText === "Texts Helper") {
        botconv = true;
        console.log("Bot conversation");
    }
    console.log("chat.js loaded");
    function getCookie(name) {
        function escape(s) { return s.replace(/([.*+?\^$(){}|\[\]\/\\])/g, '\\$1'); }
        var match = document.cookie.match(RegExp('(?:^|;\\s*)' + escape(name) + '=([^;]*)'));
        return match ? match[1] : null;
    }
    document.getElementById("input").addEventListener("keypress", function(event) {
        if (event.key === "Enter") {
            event.preventDefault();
            console.log("Event key pressed");
            document.getElementById("send_btn").click();
        }
    });

    document.getElementById("send_btn").addEventListener("click", function(event) {
        event.preventDefault();
        if (botconv === true) {
            document.getElementById("chat-area").insertAdjacentHTML("beforeend", "<div class=\"conversation inside\"> <div class=\"message inside-msg\"> <span> "+document.getElementById("input").value+"</span> <div class=\"info\"> <div class=\"time\">Just now</div><div class=\"seen-time\">Seen 1:03PM</div></div></div><div class=\"avatar\"> <img style='border-radius: 50%;' src=\""+document.getElementById("avatarpng").src+"\"> </div></div>")
            document.getElementById("chat-area").insertAdjacentHTML("beforeend","<div class=\"conversation outside margin-bottom30\"> <div class=\"avatar\"> <img style=\"border-radius: 50%\" src=\"https://pbs.twimg.com/profile_images/1017917453735682048/D6yGqGzu_400x400.jpg\"/> </div><div class=\"message outside-msg\"> <span> Please contact our <a href='mailto:admin@localhost.com'>support team</a> for further assistance </span> <div class=\"info\"> <div class=\"time\">Just now</div><div class=\"seen-time\"></div></div></div></div>")
            document.getElementById("send_btn").disabled = true;
        }else{
            send('write[{"sender": "'+document.getElementById("user-name").innerText+'", "message": "'+document.getElementById("input").value+'", "conversation": '+getQueryVariable("id")+'}]');
        }

        console.log("Send button clicked");
    });


})
