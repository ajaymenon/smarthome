var Device = React.createClass({
  getInitialState: function() {
    return {status: this.props.status}
  },

  handleClick: function(event) {
    var newState = "-1";
    if (this.state.status == "0") {
      newState = "1";
    } else if (this.state.status == "1") {
      newState = "0";
    } else {
      return;
    }
    this.setState({status: newState});
  },

  componentDidUpdate: function() {
    $.ajax({
      url: "/api/devices/toggle",//this.props.url,
      dataType: 'json',
      type: 'GET',
      data: {"name" : this.props.name, "new_state" : this.state.status, "id" : this.props.id },
      success: function(data) {
      }.bind(this),
      error: function(xkh, status, err) {
        console.error("/api/devices/toggle", status, err.toString());
      }.bind(this)
    });
  },


  render: function() {
    var status_button = <p></p>;
    if (this.state.status != "-1") {
      var divStyle = {
        borderStyle: 'solid',
      };
      status_button = (
        <button style={divStyle} onClick={this.handleClick}>
          {this.state.status == "0" ? "off" : "on"}
        </button>
      );
    }

    if (this.state.status != "-1") {
      return (
        <div className="device">
          <h2 className="deviceId">
            Device {this.props.id} - {this.props.name}
          </h2>
          Location: {this.props.room}
          <br />
          Active Energy Usage: {this.props.energy_usage_watts} Watts
          <br/>
          {status_button}
        </div>
      );
    } else {
      return (
        <div className="device">
        </div>
      );
    }
  }
});

var DeviceBox = React.createClass({
  loadDevicesFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({devices: data["devices"], rooms: data["rooms"]});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {devices: []};
  },
  componentDidMount: function() {
    this.loadDevicesFromServer();
  },
  render: function() {
    return (
      <div className="deviceBox">
        <h1>Individual Devices - WIP</h1>
        <DeviceList devices={this.state.devices} rooms={this.state.rooms}/>
      </div>
    );
  }
});

var ControllerBox = React.createClass({
  loadDevicesFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({devices: data["devices"], rooms: data["rooms"]});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {devices: []};
  },
  componentDidMount: function() {
    this.loadDevicesFromServer();
  },
  render: function() {
    return (
      <div className="deviceBox">
        <h1>Home Controller Commands - WIP</h1>
        <button> Turn all lights off - WIP </button>
      </div>
    );
  }
});


var DeviceList = React.createClass({
  render: function() {
    var devices = [];
    for (var deviceKey in this.props.devices) {
      var device = this.props.devices[deviceKey];
      var state = "-1";
      var energy_usage_watts = 0.0;
      for (var i = 0; i < device["states"].length; i++) {
        if (device["states"][i]["variable"] == "Status") {
          state = device["states"][i]["value"];
        }
        if (device["states"][i]["variable"] == "Watts") {
          energy_usage_watts = device["states"][i]["value"];
        }
      }
    var roomName = "";
    for (var roomId in this.props.rooms) {
      var room = this.props.rooms[roomId];
      if (device["room"] == room.id) {
          roomName = room.name;
        }
      }

    devices.push(
      <Device 
        room={roomName}
        key={device["id"]} 
        id={device["id"]} 
        name={device["name"]} 
        energy_usage_watts={energy_usage_watts}
        status={state}
      >
      </Device>
     );
    }
    return (
      <div className="deviceList">
        {devices}
      </div>
    );
  }
});

ReactDOM.render(
  <DeviceBox url="/api/devices" />,
  document.getElementById('devices')
);
ReactDOM.render(
  <ControllerBox url="/api/devices" />,
  document.getElementById('controls')
);
