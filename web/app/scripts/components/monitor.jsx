/*jshint indent: 2, node: true, nomen: true, browser: true*/
/*global React */

var Monitor = React.createClass({
  getInitialState: function () {
    return {
      message : 'Monitoring..coming soon.'
    };
  },

  render: function () {
    return (
      <div>{this.state.message}</div>
    );
  }
});

module.exports = Monitor;
