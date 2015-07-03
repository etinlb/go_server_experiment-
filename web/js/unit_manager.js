function UnitManager()
{
    this.units = {}; // keyed by game id?
    this.unitKeyField = 'id';
    this.stateObjects = {};  // keyed by....type of statemanager?
}

UnitManager.prototype = {
    addUnit: function( unit )
    {
        var key = unit[this.unitKeyField];
        if(key === undefined)
        {
            console.log(unit);
            throw "Yo, the unit you tried adding doesn't have a " + this.unitKeyField;
        }
        this.units[key] = unit;
    },

    getUnits: function()
    {
        // TODO: Is this really needed?
        return this.units;
    }
};