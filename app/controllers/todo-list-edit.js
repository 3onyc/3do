import Ember from 'ember';

export default Ember.Controller.extend({
    actions: {
        createTodoItem(group, itemTitle) {
            var item = this.store.createRecord('todo-item', {
                title: itemTitle,
            });

            group.get('items').pushObject(item);
            group.set('updatedAt', new Date());

            group.save();
            item.save();
        },
        createTodoGroup() {
            var list = this.get('model'),
                title = this.get('newGroupTitle');

            if (!title.trim()) {
                return;
            }

            var group = this.store.createRecord('todo-group', {
                title,
                list
            });

            list.get('groups').pushObject(group);

            this.set('newGroupTitle', '');

            group.save();
            list.save();
        }
    }
});
