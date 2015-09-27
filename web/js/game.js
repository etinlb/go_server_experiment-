$(window).load(function() {
  // add some key listeners from key_listener.js
  window.addEventListener('keyup', function(event) { KeyListener.onKeyup(event); }, false);
  window.addEventListener('keydown', function(event) { KeyListener.onKeydown(event); }, false);

  this.game = new Game();
  game.init();
});
var debug = true;
// TODO: put this somewhere else
var connectionQueue = [];
var lastNetworkingUpdate = new Date();
var lastDrawUpdate = new Date();
var networkingFps = 30;
var drawingFps = 30;

// should read teh settings here maybe?
function Game() {
  this.unitManager = new UnitManager();
  this.connection = new FancyWebSocket( settings.webSocketUrl );
  console.log(this.connection.state());
  this.events = [];

  // Probably the neatest javascript library I've seen. I hate user input with javascript...
  this.setupKeyListener();

  // TODO: Need to make a sync command from the server so it can send all the
  // game objects that were previously created.
  this.connection.bind( "createPlayer", this.createGameObj, this );
  this.connection.bind( "createObject", this.createGameObj, this );
  this.connection.bind( "update", this.updateRemoteObjects, this )
  this.connection.bind( "sync", this.sync, this )
};

Game.prototype = {
  // TODO: Why is this init and not in the constructor?
  init: function() {
    // I hate javascripts hacked up object model.
    this.drawLoop = _.bind( this.drawLoop, this );
    this.gameLoop = _.bind( this.gameLoop, this );
    this.connectAndStart = _.bind( this.connectAndStart, this );
    this.createGameObj = _.bind( this.createGameObj, this );
    this.updateRemoteObjects = _.bind( this.updateRemoteObjects, this );
    this.addGameObject = _.bind( this.addGameObject, this );
    // this.handleKeyDown = _.bind( this.handleKeyDown, this );
    // this.handleKeyUp = _.bind( this.handleKeyUp, this );

    // drawing parameters
    this.canvas = document.getElementById( settings.canvasId );
    this.context = this.canvas.getContext( '2d' );

    // frame rate debug
    this.drawFrameRate = new FrameRateTracker("drawfps");
    this.networkRate = new FrameRateTracker("networkfps");


    // TODO: read from settings maybe?
    // TODO: Actually, move to the component base system you had in the other failed game
    // and have a component called Player Controlled object or something
    this.player = new Unit();
    this.addGameObject( this.player );

    connectionQueue.push( {"event":"createPlayer", "packet" : this.player.buildPacket() })

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
  updateRemoteObjects: function( evt )
  {
    for (var i = evt.length - 1; i >= 0; i--) {
      var id = evt[i].id;
      this.unitManager.units[id].updatePositionFromPacket(evt[i]);
    };

    if(debug){
      this.debugNetwork();
    }

  },

  sync: function( objectArr )
  {
    console.log("syncing");
    console.log(objectArr.length);
    for(var i = objectArr.length -1; i >= 0; i-- )
    {
      console.log(objectArr[i]);
      this.createGameObj(objectArr[i]);
    }

  },

  // TODO: Switch based on the type field of the game object
  createGameObj: function( gameObject )
  {
    var object = new Unit();
    // object.x = gameObject.Rect.x;
    // object.y = gameObject.Rect.y;
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

      // // check if we need to update the server
      // if(this.unitManager.units[id].dirty())
      // {
      //   //TODO: this is awful
      //   // packets.push(this.gameObjects[i].buildPacket())
      //   var packet = unit.buildPacket();
      //   this.connection.send( "move", packet );
      // }
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

    if(debug){
      this.debugDraw();
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
      // TODO: Properly pool the messages to send only once
      // send the queued messages
      for (var i = connectionQueue.length - 1; i >= 0; i--) {
        this.connection.send( connectionQueue[i]["event"], connectionQueue[i]["packet"] );
      };
      this.play();
    }
  },

  processEvents: function() {
  },

  addEvent: function() {

  },

  upKeyDownEvent: function(evt){
    var key = evt.keyCode;
    switch(key) {
      case settings.KEY.LEFT:
        this.sendMoveEvent(-1, 0);
        break;
      case settings.KEY.RIGHT:
        this.sendMoveEvent(1, 0);
        break;
      case settings.KEY.UP:
        this.sendMoveEvent(0, -1);
        break;
      case settings.KEY.DOWN:
        this.sendMoveEvent(0, 1);
        break;
      case settings.KEY.SPACE:
        break;
    }
  },

  upKeyUpEvent: function(evt){
    var key = evt.keyCode;
    switch(key) {
      case settings.KEY.LEFT:
        this.sendMoveEvent(1, 0);
        break;
      case settings.KEY.RIGHT:
        this.sendMoveEvent(-1, 0);
        break;
      case settings.KEY.UP:
        this.sendMoveEvent(0, 1);
        break;
      case settings.KEY.DOWN:
        this.sendMoveEvent(0, -1);
        break;
      case settings.KEY.SPACE:
        break;
    }
  },

  // sends a move event for the player.
  // TODO: batch messages like this to send all at once from the client
  sendMoveEvent: function(forceX, forceY) {
    console.log("Sending move event with this force");
    var packet = {
      xVel: forceX,
      yVel: forceY,
      id: this.player.id
    };
    this.connection.send( "move", packet );
  },

  // TODO: component maybe?
  setupKeyListener: function() {
    var my_defaults = {
      prevent_repeat  : true,
      this            : this,
      on_keydown      : this.upKeyDownEvent,
      on_keyup        : this.upKeyUpEvent
    };

    this.keyListener = new window.keypress.Listener( "", my_defaults);
    this.keyListener.register_many([
              { "keys" : "up" },
              { "keys" : "down" },
              { "keys" : "right" },
              { "keys" : "left" }
    ]);
  },

  debugInit: function() {
    return;
  },

  debugDraw: function() {
    this.drawFrameRate.updateFrameRate();
    return;
  },

  debugNetwork: function(){
    this.networkRate.updateFrameRate();
  }
};
