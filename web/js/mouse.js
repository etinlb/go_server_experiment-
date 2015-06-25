function Mouse(game) {
  // x,y coordinates of mouse relative to top left corner of canvas
  this.x = 0;
  this.y = 0;
  // x,y coordinates of mouse relative to top left corner of game map
  this.gameX = 0;
  this.gameY = 0;
  // game grid x,y coordinates of mouse 
  this.gridX = 0;
  this.gridY = 0;
  // whether or not the left mouse button is currently pressed
  this.buttonPressed = false;
  // whether or not the player is dragging and selecting with the left mouse button pressed
  this.dragSelect = false;
  // whether or not the mouse is inside the canvas region
  this.insideCanvas = false;
  this.clearQueue(); // sets the event queue
  this.game = game;
}

Mouse.prototype = {
  click: function(ev, rightClick) {
    this.dragSelect = false;
    var clickType = rightClick ? "rightClick" : "click";
    this.addToQueue(clickType, {
      x: this.x,
      y: this.y,
      gameX: this.gameX,
      gameY: this.gameY
    })
  },

  mouseMove: function(ev) {
    var offset = this.$mouseCanvas.offset();
    // offset is game world offset
    this.x = ev.pageX - offset.left;
    this.y = ev.pageY - offset.top;

    this.calculateGameCoordinates(0, 0);
    if (this.buttonPressed) {
      if ((Math.abs(this.dragX - this.gameX) > 4 || Math.abs(this.dragY - this.gameY) > 4)) {
        this.dragSelect = true
      }
    } else {
      this.dragSelect = false;
    }
    // this.addToQueue("move", [this.x, this.y]); // TODO: put map scrolling in
  },

  mouseDown: function(ev) {
    if (ev.which == 1) {
      this.buttonPressed = true;
      this.dragX = this.gameX;
      this.dragY = this.gameY;
      ev.preventDefault();
    }
    // this.addToQueue("down", {}); // 
    return false;
  },

  mouseUp: function(ev) {
    var shiftPressed = ev.shiftKey;
    // console.log(ev);
    if (ev.which == 1) {
      if (this.dragSelect) {
        this.addToQueue("drag", this._getBox()); // give up event since it was a mouse drag
      }
      //Left key was released                
      this.buttonPressed = false;
      this.dragSelect = false;
    }
    return false;
  },

  clearQueue: function() {
    /**
     * I don't know if this is the best way to handle mouse input but it makes sense to me.
     * The events are supposed to be defined as follows, actual implementation may vary depending
     * on how much well this coffee high lasts
     * click : a click event, gets when an browser down and up even happen and the mouse didn't move in between
     * move  : when the mouse moves
     * @type {Object}
     */
    this.eventQueue = {
      click: [],
      rightClick: [],
      move: [],
      down: [],
      up: [],
      drag: []
    };
    this.eventFlag = 0; // low means nothing to process
  },

  addToQueue: function(whichQueue, infoObj) {
    /**
     * add the info obj to the queue
     */
    this.eventQueue[whichQueue].push(infoObj);
    this.eventFlag = 1;
  },

  getMouseInfo: function() {
    /**
     * Returns the current state of the mouse
     */
    var queue = this.eventQueue;
    this.clearQueue()
    return queue;
  },

  draw: function(context) {
    if (this.dragSelect) {
      // console.
      var rectBox = this._getBox();
      context.strokeStyle = 'white';
      context.strokeRect(rectBox.x - this.game.offsetX, rectBox.y - this.game.offsetY, rectBox.width, rectBox.height);
    }
    context.strokeStyle = 'white';
    context.strokeRect(this.x, this.y, 5, 5);

  },

  _getBox: function() {
    /**
     * returns the mouse drag select box
     */
    var rect = {}
    rect.x = Math.min(this.gameX, this.dragX);
    rect.y = Math.min(this.gameY, this.dragY);
    rect.width = Math.abs(this.gameX - this.dragX)
    rect.height = Math.abs(this.gameY - this.dragY)
    return rect;
  },
  // _getDrawBox: function() {
  //   /**
  //    * returns the mouse drag select box
  //    */
  //   var rect = {}
  //   rect.x = Math.min(this.gameX, this.dragX);
  //   rect.y = Math.min(this.gameY, this.dragY);
  //   rect.width = Math.abs(this.gameX - this.dragX)
  //   rect.height = Math.abs(this.gameY - this.dragY)
  //   return rect;
  // },

  calculateGameCoordinates: function(offsetX, offsetY) {
    var gridSize = 20; // Grid size is just a way to get which tile you are on
    this.gameX = this.x + offsetX;
    this.gameY = this.y + offsetY;
    // console.log(this);
    this.gameX = this.x + this.game.offsetX;
    this.gameY = this.y + this.game.offsetY;

    this.gridX = Math.floor((this.gameX) / this.game.gridSize); // do something with gridSize
    this.gridY = Math.floor((this.gameY) / this.game.gridSize);
    // this.gridX = Math.floor((this.gameX) / gridSize); // do something with gridSize
    // this.gridY = Math.floor((this.gameY) / gridSize);
  },

  init: function(canvasId, callback) {
    /**
     * Set the jquery callbacks and remember the mouse canvas
     * @type {[type]}
     */
    this.$mouseCanvas = $(canvasId);
    // console.log($mouseCanvas);
    var self = this; // Significantly faster than binding to this http://jsperf.com/bind-vs-closure-setup/6

    this.$mouseCanvas.mousemove(function(ev) {
      self.mouseMove(ev);
    });
    this.$mouseCanvas.click(function(ev) {
      self.click(ev, false);
      return false;
    });

    this.$mouseCanvas.mousedown(function(ev) {
      self.mouseDown(ev);
    });

    this.$mouseCanvas.bind('contextmenu', function(ev) {
      self.click(ev, true);
      return false;
    });

    this.$mouseCanvas.mouseup(function(ev) {
      self.mouseUp(ev);
    });

    // this.$mouseCanvas.mouseleave(function(ev) {
    //     mouse.insideCanvas = false;
    // });

    // $mouseCanvas.mouseenter(function(ev) {
    //     mouse.buttonPressed = false;
    //     mouse.insideCanvas = true;
    // });
  }
}