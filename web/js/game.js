$(window).load(function() {
  // add some key listeners from key_listener.js
  window.addEventListener('keyup', function(event) { KeyListener.onKeyup(event); }, false);
  window.addEventListener('keydown', function(event) { KeyListener.onKeydown(event); }, false);


  this.game = new Game();
  game.init();
  game.play();
});


// should read teh settings here maybe?
function Game() {
  this.gameObjects = [];
  this.connection = new Connection( settings.webSocketUrl );
};

Game.prototype = {
  init: function() {
    this.drawLoop = _.bind( this.drawLoop, this );
    this.gameLoop = _.bind( this.gameLoop, this );
    this.canvas = document.getElementById( settings.canvasId );
    this.context = this.canvas.getContext( '2d' );

    // read from settings maybe? idk
    var player = new Unit();
    this.addGameObject( player );

    var gameState = new GameState();
    this.addGameObject( gameState );

    if (settings.debug) {
      this.debugInit(); // Random function I shove shit in when I'm testing stuff
    }
  },

  addGameObject: function( gameObject )
  {
    console.log("adding " );
    console.log(gameObject );
    this.gameObjects.push( gameObject );
  },

  gameLoop: function() {
    // do mouse stuff
    for (var i = this.gameObjects.length - 1; i >= 0; i--) {
      // console.log( this.gameObjects[i] );
      this.gameObjects[i].update();
    };

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

  debugInit: function() {
    return;
  },

  debugDraw: function() {
    return;
  }
};

