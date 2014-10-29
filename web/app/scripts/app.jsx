/*jshint indent: 2, node: true, nomen: true, browser: true*/
/*global React */

var React = require('react');
var Reverter = require('./reverter');
var Nav = require('./nav');
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
  <Reverter />,
  document.getElementById('reverter')
);

React.render(
  <Nav links={links} />,
  document.getElementById('navigation')
);
