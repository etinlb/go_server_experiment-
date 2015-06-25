function applyToNestedObject(nestedObj, func) {
  for (type in nestedObj) {
    for (var i = nestedObj[type].length - 1; i >= 0; i--) {
      func(nestedObj[type][i]);
    };
  }
}

function callToNestedObject(nestedObj, func, args) {
  for (type in nestedObj) {
    for (var i = nestedObj[type].length - 1; i >= 0; i--) {
      nestedObj[type][i][func].call(nestedObj[type][i], args); // I think there could be a better way to do this 
    };
  }
}

function callOnArray(array, func, args){
  for (var i = array.length - 1; i >= 0; i--) {
    array[i][func].call(array[i], args); // I think there could be a better way to do this 
  };  
}

function createEntity(properties, components) {
  var prop;
  var entity = {
    properties: {},
    components: []
  }
  return createEntityFromTemplate(entity);
}

function createEntityFromTemplate(templateObj) {
  /**
   * A way to fake inheritance with composition. This could be very slow now that I'm doing it more
   * but it is kind of a cool idea. I need to know the speed of _.bind on subsequent calls to the function
   */
  var prop;
  console.log(templateObj);
  for (prop in templateObj.properties) {
    templateObj[prop] = templateObj.properties[prop];
  }

  templateObj.components.forEach(function(component) {
    for (prop in component) {
      if (templateObj.hasOwnProperty(prop)) {
        // check overriding
        console.log("has " + prop);
        if(templateObj.overRide !== undefined && templateObj.overRide[prop] === component.type ){
          console.log("setting precedence" + prop)
          templateObj[prop] = component[prop];          
        } 
          // throw "Entity property conflict! " + prop;
      }else if (templateObj.overload[prop] !== undefined){
        templateObj.overload[prop][component.type] = component[prop];
        templateObj.overload[prop][component.type] =  _.bind(templateObj.overload[prop][component.type], templateObj);
        console.log("here");
      } else{
        templateObj[prop] = component[prop];
      }
    }
  });
  return templateObj;
}

Engine = function(game){
  /**
   * Various engine type interactions
   * @type {[type]}
   */
  this.game = game; 
}