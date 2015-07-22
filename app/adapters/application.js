import Ember from 'ember';
import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  primaryKey: 'ID',
  namespace: 'api/v1',
  ajaxError(jqXHR) {
      Ember.$("#server-error").text(
        "Error " + jqXHR.status + ": " + jqXHR.statusText
      );

      return this._super(jqXHR);
  }
});
