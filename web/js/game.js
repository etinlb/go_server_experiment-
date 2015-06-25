$(window).load(function() {
  this.game = new Game();
  game.init();
  game.play();
});

var settings = {
  // I'm a settings driven game design guy
  gameCanvasId: 'gamecanvas',
  mouseCanvasId: 'mousecanvas',
  debug: true,
};


function Game() {
  this.loader = new Loader('#loadingscreen')
  this.mouse = new Mouse(this);
  // this.menu = new Menu();s
  // The map is broken into square tiles of this size (20 pixels x 20 pixels)
  this.gridSize = 20;
  // Store whether or not the background moved and needs to be redrawn
  this.refreshBackground = true;
  // A control loop that runs at a fixed period of time 
  this.animationTimeout = 100; // 100 milliseconds or 10 times a second
  this.offsetX = 0; // X & Y panning offsets for the map
  this.offsetY = 0;
  this.panningThreshold = 60; // Distance from edge of canvas at which panning starts
  this.panningSpeed = 10; // 
  this.running = true;
  this.gameGrid = [];
  console.log(this.gameGrid);
  this.currentMapImage = this.loader.loadImage("images/maps/level-one.png");
  this.width = this.currentMapImage.width;
  this.height = this.currentMapImage.height;

  this.gameObjects = {};
  this.actionQueue = [];
  this.selectedObjs = [];
    // console.log(this.mouse);
};

