/**
 * Contains all the components that can be used to create game objects.
 * There is no requirement of the components other than that they should
 * have a type variable as that is used for overriding specific properties in
 * the creation in unitTemplates
 */
function Unit() {
  /**
   * base unit, not used for a whole lot
   * @type {String}
   */
  this.color = "#00FF00";
  this.type = "Unit";
}

Unit.prototype = {

  click: function(ev, rightClick) {
    // Player clicked inside the canvas
    console.log("click event")
  },

  draw: function(context) {
    context.fillStyle = this.color;
    var x = (this.x) - this.game.offsetX;
    var y = (this.y) - this.game.offsetY;
    this.drawingX = x;
    this.drawingY = y;

    context.fillRect(x, y, this.width, this.height);
  },

  update: function() {

  }
}

// TODO: Evaluate if this is needed
function Selectable() {
  /**
   * Give to a unit to be "selectable" or orderable, i.e anything the player controls. Changes
   * the game menu to what ever orders you can give it? I'm not to sure
   * EX: Buildings and units.
   */
  this.type = "Selectable";
  this.selectable = true;
  this.selectColor = "#ff0000";
  this.selectWidth = 5;
}

Selectable.prototype = {
  select: function(){
    // this.unSelectColor = this.color;  // use color from unit to set as th unselectedColor
    this.selected = true;
  },

  unselect: function(){
    this.selected = false;
  },

  draw: function(context){
    context.strokeStyle = this.selectColor;
    var contextLineWidth = context.lineWidth; // significantly faster than saving the entire context.
    context.lineWidth = this.selectWidth;
    var x = (this.x) - this.game.offsetX;// - this.selectWidth;  // 
    var y = (this.y) - this.game.offsetY;// - this.selectWidth;
    // context.strokeRect(x, y, this.width+this.selectWidth*2, this.height+this.selectWidth*2);
    context.strokeRect(x, y, this.width, this.height);
    context.lineWidth = contextLineWidth;
  },
}

function Damageable() {
  /**
   * Any unit that is able to be damaged.
   * @param  {[type]} amount [description]
   * @return {[type]}        [description]
   */
  this.type = "Damageable";
}

Damageable.prototype = {
  damage: function(amount) {
    this.hp -= amount;
  }
}


function Movable(gameGrid, gridSize) {
  /**
   * Moveable unit. Utilizes A* to determine path.
   */
  this.type = "Movable";  
  this.movingTo = false;
  this.gameGrid = gameGrid; // pointer to the game grid
  this.gridSize = gridSize;
}

Movable.prototype = {
  doPath: function() {
    console.log("doing path");
    var start = this.getGrid(this.center());
    var end = this.getGrid(this.movingTo);
    var path = AStar(this.gameGrid, start, end, 'Euclidean');
    if (path.length > 1) {
      var nextStep = {
        x: path[1].x,
        y: path[1].y
      }; // next step in GRID space
      moveTowardsInGrid(this, nextStep, this.game);
      // newDirection = findAngle(nextStep,this,this.directions);    
    } else if (start[0] == end[0] && start[1] == end[1]) {
      // Reached destination grid;
      // path = [this,destination];               
      // newDirection = findAngle(destination,this,this.directions);
      return false;
    } else {
      // There is no path
      return false;
    }
    return true;
  },

  getGrid: function(point){
    /**
     * Gets the current point the unit in grid space.
     */
    var x = Math.floor(point.x / this.gridSize);
    var y = Math.floor(point.y / this.gridSize);
    return [x, y];    
  },

  getCenterOfGrid: function(gridSquare) {
    /**
     * Given a grid point, gets the center in game space.
     */
    var rect = {
      x: gridSquare.x * this.gridSize + this.gridSize / 2, // get the center,
      y: gridSquare.y * this.gridSize + this.gridSize / 2 // get the center
    }
    // rect.x = 
    // rect.y = 
    return rect;
  },


  moveTo: function(point){
    /**
     * Point to move to in pixel space
     */
    this.movingTo = point;
  },

  moveTowards: function(target) {
    var y = target.y - this.centery();
    var x = target.x - this.centerx();
    var distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2));
    var speed = 3;
    var fullCircle = Math.PI * 2;

    // what's the different between our orientation and the angle we want to face in order to move directly at our target
    var angle = Math.atan2(y, x);
    var delta = angle - this.orientation;
    var delta_abs = Math.abs(delta);

    // if the different is more than 180°, convert the angle a corresponding negative value
    if (delta_abs > Math.PI) {
      delta = delta_abs - fullCircle;
    }
    var turnSpeed = 10000;
    // if the angle is already correct, don't bother adjusting
    if (delta !== 0) {
      // do we turn left or right?
      var direction = delta / delta_abs;
      // update our orientation
      this.orientation += (direction * Math.min(turnSpeed, delta_abs));
    }
    // constrain orientation to reasonable bounds
    this.orientation %= fullCircle;

    // use orientation and speed to update our position
    this.x += Math.cos(this.orientation) * speed;
    this.y += Math.sin(this.orientation) * speed;
  },

  update: function() {
    // if (this.engaged) {
    //   this.moveTowardsEngaged();
    //   this.attack();
    //   console.log(this.collideRect(this.engaged));
    // } else 
    if (this.movingTo) {
      if (!this.doPath()) {
        this.movingTo = false;
      }
    }
  },

};


