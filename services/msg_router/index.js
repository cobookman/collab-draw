"use strict";


var PubSub = require("@google-cloud/pubsub");
var Datastore = require("@google-cloud/datastore");
var pubsub = PubSub();
var datastore = Datastore();


/**
 * Triggered on incoming drawing message.
 * This will route the incoming drawing message to all servers
 * that are subscribed to the drawing's canvas.
 */
exports.subscribe = function subscribe (context, data) {
  console.log(data.message);
  var drawing = JSON.parse(data.message);
  if (!drawing.hasOwnProperty("canvasId")) {
    context.failure("no canvas id");
    return;
  }

  var namePrefix = drawing.canvasId + ".";

  var q = datastore.createQuery("CanvasSub")
      .filter("__key__ >", namePrefix)
      .filter("__key__ <", namePrefix + "\ufffd");

  datastore.runQuery(q, (err, subs) => {
    console.log("Found " + subs.length + " subscribers to the drawing");

    subs.forEach((sub) => {
      console.log("Routing msg to topic: " + sub.topic);
      var topic = pubsub.topic(sub.topic);
      topic.publish({
        data: {
          message: data.message
        }
      }, (err) => {
        if (err) {
          console.log(err);
        }
      });
    });

    // Done sending out the messages
    context.success();
  });
};
