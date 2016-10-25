"use strict";

var PubSub = require("@google-cloud/pubsub");

exports.forwardDrawing = function forwardDrawing (context, data) {
  console.log("Received pubsub msg", data);
  try {
    if (!data.hasOwnProperty("canvasId") || !data.canvasId) {
      throw new Error("canvasId not provided. Need to know which canvas the drawing is placed in");
    }
    // TODO(bookman): Actually forward drawings
    context.error("TODO(bookman): Actually forward drawings");
    //context.success("Forwarded Drawing");
  } catch (err) {
    console.error(err);
    return context.failure(err.message);
  }
  // TODO(bookman): Implement this
};

exports.saveDrawing = function saveDrawing (context, data) {
  console.log("Received pubsub msg", data);
  context.success("Saved Drawing");
  // TODO(bookman): Implement this
};
