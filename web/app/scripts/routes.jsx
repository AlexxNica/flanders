var Router = require('react-router');
var Routes = Router.Routes;
var Route = Router.Route;
var Redirect = Router.Redirect;
var Search = require('./components/search');
var Monitor = require('./components/monitor');

module.exports = (
  <Routes>
    <Route name="search" path="/search" handler={Search} />
    <Route name="monitor" path="/monitor" handler={Monitor} />
    <Redirect path="/" to="search" />
    <Redirect path="/" to="search" />
  </Routes>
);
