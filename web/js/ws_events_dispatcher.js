/*
The MIT License (MIT)

Copyright (c) 2014 Ismael Celis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
-------------------------------*/
/*
Simplified WebSocket events dispatcher (no channels, no users)

var socket = new FancyWebSocket();

// bind to server events
socket.bind('some_event', function(data){
  alert(data.name + ' says: ' + data.message)
});

// broadcast events to all connected users
socket.send( 'some_event', {name: 'ismael', message : 'Hello world'} );
*/

var FancyWebSocket = function(url){
  var conn = new WebSocket(url);

  var callbacks = {};

  // the this binding is a complete hack
  this.bind = function(event_name, callback, thisBinding){
    callbacks[event_name] = callbacks[event_name] || [];
    if(thisBinding)
    {
      callback = _.bind(callback, thisBinding)
    }
    callbacks[event_name].push(callback);
    return this;// chainable
  };

  this.send = function(event_name, event_data){
    var payload = JSON.stringify({event:event_name, data: event_data});
    conn.send( payload ); // <= send JSON data to socket server
    return this;
  };

  // dispatch to the right handlers
  conn.onmessage = function(evt){
    var json = JSON.parse(evt.data)
    // check if multiple events were sent in this message
    // TODO: Figure out a better convention other than, "Is this event or events?"
    // But avoid type checking....
    //
    // TODO: Also, figure out if we want to enforce ordering? Should we be able
    // to read the events in any order, pick certian ones base on the type of
    // event or just read the array in order and trust they are sent properly
    if( json.events !== undefined )
    {
      for (var i = 0; i < json.events.length; i++) {
        eventName = json.events[i].event;
        eventData = json.events[i].data;
        dispatch(eventName, eventData)
      };
    } else {
      dispatch(json.event, json.data)
    }
  };

  this.state = function() {
    return conn.readyState;
  };

  conn.onclose = function(){dispatch('close',null)}
  conn.onopen = function(){dispatch('open',null)}

  var dispatch = function(event_name, message){
    var chain = callbacks[event_name];
    if(typeof chain == 'undefined') return; // no callbacks for this event
    for(var i = 0; i < chain.length; i++){
      chain[i]( message )
    }
  }
};