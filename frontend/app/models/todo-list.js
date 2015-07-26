import DS from 'ember-data';

export default DS.Model.extend({
    title: DS.attr('string'),
    description: DS.attr('string'),
    createdAt: DS.attr('date', {
        defaultValue: () => { new Date(); }
    }),
    updatedAt: DS.attr('date', {
        defaultValue: () => { new Date(); }
    }),

    // Relationships
    groups: DS.hasMany('todo-group', {async: true}),
});
