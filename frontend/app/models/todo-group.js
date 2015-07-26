import DS from 'ember-data';

export default DS.Model.extend({
    title: DS.attr('string'),
    createdAt: DS.attr('date', {
        defaultValue: () => { new Date(); }
    }),
    updatedAt: DS.attr('date', {
        defaultValue: () => { new Date(); }
    }),

    // Relationships
    list: DS.belongsTo('todo-list', {async: true}),
    items: DS.hasMany('todo-item', {async: true}),
});
