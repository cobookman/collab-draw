'use strict';

var proxyquire = require('proxyquire').noCallThru();

function getSample () {
  var topicMock = {
    publish: sinon.stub().callsArg(1)
  };
  var pubsubMock = {
    topic: sinon.stub().returns(topicMock)
  };
  var PubSubMock = sinon.stub().returns(pubsubMock);
  return {
    sample: proxyquire('../', {
      '@google-cloud/pubsub': PubSubMock
    }),
    mocks: {
      PubSub: PubSubMock,
      pubsub: pubsubMock,
      topic: topicMock
    }
  };
}

function getMockContext () {
  return {
    success: sinon.stub(),
    failure: sinon.stub()
  };
}

describe('functions:forwardDrawing', function () {
  it('Publish fails without canvasId specified', function () {
    var expectedMsg = "canvasId not provided. Need to know which canvas the drawing is placed in";
    var context = getMockContext();

    getSample().sample.forwardDrawing(context, {
      message: {}
    });

    assert.equal(context.failure.calledOnce, true);
    assert.equal(context.failure.firstCall.args[0], expectedMsg);
    assert.equal(context.success.called, false);
  });
});

