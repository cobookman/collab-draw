"use strict";

var Pubsub = require("@google-cloud/pubsub");
var Datastore = require("@google-cloud/datastore");
var datastore = Datastore();
var pubsub = Pubsub();

exports.handleDrawing = function (event, callback) {
  console.log("Received event: ", event);
  if (!event || !event.data) {
    console.log("No data, disregarding");
  }

  // parse base64 data as json obj
  var pubsubMessage = event.data;
  var data = Buffer.from(pubsubMessage.data, "base64").toString();
  if (typeof(data) != "object") {
    data = JSON.parse(data);
  }
  console.log("Received drawing: ", data);

  var counter = 0;
  var errors = [];
  var onDone = function(err) {
    --counter;
    if (err) {
      errors.push(err);
    }
  }.bind(this);

  console.log("Saving drawing: ", data);
  ++counter;
  saveDrawing(data, onDone);

  console.log("Forwarding drawing: ", data);
  ++counter;
  forwardDrawing(data, onDone);

  let intv = setInterval(()=>{
    if (counter <= 0) {
      clearInterval(intv);
      if (errors.length) {
        console.log("Errors: ", errors);
        return callback(new Error("FAIL"));
      } else {
        return callback();
      }
    }
  }, 25);
};


// Finds all records with a key which starts with given prefix.
// \ufffd is the largest UTF-8 char in regards to <, > comparisons.
function genPrefixQuery(kind, prefix) {
  return datastore.createQuery(kind)
      .filter("__key__", ">", datastore.key([kind, prefix]))
      .filter("__key__", "<", datastore.key([kind, prefix + "\ufffd"]));
}

function forwardDrawing(data, cb) {
  if (!data.hasOwnProperty("canvasId") || !data.canvasId) {
    return cb("canvasId not provided. Need to know which canvas the drawing is placed in");
  }

  let query = genPrefixQuery("Subscription", data.canvasId + ".");
  datastore.runQuery(query, (err, canvasSubs) => {
    if (err) {
      return cb(err);
    }

    console.log("Forwarding drawing to subs: ", canvasSubs);
    var count = 0;
    var errors = [];
    canvasSubs.forEach((canvasSub) => {
      ++count;
      pubsub.topic(canvasSub.serverTopicId).publish({data: data}, (err) => {
        --count;
        if (err) {
          errors.push(err);
        }
      });
    });

    let intv = setInterval(()=>{
      if (count <= 0) {
        clearInterval(intv);
        if (errors.length) {
          cb(errors);
        } else {
          console.log("Successfully forwarded drawing");
          cb();
        }
      }
    },10);
  }); // end query
}

function saveDrawing (data, cb) {
  console.log("Attempting to save drawing: ", data);

  if (!data) {
    return cb("No drawing, not saving");
  }

  if (!data.hasOwnProperty("canvasId") || !data.canvasId) {
    return cb("Invalid drawing object (no canvasId), not taking any action");
  }

  if (!data.hasOwnProperty("id") || !data.id) {
    return cb("Invalid drawing object (no id), not taking any action");
  }

  let key = datastore.key(["Drawing", data.id]);
  datastore.save({
    key: key,
    data: [{
      "name": "canvasId",
      "value": data.canvasId
    }, {
      "name": "id",
      "value": data.id
    }, {
      "name": "points",
      "value": data.points
    }]
  }, (err) => {
    if (err) {
      cb(err);
    } else {
      console.log("Successfully saved drawing with id: " + data.id);
      cb();
    }
  });
}

