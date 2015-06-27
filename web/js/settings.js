// I like settings as a simple javascript object. I should make it json so it's easiy to pass down to the client but eh?

var settings =
{
  // gameCanvasId: 'gamecanvas',
  // mouseCanvasId: 'mousecanvas',
  canvasId: "gamecanvas",
  backgroundImage: null, // no background image yet

  // network settings
  webSocketUrl: "ws://localhost:8080/ws",
  severUrl: "localhost:8080",

  // do debugy things
  debug: true,
}