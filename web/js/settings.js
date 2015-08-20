// I like settings as a simple javascript object. I should make it json so it's easiy to pass down to the client but eh?

var settings =
{
  // gameCanvasId: 'gamecanvas',
  // mouseCanvasId: 'mousecanvas',
  canvasId: "gamecanvas",
  backgroundImage: null, // no background image yet

  // network settings
  webSocketUrl: "ws://localhost:8080/ws", // the websocket passed to the
  severUrl: "localhost:8080",
  PLAYER_MOVE_IMPULSE : {
    x: 10000, // FORCE, not a velocity. So pixels per second per second
    y: 10000  // FORCE, not a velocity. So pixels per second per second
  },

  // do debugy things
  debug: true,
  // Key codes for key events
  KEY: {
    BACKSPACE: 8,
    TAB:       9,
    RETURN:   13,
    ESC:      27,
    SPACE:    32,
    PAGEUP:   33,
    PAGEDOWN: 34,
    END:      35,
    HOME:     36,
    LEFT:     37,
    UP:       38,
    RIGHT:    39,
    DOWN:     40,
    INSERT:   45,
    DELETE:   46,
    ZERO:     48, ONE: 49, TWO: 50, THREE: 51, FOUR: 52, FIVE: 53, SIX: 54, SEVEN: 55, EIGHT: 56, NINE: 57,
    A:        65, B: 66, C: 67, D: 68, E: 69, F: 70, G: 71, H: 72, I: 73, J: 74, K: 75, L: 76, M: 77, N: 78, O: 79, P: 80, Q: 81, R: 82, S: 83, T: 84, U: 85, V: 86, W: 87, X: 88, Y: 89, Z: 90,
    TILDA:    192
  }
}
