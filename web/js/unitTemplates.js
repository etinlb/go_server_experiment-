/**
 * Game object constructors. They are just a grouping of component objects that 
 * form an actual game object. 
 */
function AttackUnit(game) {
  this.game = game;
  this.properties = {
    type: "AttackUnit",
    hp: 10
  };
  this.components = [
      new Damageable(),
      new Rect(220, 500, 20, 20),
      new Unit(),
      new Attacker(),
      new Movable(this.game.gameGrid, this.game.gridSize),
      new Selectable()
    ];
  // If there is a naming conflict in properties, set which to use here
  // this.overRide = {
  //   "update": "Attacker"
  // }
  this.overload = {"update":{}, "draw":{}};       // list of functions or properties that exist in multiple components
                                       // that will be kept and keyed to the name of the component. Y
  createEntityFromTemplate(this);   // utility.js
  this.components = [];
  this.properties = [];
}
AttackUnit.prototype = {
  update: function(){
    if(this.engaged){
      // if engaged, set the move point to the the engaged unit
      this.moveTo(this.engaged.center());    
    }
    // this.overload.update.Movable(); // calls update for movable
    this.overload.update.Movable(); // 
  },
  draw: function(context){
    if(this.selected){
      this.overload.draw.Selectable(context);
    }
    this.overload.draw.Unit(context);
  }
}

function Building(game) {
  this.game = game;
  this.properties = {
    type: "Building",
    hp: 10,
    color: "#0077ff"
  }
  this.overload = {};       // list of functions or properties that exist in multiple components
  this.components = [
    new Damageable(),
    new Rect(200, 200, 20, 20),
    new Unit()
  ]
  createEntityFromTemplate(this); // utility.js
}