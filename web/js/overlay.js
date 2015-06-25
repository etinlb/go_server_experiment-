function Menu(){

}

Menu.prototype = {
  init: function($canvas){
    // Create the overlay object
    var $canvas = $('#mousecanvas')
    console.log($canvas.width());
    this.$overlay = $('#gamemenu')
      .addClass('overlay')
      .addClass(this.className)
      .css({position: 'absolute'})
      .width($canvas.width())
      .height($canvas.height())
      .offset($canvas.offset());
    
    // Insert the overlay immediately after the canvas object
    $canvas.after(this.$overlay);
  }
}