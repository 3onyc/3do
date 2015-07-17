import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  primaryKey: 'ID',
  namespace: 'api/v1'
});
