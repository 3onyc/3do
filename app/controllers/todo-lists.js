import Ember from 'ember';

export default Ember.Controller.extend({
    newTitle: '',
    creatingList: false,

    actions: {
        createList() {
            this.set('creatingList', true);
        },
        saveList() {
            var title = this.get('newTitle');
            if (!title.trim()) {
                return;
            }

            var list = this.store.createRecord('todo-list', {
                title,
            });

            this.set('newTitle', '');
            this.set('creatingList', 'false');

            list.save().then(() => {
                this.transitionTo('todo-list-edit', list);
            });
        }
    }
});
