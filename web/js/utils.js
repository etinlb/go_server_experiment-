function FrameRateTracker(frameId, filter){
    this.frameId = frameId;

    // The higher this value, the less the FPS will be affected by quick changes
    // Setting this to 1 will show you the FPS of the last sampled frame only
    if(filter === undefined){
        filter = 50;
    }

    this.filter = filter;
    this.fps = 30;
    this.lastUpdate = new Date();
}

FrameRateTracker.prototype =  {
    updateFrameRate: function() {
        var self = this;
        thisFrameFPS = 1000 / ((now=new Date()) - this.lastUpdate);

        if (now!=this.lastUpdate){
            this.fps += (thisFrameFPS - this.fps) / this.filter;
            this.lastUpdate = now;
        }

        var fpsOut = document.getElementById(this.frameId);
        setInterval(function(){
          fpsOut.innerHTML = self.frameId + ": " + self.fps.toFixed(1);
        }, 1000);
    },

    drawFrameRate: function() {

    }
};
