/**
 * Make a socket connection to the specified url.
 * @param {[type]} url           [description]
 * @param {[type]} onmessageFunc [description]
 */
function Connection( url, onmessageFunc )
{
  this.socket = new WebSocket( url )
  this.socket.onclose = function(e){
    console.log("closed");
  }
  this.socket.onopen = function(e){
    console.log("Opened");
  }

  this.socket.onmessage = function(e){
    console.log(e.data);
  }

}

Connection.prototype =
{
    send: function( data )
    {
      if ( this.socket.readyState == WebSocket.OPEN )
      {
        this.socket.send( data );
      }
    }
}