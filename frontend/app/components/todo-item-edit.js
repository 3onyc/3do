import TodoItem from 'threedo/components/todo-item';

export default TodoItem.extend({
    editingTitle: false,
    editingDescription: false,

    actions: {
        delete() {
            var item = this.get('item');
            item.deleteRecord();
            item.save();
        },
        editDescription() {
            this.set('editingDescription', true);
        },
        saveDescription() {
            this.get('item').save();
            this.set('editingDescription', false);
        },
        editTitle() {
            this.set('editingTitle', true);

        },
        saveTitle() {
            this.get('item').save();
            this.set('editingTitle', false);
        },
    }
});
