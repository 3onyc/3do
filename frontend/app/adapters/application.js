import Ember from 'ember';
import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  primaryKey: 'ID',
  namespace: 'api/v1',
  ajaxError(jqXHR) {
      Ember.$("#server-error").text(
        "Error " + jqXHR.status + ": " + jqXHR.statusText
      ).fadeIn(200).delay(2000).fadeOut(200);

      return this._super(jqXHR);
  }
});
