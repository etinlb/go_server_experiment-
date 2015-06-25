function Unit()
{
  this.x = 400;
  this.y = 400;
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

}