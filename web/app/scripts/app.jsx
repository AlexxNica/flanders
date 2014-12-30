/*jshint indent: 2, node: true, nomen: true, browser: true*/
/*global React */

var Nav = require('./components/nav');
var routes = require('./routes');
//var ClassSet = React.addons.classSet;

var links = [
  {
    href: '/search',
    title: 'Search'
  },
  {
    href: '/monitor',
    title: 'Monitor'
  }
];

React.render(
  <Nav links={links} />,
  document.getElementById('navigation')
);

React.render(
  routes,
  document.getElementById('contents')
);
