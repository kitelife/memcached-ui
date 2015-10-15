$(function(){

  function showResp($resp) {
    var toShow;
    if ($resp.status === "success") {
      toShow = $resp.data;
    } else {
      toShow = $resp.msg;
    }
    $("#action_result").text(toShow);
  }

  var targetServer = $("#target_server").val();

  var actionGetVM = new Vue({
    el: "#action_get",
    data: {
      k: ""
    },
    methods: {
      getIt: function($e) {
        $e.preventDefault();
        var myself = $e.targetVM;
        if (myself.k === "") {
          return;
        }
        var req = $.ajax({
            type: 'post',
            url: '/do',
            data: {
                action: "get",
                server: targetServer,
                key: myself.k
            },
            dataType: 'json'
        });
        req.done(showResp);
      }
    }
  });

  var actionSetVM = new Vue({
    el: "#action_set",
    data: {
      k: "",
      v: "",
      expTime: 0
    },
    methods: {
      setIt: function($e) {
        $e.preventDefault();
        var myself = $e.targetVM;
        var req = $.ajax({
            type: 'post',
            url: '/do',
            data: {
                action: "set",
                server: targetServer,
                key: myself.k,
                value: myself.v,
                exp_time: myself.expTime
            },
            dataType: 'json'
        });
        req.done(showResp);
      }
    }
  });

  var actionDeleteVM = new Vue({
    el: "#action_delete",
    data: {
      k: ""
    },
    methods: {
      deleteIt: function($e) {
        $e.preventDefault();
        var myself = $e.targetVM;
        var req = $.ajax({
            type: 'post',
            url: '/do',
            data: {
                action: "delete",
                server: targetServer,
                key: myself.k
            },
            dataType: 'json'
        });
        req.done(showResp);
      }
    }
  });

  var actionFlushAllVM = new Vue({
    el: "#action_flushall",
    data: {
    },
    methods: {
      flushIt: function($e) {
        $e.preventDefault();
        var myself = $e.targetVM;
        var req = $.ajax({
            type: 'post',
            url: '/do',
            data: {
                action: "flush_all",
                server: targetServer,
            },
            dataType: 'json'
        });
        req.done(showResp);
      }
    }
  });
});
