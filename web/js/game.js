$(window).load(function() {
  // add some key listeners from key_listener.js
  window.addEventListener('keyup', function(event) { KeyListener.onKeyup(event); }, false);
  window.addEventListener('keydown', function(event) { KeyListener.onKeydown(event); }, false);

  this.game = new Game();
  game.init();
});

// should read teh settings here maybe?
function Game() {
  this.gameObjects = [];
  this.connection = new FancyWebSocket( settings.webSocketUrl );
  console.log(this.connection.state());
  this.connection.bind( "add_object", this.addGameObject );
  this.connection.bind( "update", this.updateRemoteObject )
};

Game.prototype = {
  init: function() {
    // I hate javascripts hacked up object model.
    this.drawLoop = _.bind( this.drawLoop, this );
    this.gameLoop = _.bind( this.gameLoop, this );
    this.connectAndStart = _.bind( this.connectAndStart, this );

    // drawing parameters
    this.canvas = document.getElementById( settings.canvasId );
    this.context = this.canvas.getContext( '2d' );

    // read from settings maybe? idk
    var player = new Unit();
    this.addGameObject( player );

    var gameState = new GameState();
    this.addGameObject( gameState );

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
    console.log(this);
    console.log(evt);
  }

  addGameObject: function( gameObject )
  {
    console.log("adding " );
    console.log( gameObject );
    // TODO: make a new object
    this.gameObjects.push( gameObject );
  },

  gameLoop: function() {
    for (var i = this.gameObjects.length - 1; i >= 0; i--) {
      this.gameObjects[i].update();
    };

    // separate packet generate loop maybe?
    var packet = this.gameObjects[0].buildPacket();
    this.connection.send( "update", packet );

    window.requestAnimationFrame(this.gameLoop); //.bind(this));  
    return;
  },

  drawLoop: function() {
   this.canvas.width = this.canvas.width;
    for (var i = this.gameObjects.length - 1; i >= 0; i--) {
      this.gameObjects[i].draw( this.context );
    };
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

