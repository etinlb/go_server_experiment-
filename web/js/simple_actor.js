
/**
 * http://stackoverflow.com/a/105074
 * Thanks stack overflow
 */
function guid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return s4();
}

function Unit()
{
  this.x = 200;
  this.y = 200;
  this.previous_x = this.x;
  this.previous_y = this.y;
  this.id = guid();
  this.color = "#00FF00";
}

Unit.prototype =
{
  draw: function(context)
  {
    context.fillStyle = this.color;

    var x = (this.x);
    var y = (this.y);

    context.fillRect(x, y, 20, 20);
  },

  update: function()
  {
    this.saveOldState();
    // key comes from the global
    // TODO: THIS IS UGLY AS FUCK
    if (KeyListener.isDown(KeyListener.UP)) this.moveUp();
    if (KeyListener.isDown(KeyListener.LEFT)) this.moveLeft();
    if (KeyListener.isDown(KeyListener.DOWN)) this.moveDown();
    if (KeyListener.isDown(KeyListener.RIGHT)) this.moveRight();
  },

  moveLeft: function()
  {
    this.x -= 1;
  },

  moveRight: function()
  {
    this.x += 1;
  },

  moveUp: function()
  {
    this.y -= 1;
  },

  moveDown: function()
  {
    this.y += 1;
  },

  saveOldState: function()
  {
    this.previous_x = this.x;
    this.previous_y = this.y;
  },

  dirty: function()
  {
    // returns if the object has been updated at all
    // Currently just for movement
    return this.x != this.previous_x ||
           this.y != this.previous_y ;
  },

  buildPacket: function()
  {
    var packet =
    {
      x: this.x,
      y: this.y,
      id: this.id
    }
    return packet;
  },

  updatePositionFromPacket: function(packet)
  {
    this.x = packet.Rect.x;
    this.y = packet.Rect.y;
  }

}