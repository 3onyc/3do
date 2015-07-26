import Ember from 'ember';

export default Ember.TextArea.extend({
    didInsertElement() {
        this.$().focus();
    },
});
