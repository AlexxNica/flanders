
var SearchForm = React.createClass({
  getInitialState: function () {
    return {
      message : 'Searching..coming soon.'
    };
  },

  render: function() {
    return (
      <div>{this.state.message}</div>
    );
  }
});

module.exports = SearchForm;
