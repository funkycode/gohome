 var React = require('react');

var ZoneTypePicker = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.type || 'unknown'
        };
    },

    componentDidMount: function() {
        // If a value wasn't passed in, raise a changed notification so callers
        // can set their value accordingly since we default to unknown
        if (this.state.value === 'unknown') {
            this.props.changed && this.props.changed(this.state.value);
        }
    },
    
    selected: function(evt) {
        this.setType(evt.target.value);
    },

    setType: function(type) {
        this.setState({ value: type });
        this.props.changed && this.props.changed(type);
    },
    
    render: function() {
        var types = [
            { str: "Unknown", val:"unknown" },
            { str: "Light", val:"light" },
            { str: "Switch", val:"switch" },
            { str: "Shade", val:"shade" }
        ];
        var self = this;
        var nodes = types.map(function(type) {
            return <option value={type.val} key={type.val}>{type.str}</option>
        });
        return (
            <div className="cmp-ZoneTypePicker">
                <select
                    className="form-control"
                    onChange={this.selected}
                    value={this.state.value}>
                {nodes}
              </select>
            </div>
        );
    }
});
module.exports = ZoneTypePicker;
