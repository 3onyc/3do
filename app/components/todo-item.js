import Ember from 'ember';

export default Ember.Component.extend({
    actions: {
        toggle() {
            this.send(
                this.get('item').get('done') ? 'todo' : 'done'
            );
        },
        done() {
            var item = this.get('item');
            item.set('done', true);
            item.set('doneAt', new Date());
            item.save();
        },
        todo() {
            var item = this.get('item');
            item.set('done', false);
            item.set('doneAt', null);
            item.save();
        },
    }
});
