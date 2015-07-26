import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.resource('todo-lists', { path: '/todo' });
  this.resource('todo-list', { path: '/todo/:id/:title' });
  this.resource('todo-list-edit', { path: '/todo/edit/:id/:title' });
});

export default Router;
