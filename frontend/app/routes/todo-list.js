import Ember from 'ember';

export default Ember.Route.extend({
    setupController(controller, model) {
        if (!model) {
            return controller.set('model', this.store.find('todo-list', model.id));
        }

        if (!model._internalModel ||
            !model._internalModel._relationships ||
            !model._internalModel._relationships.initializedRelationships
        ) {
            return controller.set('model', this.store.fetchById('todo-list', model.id));
        }

        var relationships = model._internalModel._relationships.initializedRelationships;
        if (relationships.groups.canonicalState[0].dataHasInitialized) {
            return;
        }

        return controller.set('model', this.store.fetchById('todo-list', model.id));
    }
});