Game.prototype = {
  init: function() {
    this.backgroundCanvas = document.getElementById(settings.gameCanvasId);
    this.foregroundCanvas = document.getElementById(settings.mouseCanvasId);
    this.backgroundContext = this.backgroundCanvas.getContext('2d');
    this.foregroundcontext = this.foregroundCanvas.getContext('2d');
    this.loader.init();
    this.mouse.init('#' + settings.mouseCanvasId); // 
    $('.gamelayer').hide();
    $('#gamestartscreen').show();
    this.state = 0;
    this.running = true;

    // Bind this forever!!!
    this.drawLoop = _.bind(this.drawLoop, this);
    this.gameLoop = _.bind(this.gameLoop, this);
    if (settings.debug) {
      this.debugInit(); // Random function I shove shit in when I'm testing stuff
    }
  },

  gameLoop: function() {
    // do mouse stuff
    if (this.mouse.eventFlag) {
      var mouseEvents = this.mouse.getMouseInfo();
      console.log(mouseEvents);
      this.handleMouseEvents(mouseEvents);
    }
    callToNestedObject(this.gameObjects, 'update');  //, this.foregroundcontext);
    window.requestAnimationFrame(this.gameLoop); //.bind(this));  
    return;
  },

  drawLoop: function() {
    if (this.running) {
      this.handlePanning();
      // fast way to clear the foreground canvas
      this.foregroundCanvas.width = this.foregroundCanvas.width;

      // draw the game objects
      callToNestedObject(this.gameObjects, 'draw', this.foregroundcontext);
      if (this.refreshBackground) {
        this.backgroundContext.drawImage(this.currentMapImage, this.offsetX, this.offsetY, this.foregroundCanvas.width,
          this.foregroundCanvas.height, 0, 0, this.foregroundCanvas.width, this.foregroundCanvas.height);
        this.refreshBackground = false;
      }

      if (settings.debug) {
        this.debugDraw();
      }
      this.mouse.draw(this.foregroundcontext);
      window.requestAnimationFrame(this.drawLoop); //.bind(this));  
    }
  },

  handleMouseEvents: function(mouseEvents){
    if (mouseEvents.drag.length > 0) {
      // should only need the last drag event
      this.selectedObjs = this.selectObjs(mouseEvents.drag[mouseEvents.drag.length - 1]);
    } else if (mouseEvents.rightClick.length > 0 && this.selectedObjs.length > 0) {
      var mouseGridPoint = mouseEvents.rightClick[mouseEvents.rightClick.length - 1];
      var movePoint = {
        x: mouseGridPoint.gameX,
        y: mouseGridPoint.gameY
      };
      for (var i = this.selectedObjs.length - 1; i >= 0; i--) {
        this.selectedObjs[i].moveTo(movePoint);
      };
    } else if (mouseEvents.click.length > 0) {
      if(this.selectedObjs.length > 0){
        this.clearSelected();  
      }
      this.getGrid(mouseEvents.click[0]);
    }
  },

  clearSelected: function(){
    callOnArray(this.selectedObjs, 'unselect');
    // this.selectObjs = [];
  },

  selectObjs: function(selectArea) {
    /**
     * Select any game objects in the selected area. Uses the center of the
     * unit to determine in in the selected area.
     */
    // make a rect to use the contains point function
    this.clearSelected();
    var rect = new Rect(selectArea.x, selectArea.y, selectArea.width, selectArea.height);
    var selectedObjs = [];
    for (var i = this.gameObjects['unit'].length - 1; i >= 0; i--) {
      if (this.gameObjects['unit'][i].selectable && rect.containsPoint(this.gameObjects['unit'][i].center())) {
        this.gameObjects['unit'][i].select();
        selectedObjs.push(this.gameObjects['unit'][i])
      }
    };
    return selectedObjs;
  },

  play: function() {
    /**
     * A bit misleading, simply starts the game and draw loop.
     */
    this.gameLoop();
    this.drawLoop();
  },

  getGrid: function(point) {
    /**
     * get the grid square of the point
     */
    var x = Math.floor(point.x / this.gridSize);
    var y = Math.floor(point.y / this.gridSize);
    return [x, y]
  },

  handlePanning: function() {
    if (this.mouse.x <= this.panningThreshold) {
      if (this.offsetX >= this.panningSpeed) {
        this.refreshBackground = true;
        this.offsetX -= this.panningSpeed;
      }
    } else if (this.mouse.x >= this.foregroundCanvas.width - this.panningThreshold) {
      if (this.offsetX + this.foregroundCanvas.width + this.panningSpeed <= this.currentMapImage.width) {
        this.refreshBackground = true;
        this.offsetX += this.panningSpeed;
      }
    }

    if (this.mouse.y <= this.panningThreshold) {
      if (this.offsetY >= this.panningSpeed) {
        this.refreshBackground = true;
        this.offsetY -= this.panningSpeed;
      }
    } else if (this.mouse.y >= this.foregroundCanvas.height - this.panningThreshold) {
      if (this.offsetY + this.foregroundCanvas.height + this.panningSpeed <= this.currentMapImage.height) {
        this.refreshBackground = true;
        this.offsetY += this.panningSpeed;
      }
    }
  },

  getCenterOfGrid: function(gridSquare) {
    /**
     * Opposite of getGrid
     */
    var rect = {
      x: 0,
      y: 0
    }
    rect.x = gridSquare.x * this.gridSize + this.gridSize / 2; // get the center
    rect.y = gridSquare.y * this.gridSize + this.gridSize / 2; // get the center
    return rect;
  },


  debugInit: function() {
    // var entity1 = this.createEntity({hp:10}, 
    //  [new Damageable(), new Rect(20, 20, 20, 20), new Unit(), new Attacker()]);
    var entity2 = new Building(this);
    var entity1 = new AttackUnit(this);
    // entity1.engageSpecific(entity2);
    var entity3 = new AttackUnit(this);
    entity3.x = entity1.x + 100;

    var addOne;
    for (var y = 0; y < this.height / this.gridSize; y++) {
      this.gameGrid[y] = [];
      for (var x = 0; x < this.width / this.gridSize; x++) {
        addOne = _.random(0, 100);
        if (addOne > 90) {
          this.gameGrid[y][x] = 1;
        } else {
          this.gameGrid[y][x] = 0;
        }
        addOne--;
      }
    }
    this.gameObjects['unit'] = [entity1, entity3]; // entity2];
    this.gameObjects['building'] = [entity2];
  },

  debugDraw: function() {
    /**
     * Draw debug stuff.
     */
    var offsetGrid = this.getGrid({
      x: this.offsetX,
      y: this.offsetY
    });

    for (var y = offsetGrid[1]; y < this.gameGrid.length; y++) {
      for (var x = offsetGrid[0]; x < this.gameGrid[y].length; x++) {
        this.backgroundContext.fillStyle = "#ffffff"
        if (this.gameGrid[y][x] == 1) {
          this.backgroundContext.fillRect(x * this.gridSize - this.offsetX, y * this.gridSize - this.offsetY, this.gridSize,
            this.gridSize);
        } else if (this.gameGrid[y][x] == 2) {
          this.backgroundContext.fillStyle = "#777777"
          this.backgroundContext.fillRect(x * this.gridSize - this.offsetX, y * this.gridSize - this.offsetY, this.gridSize,
            this.gridSize);
        } else {
          this.backgroundContext.strokeRect(x * this.gridSize - this.offsetX, y * this.gridSize - this.offsetY,
            this.gridSize, this.gridSize);
        }
      }
    };

  }
};

