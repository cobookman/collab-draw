"use strict";
var Pubsub = require("@google-cloud/pubsub");
var Datastore = require("@google-cloud/datastore");
var datastore = Datastore();
var pubsub = Pubsub();

exports.forwardDrawing = function forwardDrawing (context, data) {
  console.log("Forwarding drawing");
  if (!data.hasOwnProperty("canvasId") || !data.canvasId) {
    return context.failure("canvasId not provided. Need to know which canvas the drawing is placed in");
  }

  console.log("Finding subscriptions to canvas: ", data.canvasId);
  var subscriptionPrefix = data.canvasId + ".";
  var smallestKey = datastore.key(["CanvasSub", subscriptionPrefix]);
  var largestKey = datastore.key(["CanvasSub", subscriptionPrefix + "\ufffd"]);

  var query = datastore.createQuery("CanvasSub")
      .filter("__key__", ">", smallestKey)
      .filter("__key__", "<", largestKey);

  console.log("Running query now");
  datastore.runQuery(query, (err, canvasSubs) => {
    if (err) {
      console.log("Failed to run query:", err);
      return context.error(err);
    }

    // async publish pubsub messages
    var counter = 0;
    console.log("Fetched subscriptions: ", canvasSubs);
    canvasSubs.forEach((canvasSub) => {
      ++counter;
      var topic = pubsub.topic(canvasSub.topicId);
      console.log("Telling canvasSub about drawing", canvasSub.topicId);
      topic.publish({
        data: data,
      }, (err) => {
        --counter;
        if (err) {
          console.log("Failed to send notification");
          return context.error(err);
        }
      });
    });

    // submit an a-ok once all pubsub messages have been sent
    var inv = setInterval(()=>{
      if (counter == 0) {
        clearInterval(inv);
        context.success("Done sending messages");
      }
    }, 25);
  });
};

exports.saveDrawing = function saveDrawing (context, data) {
  console.log("Saving drawing");
  if (!data.hasOwnProperty("canvasId") || !data.canvasId) {
    return context.failure("canvasId not provided. Need to know which canvas the drawing is placed in");
  }

  if (!data.hasOwnProperty("drawingId") || !data.drawingId) {
    return context.failure("canvasId not provided. Need to know which canvas the drawing is placed in");
  }

  var key = datastore.key(["Drawing", data.drawingId]);
  console.log("Saving drawing to datastore now");
  datastore.save({
    key: key,
    data: data
  }, (err) => {
    if (err) {
      console.log("Failed to save drawing to datastore");
      context.failure(err);
    } else {
      console.log("Successfully saved drawing to datastore");
      context.success("Saved drawing");
    }
  });
};
