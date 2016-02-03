var React = require('react');
var ReactDOM = require('react-dom');
var Scene = require('./Scene.jsx');
var SceneInfo = require('./SceneInfo.jsx');

var SceneList = React.createClass({
    getInitialState: function() {
        return {
            editMode: false,
            scenes: this.props.scenes,
            zones: this.props.zones,
        };
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.scenes) {
            this.setState({ scenes: nextProps.scenes });
        }
        if (nextProps.zones) {
            this.setState({ zones: nextProps.zones });
        }
    },
    
    componentDidMount: function() {

        //TODO: Needed?
        $.ajax({
            url: '/api/v1/systems/123/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({zones: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });

        return;
        //TODO: Enable as part of a mode
        var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        Sortable.create(el);
    },

    edit: function() {
        this.setState({ editMode: true });
    },

    sceneDeleted: function(sceneId) {
        var scenes = this.state.scenes;
        for (var i=0; i<scenes.length; ++i) {
            console.log(sceneId + " : " + scenes[i].id);
            if (scenes[i].id === sceneId) {
                scenes.splice(i, 1);
                console.log(scenes);
                this.setState({ scenes: scenes });
                break;
            }
        }
    },

    render: function() {

        var body;
        var self = this;
        if (this.state.editMode) {
            body = this.state.scenes.map(function(scene) {
                return (
                    <SceneInfo
                      onDestroy={self.sceneDeleted}
                      scenes={self.state.scenes}
                      zones={self.state.zones}
                      buttons={self.props.buttons}
                      scene={scene}
                      key={scene.id} />
                );
            });
        } else {
            body = this.state.scenes.map(function(scene) {
                return (
                    <Scene scene={scene} key={scene.id}/>
                );
            });
        }
        
        //TODO: Add loading
        return (
            <div className="cmp-SceneList">
              <div className="clearfix editButtonWrapper">
                <button className="btn btn-primary btnEdit pull-right" onClick={this.edit}>Edit</button>
              </div>
              {body}
            </div>
        );
    }
});
module.exports = SceneList;

//TODO existing scene edit:
//2. edit name + save
//3. make id readonly
//4. set address
//5. delete existing command
//8. Test button

//TODO: Add new scene