function Attacker() {
  /**
   * Attacking unit. Determines if they are close enough to attack  and engages units
   */
  this.delay = 10; // number of frames to wait until able to attack again
  this.attackDistance = 20;
  this.safeDistance = 10; // distance to keep from the 
  this.type = "Attacker";
}

Attacker.prototype = {
  engageClosest: function(entities) {
    /**
     * Find closest entity and set it to the engaged target. Entities is an
     * array of Unit Objects
     */

    // TODO: check if has a proper length
    var closest = entities[0];
    var closestDistance = distance(this, closest);
    for (var i = entities.length - 1; i >= 0; i--) {
      var distance = distance(this, entities[i]);
      if (distance < closestDistance) {
        closest = entities[i];
        closestDistance = distance;
      }
    };
    this.engaged = closest;
  },

  engageSpecific: function(entity) {
    this.engaged = entity;
  },

  attack: function() {
    if (distance(this, this.engaged) < this.attackDistance &&
      this.attackTimer == this.delay) {
      // close enough to attack
      this.engaged.damage()
      this.attackerTimer = 0;
    }
    if (this.attackTimer < this.delay) {
      this.attackTimer++;
    }
  },
}


// Utility functions
function distance(point1, point2) {
  var xs = 0;
  var ys = 0;
  // go duck typing go!
  try {
    xs = point2.centerx() - point1.centerx();
  } catch (err) {
    // I like to use centerx but not always an option, use x instead
    xs = point2.x - point1.x;
  }
  xs = xs * xs;

  try {
    ys = point2.centery() - point1.centery();
  } catch (err) {
    // I like to use centerx but not always an option, use x instead
    ys = point2.y - point1.y;
  }
  ys = ys * ys;

  return Math.sqrt(xs + ys);
}

function moveTowardsInGrid(obj, target, game) {
  // var start = game.getCenterOfGrid(obj);
  var end = game.getCenterOfGrid(target);
  moveTowards(obj, end);
}

function moveTowards(obj, target) {
  var y = target.y - obj.centery();
  var x = target.x - obj.centerx();
  var distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2));
  var speed = 3;
  var fullCircle = Math.PI * 2;

  // what's the different between our orientation
  // and the angle we want to face in order to 
  // move directly at our target
  var angle = Math.atan2(y, x);
  var delta = angle - obj.orientation;
  var delta_abs = Math.abs(delta);

  // if the different is more than 180°, convert
  // the angle a corresponding negative value
  if (delta_abs > Math.PI) {
    delta = delta_abs - fullCircle;
  }
  var turnSpeed = 10000;
  // if the angle is already correct,
  // don't bother adjusting
  if (delta !== 0) {
    // do we turn left or right?
    var direction = delta / delta_abs;
    // update our orientation
    obj.orientation += (direction * Math.min(turnSpeed, delta_abs));
  }
  // constrain orientation to reasonable bounds
  obj.orientation %= fullCircle;

  // use orientation and speed to update our position
  obj.x += Math.cos(obj.orientation) * speed;
  obj.y += Math.sin(obj.orientation) * speed;
}

