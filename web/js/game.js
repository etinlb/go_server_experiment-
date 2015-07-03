$(window).load(function() {
  // add some key listeners from key_listener.js
  window.addEventListener('keyup', function(event) { KeyListener.onKeyup(event); }, false);
  window.addEventListener('keydown', function(event) { KeyListener.onKeydown(event); }, false);

  this.game = new Game();
  game.init();
});

// TODO: put this somewhere else
var connectionQueue = [];
// should read teh settings here maybe?
function Game() {
  this.unitManager = new UnitManager();
  this.connection = new FancyWebSocket( settings.webSocketUrl );
  console.log(this.connection.state());

  this.connection.bind( "createUnit", this.createGameObj, this );
  this.connection.bind( "update", this.updateRemoteObject, this )
};

Game.prototype = {
  init: function() {
    // I hate javascripts hacked up object model.
    this.drawLoop = _.bind( this.drawLoop, this );
    this.gameLoop = _.bind( this.gameLoop, this );
    this.connectAndStart = _.bind( this.connectAndStart, this );
    this.createGameObj = _.bind( this.createGameObj, this );
    this.updateRemoteObject = _.bind( this.updateRemoteObject, this );
    this.addGameObject = _.bind( this.addGameObject, this );

    // drawing parameters
    this.canvas = document.getElementById( settings.canvasId );
    this.context = this.canvas.getContext( '2d' );

    // read from settings maybe? idk
    var player = new Unit();
    this.addGameObject( player );

    connectionQueue.push( {"event":"createUnit", "packet" : player.buildPacket() })

    // var gameState = new GameState();
    // this.addGameObject( gameState );

    this.connectAndStart();
    if (settings.debug) {
      this.debugInit(); // Random function I shove shit in when I'm testing stuff
    }
  },

  /**
   * Updates a object from data from the server
   * @param  {[type]} evt [description]
   * @return {[type]}     [description]
   */
  updateRemoteObject: function( evt )
  {
    this.unitManager.units[evt.id].updatePositionFromPacket(evt);
  },

  createGameObj: function( gameObject )
  {
    var object = new Unit();
    object.x = gameObject.Rect.x;
    object.y = gameObject.Rect.y;
    object.id = gameObject.id;

    this.addGameObject( object );
  },

  addGameObject: function( gameObject )
  {
    this.unitManager.addUnit( gameObject );
  },

  gameLoop: function() {
    // TODO: Abstract message sending better
    var packets = [];
    for( var id in this.unitManager.units )
    {
      var unit = this.unitManager.units[id];
      unit.update();

      // check if we need to update the server
      if(this.unitManager.units[id].dirty())
      {
        //TODO: this is awful
        // packets.push(this.gameObjects[i].buildPacket())
        var packet = unit.buildPacket();
        this.connection.send( "update", packet );
      }
    }

    window.requestAnimationFrame(this.gameLoop); //.bind(this));  
    return;
  },

  drawLoop: function() {
    this.canvas.width = this.canvas.width;
    for( var id in this.unitManager.units )
    {
      this.unitManager.units[id].draw( this.context );
    }
    window.requestAnimationFrame( this.drawLoop ); //.bind(this));
  },

  play: function() {
    /**
     * A bit misleading, simply starts the game and draw loop.
     */
    this.gameLoop();
    this.drawLoop();
  },

  /**
   * Waits for the connection to be open, then calls the play function
   */
  connectAndStart: function() {
    if (this.connection.state() != WebSocket.OPEN){
      setTimeout( this.connectAndStart, 100);
    } else {
      // send the queued messages
      for (var i = connectionQueue.length - 1; i >= 0; i--) {
        this.connection.send( connectionQueue[i]["event"], connectionQueue[i]["packet"] );
      };
      this.play();
    }
  },

  debugInit: function() {
    return;
  },

  debugDraw: function() {
    return;
  }
};

