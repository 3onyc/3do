import Ember from 'ember';

export default Ember.Component.extend({
    newItemTitle: '',
    editingTitle: false,

    actions: {
        delete() {
            var group = this.get('group'),
                items = group.get('items');

            items.toArray().forEach((item) => {
                item.destroyRecord();
            });

            group.destroyRecord();
        },
        // TODO: Remove need to call outside action
        newItem() {
            var title = this.get('newItemTitle');
            if (!title.trim()) {
                return;
            }
            this.set('newItemTitle', '');
            this.sendAction('action', this.get('group'), title);
        },
        editTitle() {
            this.set('editingTitle', true);
        },
        saveTitle() {
            this.get('group').save();
            this.set('editingTitle', false);
        },
    }
});
