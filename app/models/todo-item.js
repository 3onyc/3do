import DS from 'ember-data';

export default DS.Model.extend({
    title: DS.attr('string'),
    description: DS.attr('string'),
    done: DS.attr('boolean', {defaultValue: false}),
    doneAt: DS.attr('date', {defaultValue: null}),
    createdAt: DS.attr('date', {
        defaultValue() { new Date(); }
    }),
    updatedAt: DS.attr('date', {
        defaultValue() { new Date(); }
    }),

    // Relationships
    group: DS.belongsTo('todo-group', {async: true}),

    // Computed
    hasDescription: function() {
        var description = this.get('description');

        return description && !!description.trim();
    }.property('description'),
});
