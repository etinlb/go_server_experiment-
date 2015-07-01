var isPaused; // global because fuck you. No but really, this should all be a game class

function GameState()
{
  this.paused = false;
}

GameState.prototype =
{
  update: function()
  {
    // key comes from the global
    // TODO: THIS IS UGLY AS FUCK
    if (KeyListener.isDown(KeyListener.DOWN))
    {
      this.paused = !this.paused;
    }
  },

  draw: function( canvas )
  {
    return;
  },

  //TODO: FIgure out wtf is going on with this thing
  dirty: function()
  {
    return false;
  }

}