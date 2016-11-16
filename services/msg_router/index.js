"use strict";

var Pubsub = require("@google-cloud/pubsub");
var Datastore = require("@google-cloud/datastore");
var datastore = Datastore();
var pubsub = Pubsub();

exports.handleDrawing = function (context, data) {
  if (!data) {
    console.log("No data, disregarding");
  }

  if (typeof(data) != "object") {
    try {
      data = JSON.parse(data);
    } catch(e) {
      console.log("Input is not json: ", e);
      console.log("Data type is: ", typeof(data));
      return context.success("input is not json, disregarding");
    }
  }

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
        return context.failure("FAIL");
      } else {
        return context.success("DONE");
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

  let query = genPrefixQuery("CanvasSub", data.canvasId + ".");
  datastore.runQuery(query, (err, canvasSubs) => {
    if (err) {
      return cb(err);
    }

    console.log("Forwarding drawing to subs: ", canvasSubs);
    var count = 0;
    var errors = [];
    canvasSubs.forEach((canvasSub) => {
      ++count;
      pubsub.topic(canvasSub.topicId).publish({data: data}, (err) => {
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

  if (!data.hasOwnProperty("drawingId") || !data.drawingId) {
    return cb("Invalid drawing object (no drawingId), not taking any action");
  }

  let key = datastore.key(["Drawing", data.drawingId]);
  datastore.save({
    key: key,
    data: [{
      "name": "canvasId",
      "value": data.canvasId
    }, {
      "name": "drawingId",
      "value": data.drawingId
    }, {
      "name": "points",
      "value": data.points
    }]
  }, (err) => {
    if (err) {
      cb(err);
    } else {
      console.log("Successfully saved drawing with id: " + data.drawingId);
      cb();
    }
  });
}

