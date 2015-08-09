import Ember from 'ember';

export default Ember.Controller.extend({
    importText: false,
    actions: {
        importList() {
            var formData = new FormData(Ember.$("#import-form")[0]);
            Ember.$.ajax({
                url: "/api/v1/todoLists/import",
                type: "POST",
                data: formData,
                processData: false,
                contentType: false,
                dataType: "json",
            }).done((data) => {
                var list = data.todoList;
                this.transitionToRoute('todo-list', list.id, list.title);
            });
        },
        fileSelected() {
            this.send('updateTitle');
        },
        updateTitle() {
            if (this.get('customTitle') || this.get('importText')) {
                return;
            }

            var fName = Ember.$("#file").val().split("\\").pop(),
                lastDot = fName.lastIndexOf(".");

            if (lastDot === -1) {
                this.set('title', fName);
            } else {
                this.set('title', fName.substr(0, lastDot));
            }
        },
        checkCustom(val) {
            this.set("customTitle", val);
            this.send('updateTitle');
        },
        checkImportText(val) {
            this.set("importText", val);
            this.set("customTitle", val);
        }
    }
});

